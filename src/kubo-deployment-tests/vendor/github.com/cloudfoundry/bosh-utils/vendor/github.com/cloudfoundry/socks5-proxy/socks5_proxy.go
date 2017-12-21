package proxy

import (
	"errors"
	"fmt"
	"net"
	"strconv"

	socks5 "github.com/armon/go-socks5"

	"golang.org/x/crypto/ssh"
	"golang.org/x/net/context"
)

var netListen = net.Listen

//go:generate counterfeiter . Proxy
type Proxy interface {
	Start(string, string) error
	Addr() (string, error)
}

type DialFunc func(network, address string) (net.Conn, error)

type Socks5Proxy struct {
	hostKeyGetter KeyGetter
	port          int
	started       bool
}

func NewSocks5Proxy(hostKeyGetter KeyGetter) *Socks5Proxy {
	return &Socks5Proxy{
		hostKeyGetter: hostKeyGetter,
		started:       false,
	}
}

func (s *Socks5Proxy) Start(key, url string) error {
	if s.started {
		return nil
	}

	dialer, err := s.Dialer(key, url)
	if err != nil {
		return err
	}

	err = s.StartWithDialer(dialer)
	if err != nil {
		return err
	}

	return nil
}

func (s *Socks5Proxy) Dialer(key, url string) (DialFunc, error) {
	signer, err := ssh.ParsePrivateKey([]byte(key))
	if err != nil {
		return nil, fmt.Errorf("parse private key: %s", err)
	}

	hostKey, err := s.hostKeyGetter.Get(key, url)
	if err != nil {
		return nil, fmt.Errorf("get host key: %s", err)
	}

	clientConfig := &ssh.ClientConfig{
		User: "jumpbox",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.FixedHostKey(hostKey),
	}

	conn, err := ssh.Dial("tcp", url, clientConfig)
	if err != nil {
		return nil, fmt.Errorf("ssh dial: %s", err)
	}

	return conn.Dial, nil
}

func (s *Socks5Proxy) StartWithDialer(dialer DialFunc) error {
	conf := &socks5.Config{
		Dial: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer(network, addr)
		},
	}

	server, err := socks5.New(conf)
	if err != nil {
		return fmt.Errorf("new socks5 server: %s", err) // not tested
	}

	if s.port == 0 {
		s.port, err = openPort()
		if err != nil {
			return fmt.Errorf("open port: %s", err)
		}
	}

	go func() {
		err = server.ListenAndServe("tcp", fmt.Sprintf("127.0.0.1:%d", s.port))
		if err != nil {
			// untested; commands that require the proxy will return errors
			fmt.Printf("socks5 proxy: %s", err.Error())
		}
	}()

	s.started = true
	return nil
}

func (s *Socks5Proxy) Addr() (string, error) {
	if s.port == 0 {
		return "", errors.New("socks5 proxy is not running")
	}
	return fmt.Sprintf("127.0.0.1:%d", s.port), nil
}

func openPort() (int, error) {
	l, err := netListen("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	defer l.Close()
	_, port, err := net.SplitHostPort(l.Addr().String())
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(port)
}
