package proxy_test

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"testing"

	"golang.org/x/crypto/ssh"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	sshPrivateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAoxA4a1jmeJHabb+ACuIjPdAEURB9GTbwjHJzJrGWw6ppQCKS
QtCzgLY2WdKzZef9G0eSKF7YgVtLvMHM6O0ph4eChtH1gArMpU46DcOZH8LXvrlX
d4z2aD4zvDwDbCP8HszOSCmQIy9UVxYD6wsOw4xBvOP0EDGK5jtM7nZ/vg9bMRtK
SP8Z59cDtiy2EeDezRAW2t1i8v5Uy9Bg2fjPl6g1b1PI64F1nm5IOKvb8vyT9GRT
VnSRKUypREFJr6QFVo8xvLc35f8pe9sJSgxJXzEjf43/FLL5eAIgmxpwMhAks5yd
/a4jv8EBXL9ZUMUjXJHzf4xmkNkkHTmxEhgD7wIDAQABAoIBAAmDuMcKuOfwGr6s
ndwEtem1aYsRWztNaVvIkc+ALTvdhaaoXcBoTREFkMZM9QrNLoeY9X5FinyBxzmM
VViB/hpaXdNgDOMbvjUnC1wiPZ0M0WnfhqsDHp2Wg45IMirtLpjdemvbgP2MlW8/
aZsdWg8u7+cFpggL2/7zFtoTMAD27JM7gk/06IN4vx+G2G7c0g4T0Xn+redTj9Hq
V9mpy/fke/esi4JgLqdw8l7P1qqQZtrl6imbPrvBdb0/moIW6w8txUofTOdlBUN2
XJlaL7JjdWfTN7P80KGIklLJiDmlPyS267fNxOWpDOjGMxe1Ctzx12MzS8Jjn2S2
6zq6XIECgYEAzbOZP9gg+3/PTfpvtCujZSCYuk8GXMjgoxsx4/ongaA1cLmpddvW
My/lGvGkwJ7b2JviTUQLV+TzvU8f6S8J4yOeYp6IuUIerZaEn78YFWXWXbcUNdkv
TIZ/j6YrCGSsX27hI7OJcK8ZjdGVgFGJh4tShDcWFzioNMO4wjws8YsCgYEAyu+T
4Phq8NugtS3zGic6M1IZIcfQEYt9ngr8F+t9YSJDYVoMp8CxngngWkfzda2hLCEG
gUitN7oonuw0I1WpUf3dOsYKrVJXmmQGEszo4swn6F//tXhEK3KAKnBkfhVTHYFh
eU3K8gg48/T/R3annuXR47LcIwwurGGmbR5he60CgYEAgHBo+yVfisoGTiFWiEBL
OQS+eG6JgXvoT8/WOgxjiJvZYnZ7Kl1HBRUdz9IcVi2bBkhnaGlZT9tkmcsDGN3H
Ja2C4v8sTcjMUQVP8FMonYvF6yQ6mVjwIK9GjRJrgkUiIECikWE0K0kaAqRf3gyL
fDfxIR8oSv2UgcXH4ngic/sCgYA3/oD8Ky8+xCsEsugICFjbvkRm+L4liSqhCADl
DLosqgqTewhQ5S9dHvaDkqTPjJgTGA22cHozDS+WIjCEq2cr03NOe0SI7FZ1qDGw
0E9V/OTqDkr9JHES1+YbT6W60GF9m6xsjxV3UON+FNS3QDsh8eHHBRwOo5bhQ5Rr
OV3GhQKBgH5ZxIGKybkjkqrrzY/sDVrmteTHAPrBCcrCW3bWHKjRqCAUfML3ixvV
5wTm9J7ak28ylLR+ESqhNE5Shqga3cc7jvZQJ3MEg3oKgrItpCH0JaOtQJ+g2S6V
mfK1ysRq5wxNtSQoADf1XklMhEUWGUEh/8LnkP/DceWhqPAMGyOY
-----END RSA PRIVATE KEY-----`
)

func TestProxy(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Proxy")
}

func startSSHServer(httpServerURL string) string {
	signer, err := ssh.ParsePrivateKey([]byte(sshPrivateKey))
	if err != nil {
		log.Fatal("Failed to parse private key: ", err)
	}

	config := &ssh.ServerConfig{
		PublicKeyCallback: func(c ssh.ConnMetadata, pubKey ssh.PublicKey) (*ssh.Permissions, error) {
			if string(signer.PublicKey().Marshal()) == string(pubKey.Marshal()) {
				return nil, nil
			}
			return nil, fmt.Errorf("unknown public key for %q", c.User())
		},
	}

	config.AddHostKey(signer)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal("failed to listen for connection: ", err)
	}

	go func() {
		nConn, err := listener.Accept()
		if err != nil {
			log.Fatal("failed to accept incoming connection: ", err)
		}

		_, chans, reqs, err := ssh.NewServerConn(nConn, config)
		if err != nil {
			log.Fatal("failed to handshake: ", err)
		}
		go ssh.DiscardRequests(reqs)

		for newChannel := range chans {
			if newChannel.ChannelType() != "direct-tcpip" {
				newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
				continue
			}
			channel, _, err := newChannel.Accept()
			if err != nil {
				log.Fatalf("Could not accept channel: %v", err)
			}
			defer channel.Close()

			data, err := bufio.NewReader(channel).ReadString('\n')
			if err != nil {
				log.Fatalf("Can't read data from channel: %v", err)
			}

			httpConn, err := net.Dial("tcp", httpServerURL)
			if err != nil {
				log.Fatalf("Could not open connection to http server: %v", err)
			}
			defer httpConn.Close()

			_, err = httpConn.Write([]byte(data + "\r\n\r\n"))
			if err != nil {
				log.Fatalf("Could not write to http server: %v", err)
			}

			data, err = bufio.NewReader(httpConn).ReadString('\n')
			if err != nil {
				log.Fatalf("Can't read data from http conn: %v", err)
			}

			_, err = channel.Write([]byte(data))
			if err != nil {
				log.Fatalf("Can't write data to channel: %v", err)
			}
		}
	}()

	return listener.Addr().String()
}
