package proxy

import (
	"net"

	"golang.org/x/crypto/ssh"
)

//go:generate counterfeiter . KeyGetter
type KeyGetter interface {
	Get(string, string) (ssh.PublicKey, error)
}

type HostKeyGetter struct {
	publicKeyChannel chan ssh.PublicKey
	dialErrorChannel chan error
}

func NewHostKeyGetter() HostKeyGetter {
	return HostKeyGetter{
		publicKeyChannel: make(chan ssh.PublicKey),
		dialErrorChannel: make(chan error),
	}
}

func (h HostKeyGetter) Get(key, serverURL string) (ssh.PublicKey, error) {
	signer, err := ssh.ParsePrivateKey([]byte(key))
	if err != nil {
		return nil, err
	}

	clientConfig := &ssh.ClientConfig{
		User: "jumpbox",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: h.keyScanCallback,
	}

	go func() {
		conn, err := ssh.Dial("tcp", serverURL, clientConfig)
		if err != nil {
			h.publicKeyChannel <- nil
			h.dialErrorChannel <- err
			return
		}
		defer conn.Close()
		h.dialErrorChannel <- nil
	}()

	return <-h.publicKeyChannel, <-h.dialErrorChannel
}

func (h HostKeyGetter) keyScanCallback(hostname string, remote net.Addr, key ssh.PublicKey) error {
	h.publicKeyChannel <- key
	return nil
}
