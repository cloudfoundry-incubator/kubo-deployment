package proxy_test

import (
	proxy "github.com/cloudfoundry/socks5-proxy"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/crypto/ssh"
)

var _ = Describe("HostKeyGetter", func() {
	Describe("Get", func() {
		var (
			hostKeyGetter proxy.HostKeyGetter
			key           ssh.PublicKey
			sshServerAddr string
		)

		BeforeEach(func() {
			signer, err := ssh.ParsePrivateKey([]byte(sshPrivateKey))
			Expect(err).NotTo(HaveOccurred())
			key = signer.PublicKey()

			sshServerAddr = startSSHServer("")

			hostKeyGetter = proxy.NewHostKeyGetter()
		})

		It("returns the host key", func() {
			hostKey, err := hostKeyGetter.Get(sshPrivateKey, sshServerAddr)
			Expect(err).NotTo(HaveOccurred())
			Expect(hostKey).To(Equal(key))
		})

		Context("failure cases", func() {
			Context("when parse private key fails", func() {
				It("returns an error", func() {
					_, err := hostKeyGetter.Get("%%%", sshServerAddr)
					Expect(err).To(MatchError("ssh: no key found"))
				})
			})

			Context("when dial fails", func() {
				It("returns an error", func() {
					_, err := hostKeyGetter.Get(sshPrivateKey, "some-bad-url")
					Expect(err).To(MatchError("dial tcp: address some-bad-url: missing port in address"))
				})
			})
		})
	})
})
