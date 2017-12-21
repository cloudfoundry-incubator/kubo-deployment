package ssh_test

import (
	"errors"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
	fakesys "github.com/cloudfoundry/bosh-utils/system/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	boshdir "github.com/cloudfoundry/bosh-cli/director"
	. "github.com/cloudfoundry/bosh-cli/ssh"
)

var _ = Describe("SessionImpl", func() {
	var (
		connOpts       ConnectionOpts
		sessOpts       SessionImplOpts
		result         boshdir.SSHResult
		privKeyFile    *fakesys.FakeFile
		knownHostsFile *fakesys.FakeFile
		fs             *fakesys.FakeFileSystem
		session        *SessionImpl
	)

	BeforeEach(func() {
		connOpts = ConnectionOpts{}
		sessOpts = SessionImplOpts{}
		result = boshdir.SSHResult{}
		fs = fakesys.NewFakeFileSystem()
		privKeyFile = fakesys.NewFakeFile("/tmp/priv-key", fs)
		knownHostsFile = fakesys.NewFakeFile("/tmp/known-hosts", fs)
		fs.ReturnTempFilesByPrefix = map[string]boshsys.File{
			"ssh-priv-key":    privKeyFile,
			"ssh-known-hosts": knownHostsFile,
		}
		session = NewSessionImpl(connOpts, sessOpts, result, fs)
	})

	Describe("Start", func() {
		act := func() *SessionImpl { return NewSessionImpl(connOpts, sessOpts, result, fs) }

		It("writes out private key", func() {
			connOpts.PrivateKey = "priv-key"

			_, err := act().Start()
			Expect(err).ToNot(HaveOccurred())
			Expect(fs.ReadFileString("/tmp/priv-key")).To(Equal("priv-key"))
		})

		It("returns error if cannot create private key temp file", func() {
			fs.TempFileErrorsByPrefix = map[string]error{
				"ssh-priv-key": errors.New("fake-err"),
			}

			_, err := act().Start()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-err"))
		})

		It("returns error if writing public key failed", func() {
			privKeyFile.WriteErr = errors.New("fake-err")

			_, err := act().Start()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-err"))
		})

		It("writes out all known hosts", func() {
			result.Hosts = []boshdir.Host{
				{Host: "127.0.0.1", HostPublicKey: "pub-key1"},
				{Host: "127.0.0.2", HostPublicKey: "pub-key2"},
				{Host: "::1", HostPublicKey: "pub-key3"},
			}

			_, err := act().Start()
			Expect(err).ToNot(HaveOccurred())
			Expect(fs.ReadFileString("/tmp/known-hosts")).To(Equal(
				"127.0.0.1 pub-key1\n127.0.0.2 pub-key2\n::1 pub-key3\n"))
		})

		It("returns error if cannot create known hosts temp file and deletes private key", func() {
			fs.TempFileErrorsByPrefix = map[string]error{
				"ssh-known-hosts": errors.New("fake-err"),
			}

			_, err := act().Start()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-err"))

			Expect(fs.FileExists("/tmp/priv-key")).To(BeFalse())
		})

		It("returns error if writing known hosts failed and deletes private key", func() {
			result.Hosts = []boshdir.Host{
				{Host: "127.0.0.1", HostPublicKey: "pub-key1"},
			}
			knownHostsFile.WriteErr = errors.New("fake-err")

			_, err := act().Start()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-err"))

			Expect(fs.FileExists("/tmp/priv-key")).To(BeFalse())
		})

		It("returns ssh arguments with appropriate configuration", func() {
			result.Hosts = []boshdir.Host{{Host: "127.0.0.1"}} // populate results
			connOpts.PrivateKey = "priv-key"                   // populate connOpts

			args, err := act().Start()
			Expect(err).ToNot(HaveOccurred())
			Expect(args).To(Equal(SSHArgs{
				ConnOpts:       connOpts,
				Result:         result,
				ForceTTY:       false,
				PrivKeyFile:    privKeyFile,
				KnownHostsFile: knownHostsFile,
			}))

			sessOpts.ForceTTY = true

			args, err = act().Start()
			Expect(err).ToNot(HaveOccurred())
			Expect(args).To(Equal(SSHArgs{
				ConnOpts:       connOpts,
				Result:         result,
				ForceTTY:       true,
				PrivKeyFile:    privKeyFile,
				KnownHostsFile: knownHostsFile,
			}))
		})
	})

	Describe("Finish", func() {
		BeforeEach(func() {
			_, err := session.Start()
			Expect(err).ToNot(HaveOccurred())
		})

		It("removes private key and known hosts files", func() {
			err := session.Finish()
			Expect(err).ToNot(HaveOccurred())
			Expect(fs.FileExists("/tmp/priv-key")).To(BeFalse())
			Expect(fs.FileExists("/tmp/known-hosts")).To(BeFalse())
		})

		It("returns error if deleting private key file fails but still deletes known hosts file", func() {
			fs.RemoveAllStub = func(path string) error {
				if path == "/tmp/priv-key" {
					return errors.New("fake-err")
				}
				return nil
			}
			err := session.Finish()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-err"))
			Expect(fs.FileExists("/tmp/known-hosts")).To(BeFalse())
		})

		It("returns error if deleting known hosts file fails but still deletes private key file", func() {
			fs.RemoveAllStub = func(path string) error {
				if path == "/tmp/known-hosts" {
					return errors.New("fake-err")
				}
				return nil
			}
			err := session.Finish()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-err"))
			Expect(fs.FileExists("/tmp/priv-key")).To(BeFalse())
		})
	})
})
