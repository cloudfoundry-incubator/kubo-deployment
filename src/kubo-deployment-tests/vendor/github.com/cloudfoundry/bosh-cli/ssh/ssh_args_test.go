package ssh_test

import (
	fakesys "github.com/cloudfoundry/bosh-utils/system/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	boshdir "github.com/cloudfoundry/bosh-cli/director"
	. "github.com/cloudfoundry/bosh-cli/ssh"
)

var _ = Describe("SSHArgs", func() {
	var (
		connOpts       ConnectionOpts
		result         boshdir.SSHResult
		forceTTY       bool
		privKeyFile    *fakesys.FakeFile
		knownHostsFile *fakesys.FakeFile
		fs             *fakesys.FakeFileSystem
		host           boshdir.Host
	)

	BeforeEach(func() {
		connOpts = ConnectionOpts{}
		result = boshdir.SSHResult{}
		forceTTY = false
		fs = fakesys.NewFakeFileSystem()
		privKeyFile = fakesys.NewFakeFile("/tmp/priv-key", fs)
		knownHostsFile = fakesys.NewFakeFile("/tmp/known-hosts", fs)
		host = boshdir.Host{Host: "127.0.0.1", Username: "user"}
	})

	Describe("LoginForHost", func() {
		act := func() []string {
			return SSHArgs{}.LoginForHost(host)
		}

		It("returns login details with IPv4", func() {
			Expect(act()).To(Equal([]string{"127.0.0.1", "-l", "user"}))
		})

		It("returns login details with IPv6 non-bracketed", func() {
			host.Host = "::1"
			Expect(act()).To(Equal([]string{"::1", "-l", "user"}))
		})
	})

	Describe("OptsForHost", func() {
		act := func() []string {
			args := SSHArgs{
				ConnOpts:       connOpts,
				Result:         result,
				ForceTTY:       forceTTY,
				PrivKeyFile:    privKeyFile,
				KnownHostsFile: knownHostsFile,
			}
			return args.OptsForHost(host)
		}

		It("returns ssh options with correct paths to private key and known hosts", func() {
			Expect(act()).To(Equal([]string{
				"-o", "ServerAliveInterval=30",
				"-o", "ForwardAgent=no",
				"-o", "PasswordAuthentication=no",
				"-o", "IdentitiesOnly=yes",
				"-o", "IdentityFile=/tmp/priv-key",
				"-o", "StrictHostKeyChecking=yes",
				"-o", "UserKnownHostsFile=/tmp/known-hosts",
			}))
		})

		It("returns ssh options with forced tty option if requested", func() {
			forceTTY = true

			Expect(act()).To(Equal([]string{
				"-tt",
				"-o", "ServerAliveInterval=30",
				"-o", "ForwardAgent=no",
				"-o", "PasswordAuthentication=no",
				"-o", "IdentitiesOnly=yes",
				"-o", "IdentityFile=/tmp/priv-key",
				"-o", "StrictHostKeyChecking=yes",
				"-o", "UserKnownHostsFile=/tmp/known-hosts",
			}))
		})

		It("returns ssh options with custom raw options specified", func() {
			connOpts.RawOpts = []string{"raw1", "raw2"}

			Expect(act()).To(Equal([]string{
				"-o", "ServerAliveInterval=30",
				"-o", "ForwardAgent=no",
				"-o", "PasswordAuthentication=no",
				"-o", "IdentitiesOnly=yes",
				"-o", "IdentityFile=/tmp/priv-key",
				"-o", "StrictHostKeyChecking=yes",
				"-o", "UserKnownHostsFile=/tmp/known-hosts",
				"raw1", "raw2",
			}))
		})

		It("returns ssh options with gateway settings returned from the Director", func() {
			result.GatewayUsername = "gw-user"
			result.GatewayHost = "gw-host"

			Expect(act()).To(Equal([]string{
				"-o", "ServerAliveInterval=30",
				"-o", "ForwardAgent=no",
				"-o", "PasswordAuthentication=no",
				"-o", "IdentitiesOnly=yes",
				"-o", "IdentityFile=/tmp/priv-key",
				"-o", "StrictHostKeyChecking=yes",
				"-o", "UserKnownHostsFile=/tmp/known-hosts",
				"-o", "ProxyCommand=ssh -tt -W %h:%p -l gw-user gw-host -o ServerAliveInterval=30 -o ForwardAgent=no -o ClearAllForwardings=yes -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null",
			}))
		})

		It("returns ssh options with gateway settings returned from the Director and private key set by user", func() {
			connOpts.GatewayPrivateKeyPath = "/tmp/gw-priv-key"

			result.GatewayUsername = "gw-user"
			result.GatewayHost = "gw-host"

			Expect(act()).To(Equal([]string{
				"-o", "ServerAliveInterval=30",
				"-o", "ForwardAgent=no",
				"-o", "PasswordAuthentication=no",
				"-o", "IdentitiesOnly=yes",
				"-o", "IdentityFile=/tmp/priv-key",
				"-o", "StrictHostKeyChecking=yes",
				"-o", "UserKnownHostsFile=/tmp/known-hosts",
				"-o", "ProxyCommand=ssh -tt -W %h:%p -l gw-user gw-host -o ServerAliveInterval=30 -o ForwardAgent=no -o ClearAllForwardings=yes -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o PasswordAuthentication=no -o IdentitiesOnly=yes -o IdentityFile=/tmp/gw-priv-key",
			}))
		})

		It("returns ssh options with gateway settings overridden by user even if the Director specifies some", func() {
			connOpts.GatewayUsername = "user-gw-user"
			connOpts.GatewayHost = "user-gw-host"

			result.GatewayUsername = "gw-user"
			result.GatewayHost = "gw-host"

			Expect(act()).To(Equal([]string{
				"-o", "ServerAliveInterval=30",
				"-o", "ForwardAgent=no",
				"-o", "PasswordAuthentication=no",
				"-o", "IdentitiesOnly=yes",
				"-o", "IdentityFile=/tmp/priv-key",
				"-o", "StrictHostKeyChecking=yes",
				"-o", "UserKnownHostsFile=/tmp/known-hosts",
				"-o", "ProxyCommand=ssh -tt -W %h:%p -l user-gw-user user-gw-host -o ServerAliveInterval=30 -o ForwardAgent=no -o ClearAllForwardings=yes -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null",
			}))
		})

		It("returns ssh options without gateway settings if disabled even if user or the Director specifies some", func() {
			connOpts.GatewayDisable = true
			connOpts.GatewayUsername = "user-gw-user"
			connOpts.GatewayHost = "user-gw-host"

			result.GatewayUsername = "gw-user"
			result.GatewayHost = "gw-host"

			Expect(act()).To(Equal([]string{
				"-o", "ServerAliveInterval=30",
				"-o", "ForwardAgent=no",
				"-o", "PasswordAuthentication=no",
				"-o", "IdentitiesOnly=yes",
				"-o", "IdentityFile=/tmp/priv-key",
				"-o", "StrictHostKeyChecking=yes",
				"-o", "UserKnownHostsFile=/tmp/known-hosts",
			}))
		})

		It("returns ssh options without socks5 settings if SOCKS5Proxy is set", func() {
			connOpts.GatewayDisable = true
			connOpts.GatewayUsername = "user-gw-user"
			connOpts.GatewayHost = "user-gw-host"
			connOpts.SOCKS5Proxy = "socks5://some-proxy"

			result.GatewayUsername = "gw-user"
			result.GatewayHost = "gw-host"

			Expect(act()).To(Equal([]string{
				"-o", "ServerAliveInterval=30",
				"-o", "ForwardAgent=no",
				"-o", "PasswordAuthentication=no",
				"-o", "IdentitiesOnly=yes",
				"-o", "IdentityFile=/tmp/priv-key",
				"-o", "StrictHostKeyChecking=yes",
				"-o", "UserKnownHostsFile=/tmp/known-hosts",
				"-o", "ProxyCommand=nc -x some-proxy %h %p",
			}))
		})

		It("returns ssh options with bracketed gateway proxy command if host IP is IPv6", func() {
			result.GatewayUsername = "gw-user"
			result.GatewayHost = "gw-host"
			host = boshdir.Host{Host: "::1", Username: "user"}

			Expect(act()).To(Equal([]string{
				"-o", "ServerAliveInterval=30",
				"-o", "ForwardAgent=no",
				"-o", "PasswordAuthentication=no",
				"-o", "IdentitiesOnly=yes",
				"-o", "IdentityFile=/tmp/priv-key",
				"-o", "StrictHostKeyChecking=yes",
				"-o", "UserKnownHostsFile=/tmp/known-hosts",
				"-o", "ProxyCommand=ssh -tt -W [%h]:%p -l gw-user gw-host -o ServerAliveInterval=30 -o ForwardAgent=no -o ClearAllForwardings=yes -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null",
			}))
		})

		It("returns ssh options with non-bracketed IPs if gateway IP is IPv6", func() {
			result.GatewayUsername = "gw-user"
			result.GatewayHost = "::1"

			Expect(act()).To(Equal([]string{
				"-o", "ServerAliveInterval=30",
				"-o", "ForwardAgent=no",
				"-o", "PasswordAuthentication=no",
				"-o", "IdentitiesOnly=yes",
				"-o", "IdentityFile=/tmp/priv-key",
				"-o", "StrictHostKeyChecking=yes",
				"-o", "UserKnownHostsFile=/tmp/known-hosts",
				"-o", "ProxyCommand=ssh -tt -W %h:%p -l gw-user ::1 -o ServerAliveInterval=30 -o ForwardAgent=no -o ClearAllForwardings=yes -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null",
			}))
		})

		It("returns ssh options with bracketed and non-bracketed IPs if host and gateway IP is IPv6", func() {
			result.GatewayUsername = "gw-user"
			result.GatewayHost = "::1"
			host = boshdir.Host{Host: "::2", Username: "user"}

			Expect(act()).To(Equal([]string{
				"-o", "ServerAliveInterval=30",
				"-o", "ForwardAgent=no",
				"-o", "PasswordAuthentication=no",
				"-o", "IdentitiesOnly=yes",
				"-o", "IdentityFile=/tmp/priv-key",
				"-o", "StrictHostKeyChecking=yes",
				"-o", "UserKnownHostsFile=/tmp/known-hosts",
				"-o", "ProxyCommand=ssh -tt -W [%h]:%p -l gw-user ::1 -o ServerAliveInterval=30 -o ForwardAgent=no -o ClearAllForwardings=yes -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null",
			}))
		})

		It("returns ssh options non-bracketed if host is IPv6 and SOCKS5Proxy is set", func() {
			host = boshdir.Host{Host: "::1", Username: "user"}
			connOpts.SOCKS5Proxy = "socks5://some-proxy"

			Expect(act()).To(Equal([]string{
				"-o", "ServerAliveInterval=30",
				"-o", "ForwardAgent=no",
				"-o", "PasswordAuthentication=no",
				"-o", "IdentitiesOnly=yes",
				"-o", "IdentityFile=/tmp/priv-key",
				"-o", "StrictHostKeyChecking=yes",
				"-o", "UserKnownHostsFile=/tmp/known-hosts",
				"-o", "ProxyCommand=nc -x some-proxy %h %p",
			}))
		})
	})
})
