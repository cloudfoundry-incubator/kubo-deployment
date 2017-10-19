package cmd_test

import (
	"errors"

	semver "github.com/cppforlife/go-semi-semantic/version"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry/bosh-cli/cmd"
	fakecmd "github.com/cloudfoundry/bosh-cli/cmd/cmdfakes"
	boshdir "github.com/cloudfoundry/bosh-cli/director"
	boshrel "github.com/cloudfoundry/bosh-cli/release"
	fakerel "github.com/cloudfoundry/bosh-cli/release/releasefakes"
)

var _ = Describe("ReleaseManager", func() {
	var (
		createReleaseCmd *fakecmd.FakeReleaseCreatingCmd
		uploadReleaseCmd *fakecmd.FakeReleaseUploadingCmd
		releaseManager   ReleaseManager
	)

	BeforeEach(func() {
		createReleaseCmd = &fakecmd.FakeReleaseCreatingCmd{
			RunStub: func(opts CreateReleaseOpts) (boshrel.Release, error) {
				release := &fakerel.FakeRelease{
					NameStub:    func() string { return opts.Name },
					VersionStub: func() string { return opts.Name + "-created-ver" },
				}
				return release, nil
			},
		}

		uploadReleaseCmd = &fakecmd.FakeReleaseUploadingCmd{}

		threadCount := 5
		releaseManager = NewReleaseManager(createReleaseCmd, uploadReleaseCmd, threadCount)
	})

	Describe("UploadReleases", func() {
		It("uploads remote releases skipping releases without url", func() {
			bytes := []byte(`
releases:
- name: capi
  sha1: capi-sha1
  url: https://capi-url
  version: 1+capi
- name: rel-without-upload
  version: 1+rel
- name: consul
  sha1: consul-sha1
  url: https://consul-url
  version: 1+consul
- name: compiled-release
  url: https://compiled-release-url
  sha1: compiled-release-sha1
  version: 1+compiled-release
  stemcell:
    os: ubuntu-trusty
    version: 3421
- name: local
  url: file:///local-dir
  version: create
`)

			_, err := releaseManager.UploadReleases(bytes)
			Expect(err).ToNot(HaveOccurred())

			Expect(uploadReleaseCmd.RunCallCount()).To(Equal(4))
			runArgs := []UploadReleaseOpts{
				uploadReleaseCmd.RunArgsForCall(0),
				uploadReleaseCmd.RunArgsForCall(1),
				uploadReleaseCmd.RunArgsForCall(2),
				uploadReleaseCmd.RunArgsForCall(3),
			}

			var capiRelease UploadReleaseOpts
			var consulRelease UploadReleaseOpts
			var compiledRelease UploadReleaseOpts
			var localRelease UploadReleaseOpts
			for _, opts := range runArgs {
				switch opts.Name {
				case "capi":
					capiRelease = opts
				case "consul":
					consulRelease = opts
				case "compiled-release":
					compiledRelease = opts
				case "local":
					localRelease = opts
				}
			}

			Expect(capiRelease).To(Equal(UploadReleaseOpts{
				Name:    "capi",
				Args:    UploadReleaseArgs{URL: URLArg("https://capi-url")},
				SHA1:    "capi-sha1",
				Version: VersionArg(semver.MustNewVersionFromString("1+capi")),
			}))
			Expect(consulRelease).To(Equal(UploadReleaseOpts{
				Name:    "consul",
				Args:    UploadReleaseArgs{URL: URLArg("https://consul-url")},
				SHA1:    "consul-sha1",
				Version: VersionArg(semver.MustNewVersionFromString("1+consul")),
			}))
			Expect(compiledRelease).To(Equal(UploadReleaseOpts{
				Name:    "compiled-release",
				Args:    UploadReleaseArgs{URL: URLArg("https://compiled-release-url")},
				SHA1:    "compiled-release-sha1",
				Version: VersionArg(semver.MustNewVersionFromString("1+compiled-release")),

				Stemcell: boshdir.NewOSVersionSlug("ubuntu-trusty", "3421"),
			}))
			Expect(localRelease).To(Equal(UploadReleaseOpts{
				Release: localRelease.Release, // only Release should be set
			}))
		})

		It("skips uploading releases if url is not provided, even if the version is invalid", func() {
			bytes := []byte(`
releases:
- name: capi
  version: ((/blah_interpolate_me_with_config_server))
`)

			_, err := releaseManager.UploadReleases(bytes)
			Expect(err).ToNot(HaveOccurred())
			Expect(uploadReleaseCmd.RunCallCount()).To(Equal(0))
		})

		It("creates releases if version is 'create' skipping others", func() {
			bytes := []byte(`
releases:
- name: capi
  url: file:///capi-dir
  version: create
- name: rel-without-upload
  version: 1+rel
- name: consul
  url: /consul-dir # doesn't require file://
  version: create
`)

			bytes, err := releaseManager.UploadReleases(bytes)
			Expect(err).ToNot(HaveOccurred())

			Expect(createReleaseCmd.RunCallCount()).To(Equal(2))
			runArgs := []CreateReleaseOpts{
				createReleaseCmd.RunArgsForCall(0),
				createReleaseCmd.RunArgsForCall(1),
			}

			var capiRelease CreateReleaseOpts
			var consulRelease CreateReleaseOpts
			for _, opts := range runArgs {
				switch opts.Name {
				case "capi":
					capiRelease = opts
				case "consul":
					consulRelease = opts
				}
			}

			Expect(capiRelease).To(Equal(CreateReleaseOpts{
				Name:             "capi",
				Directory:        DirOrCWDArg{Path: "/capi-dir"},
				TimestampVersion: true,
				Force:            true,
			}))

			Expect(consulRelease).To(Equal(CreateReleaseOpts{
				Name:             "consul",
				Directory:        DirOrCWDArg{Path: "/consul-dir"},
				TimestampVersion: true,
				Force:            true,
			}))

			Expect(bytes).To(Equal([]byte(`releases:
- name: capi
  version: capi-created-ver
- name: rel-without-upload
  version: 1+rel
- name: consul
  version: consul-created-ver
`)))
		})

		It("returns error and does not upload if creating release fails", func() {
			bytes := []byte(`
releases:
- name: capi
  url: file:///capi-dir
  version: create
`)
			createReleaseCmd.RunReturns(nil, errors.New("fake-err"))

			_, err := releaseManager.UploadReleases(bytes)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-err"))

			Expect(uploadReleaseCmd.RunCallCount()).To(Equal(0))
		})

		It("returns error if uploading release fails", func() {
			bytes := []byte(`
releases:
- name: capi
  sha1: capi-sha1
  url: https://capi-url
  version: 1+capi
`)
			uploadReleaseCmd.RunReturns(errors.New("fake-err"))

			_, err := releaseManager.UploadReleases(bytes)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-err"))
		})

		It("returns an error and does not upload if release version cannot be parsed", func() {
			bytes := []byte(`
releases:
- name: capi
  sha1: capi-sha1
  url: https://capi-url
  version: 1+capi+capi
`)

			_, err := releaseManager.UploadReleases(bytes)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Expected version '1+capi+capi' to match version format"))

			Expect(uploadReleaseCmd.RunCallCount()).To(Equal(0))
		})

		It("returns an error if bytes cannot be parsed to find releases", func() {
			bytes := []byte(`-`)

			_, err := releaseManager.UploadReleases(bytes)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Parsing manifest"))

			Expect(createReleaseCmd.RunCallCount()).To(Equal(0))
			Expect(uploadReleaseCmd.RunCallCount()).To(Equal(0))
		})
	})
})
