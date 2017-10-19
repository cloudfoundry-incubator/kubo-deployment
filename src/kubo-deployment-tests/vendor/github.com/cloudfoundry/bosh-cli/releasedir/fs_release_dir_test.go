package releasedir_test

import (
	"errors"
	"time"

	"code.cloudfoundry.org/clock"
	"code.cloudfoundry.org/clock/fakeclock"
	fakesys "github.com/cloudfoundry/bosh-utils/system/fakes"
	semver "github.com/cppforlife/go-semi-semantic/version"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	boshrel "github.com/cloudfoundry/bosh-cli/release"
	boshman "github.com/cloudfoundry/bosh-cli/release/manifest"
	boshpkg "github.com/cloudfoundry/bosh-cli/release/pkg"
	fakerel "github.com/cloudfoundry/bosh-cli/release/releasefakes"
	fakeres "github.com/cloudfoundry/bosh-cli/release/resource/resourcefakes"
	. "github.com/cloudfoundry/bosh-cli/releasedir"
	fakereldir "github.com/cloudfoundry/bosh-cli/releasedir/releasedirfakes"
)

var _ = Describe("FSGenerator", func() {
	var (
		config        *fakereldir.FakeConfig
		gitRepo       *fakereldir.FakeGitRepo
		blobsDir      *fakereldir.FakeBlobsDir
		gen           *fakereldir.FakeGenerator
		devReleases   *fakereldir.FakeReleaseIndex
		finalReleases *fakereldir.FakeReleaseIndex
		finalIndicies boshrel.ArchiveIndicies
		reader        *fakerel.FakeReader
		timeService   clock.Clock
		fs            *fakesys.FakeFileSystem
		releaseDir    FSReleaseDir
	)

	BeforeEach(func() {
		config = &fakereldir.FakeConfig{}
		gitRepo = &fakereldir.FakeGitRepo{}
		blobsDir = &fakereldir.FakeBlobsDir{}
		gen = &fakereldir.FakeGenerator{}
		devReleases = &fakereldir.FakeReleaseIndex{}
		finalReleases = &fakereldir.FakeReleaseIndex{}
		finalIndicies = boshrel.ArchiveIndicies{
			Jobs:     &fakeres.FakeArchiveIndex{},
			Packages: &fakeres.FakeArchiveIndex{},
		}
		reader = &fakerel.FakeReader{}
		timeService = fakeclock.NewFakeClock(time.Date(2009, time.November, 10, 23, 1, 2, 333, time.UTC))
		fs = fakesys.NewFakeFileSystem()
		releaseDir = NewFSReleaseDir(
			"/dir",
			config,
			gitRepo,
			blobsDir,
			gen,
			devReleases,
			finalReleases,
			finalIndicies,
			reader,
			timeService,
			fs,
			2,
		)
	})

	Describe("Init", func() {
		It("creates commont jobs, packages and src directories", func() {
			err := releaseDir.Init(true)
			Expect(err).ToNot(HaveOccurred())

			Expect(fs.FileExists("/dir/jobs")).To(BeTrue())
			Expect(fs.FileExists("/dir/packages")).To(BeTrue())
			Expect(fs.FileExists("/dir/src")).To(BeTrue())
		})

		It("returns error if creating common dirs fails", func() {
			fs.MkdirAllError = errors.New("fake-err")

			err := releaseDir.Init(true)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-err"))
		})

		It("saves release name to directory base name", func() {
			err := releaseDir.Init(true)
			Expect(err).ToNot(HaveOccurred())

			Expect(config.SaveNameCallCount()).To(Equal(1))
			Expect(config.SaveNameArgsForCall(0)).To(Equal("dir"))
		})

		It("saves release name to directory base name stripping '-release' suffix from the name", func() {
			releaseDir := NewFSReleaseDir(
				"/dir-release",
				config,
				gitRepo,
				blobsDir,
				gen,
				devReleases,
				finalReleases,
				finalIndicies,
				reader,
				timeService,
				fs,
				2,
			)

			err := releaseDir.Init(true)
			Expect(err).ToNot(HaveOccurred())

			Expect(config.SaveNameCallCount()).To(Equal(1))
			Expect(config.SaveNameArgsForCall(0)).To(Equal("dir"))
		})

		It("returns error if saving final name fails", func() {
			config.SaveNameReturns(errors.New("fake-err"))

			err := releaseDir.Init(true)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-err"))
		})

		It("inits blobs", func() {
			err := releaseDir.Init(true)
			Expect(err).ToNot(HaveOccurred())

			Expect(blobsDir.InitCallCount()).To(Equal(1))
		})

		It("returns error if initing blobs fails", func() {
			blobsDir.InitReturns(errors.New("fake-err"))

			err := releaseDir.Init(true)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-err"))
		})

		It("inits git repo if requested", func() {
			err := releaseDir.Init(true)
			Expect(err).ToNot(HaveOccurred())

			Expect(gitRepo.InitCallCount()).To(Equal(1))
		})

		It("does not init git repo if not requested", func() {
			err := releaseDir.Init(false)
			Expect(err).ToNot(HaveOccurred())

			Expect(gitRepo.InitCallCount()).To(Equal(0))
		})

		It("returns error if initing git repo fails", func() {
			gitRepo.InitReturns(errors.New("fake-err"))

			err := releaseDir.Init(true)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-err"))
		})
	})

	Describe("GenerateJob", func() {
		It("delegates to generator", func() {
			gen.GenerateJobStub = func(name string) error {
				Expect(name).To(Equal("job1"))
				return errors.New("fake-err")
			}
			Expect(releaseDir.GenerateJob("job1")).To(Equal(errors.New("fake-err")))
		})
	})

	Describe("GeneratePackage", func() {
		It("delegates to generator", func() {
			gen.GeneratePackageStub = func(name string) error {
				Expect(name).To(Equal("job1"))
				return errors.New("fake-err")
			}
			Expect(releaseDir.GeneratePackage("job1")).To(Equal(errors.New("fake-err")))
		})
	})

	Describe("Reset", func() {
		It("removes .blobs, blobs, .dev_builds and dev_releases", func() {
			fs.WriteFileString("/dir/.dev_builds/sub-dir", "")
			fs.WriteFileString("/dir/dev_releases/sub-dir", "")
			fs.WriteFileString("/dir/.blobs/sub-dir", "")
			fs.WriteFileString("/dir/blobs/sub-dir", "")

			err := releaseDir.Reset()
			Expect(err).ToNot(HaveOccurred())

			Expect(fs.FileExists("/dir/.dev_builds")).To(BeFalse())
			Expect(fs.FileExists("/dir/dev_releases")).To(BeFalse())
			Expect(fs.FileExists("/dir/.blobs")).To(BeFalse())
			Expect(fs.FileExists("/dir/blobs")).To(BeFalse())
		})

		It("returns error when deleting directory fails", func() {
			fs.RemoveAllStub = func(_ string) error { return errors.New("fake-err") }

			err := releaseDir.Reset()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-err"))
		})
	})

	Describe("DefaultName", func() {
		It("delegates to config", func() {
			config.NameReturns("name", errors.New("fake-err"))

			name, err := releaseDir.DefaultName()
			Expect(name).To(Equal("name"))
			Expect(err).To(Equal(errors.New("fake-err")))
		})
	})

	Describe("NextFinalVersion", func() {
		It("returns incremented last final version for specific release name", func() {
			finalReleases.LastVersionStub = func(name string) (*semver.Version, error) {
				Expect(name).To(Equal("rel1"))
				lastVer := semver.MustNewVersionFromString("1.1")
				return &lastVer, nil
			}

			ver, err := releaseDir.NextFinalVersion("rel1")
			Expect(err).ToNot(HaveOccurred())
			Expect(ver.String()).To(Equal(semver.MustNewVersionFromString("1.2").String()))
		})

		It("returns '1' if there are no versions so that when it's finalized it will be '1'", func() {
			finalReleases.LastVersionReturns(nil, nil)

			ver, err := releaseDir.NextFinalVersion("rel1")
			Expect(err).ToNot(HaveOccurred())
			Expect(ver.String()).To(Equal(semver.MustNewVersionFromString("1").String()))
		})

		It("returns error if cannot find out last version", func() {
			finalReleases.LastVersionReturns(nil, errors.New("fake-err"))

			_, err := releaseDir.NextFinalVersion("rel1")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-err"))
		})

		It("returns error if incrementing fails", func() {
			lastVer := semver.MustNewVersionFromString("a")
			finalReleases.LastVersionReturns(&lastVer, nil)

			_, err := releaseDir.NextFinalVersion("rel1")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Incrementing last final version"))
		})
	})

	Describe("NextDevVersion", func() {
		It("returns incremented last final version for specific release name", func() {
			finalReleases.LastVersionStub = func(name string) (*semver.Version, error) {
				Expect(name).To(Equal("rel1"))
				lastVer := semver.MustNewVersionFromString("1.1")
				return &lastVer, nil
			}

			ver, err := releaseDir.NextDevVersion("rel1", false)
			Expect(err).ToNot(HaveOccurred())
			Expect(ver.String()).To(Equal(semver.MustNewVersionFromString("1.1+dev.1").String()))
		})

		It("returns incremented last dev version for specific release name", func() {
			devReleases.LastVersionStub = func(name string) (*semver.Version, error) {
				Expect(name).To(Equal("rel1"))
				lastVer := semver.MustNewVersionFromString("1.1+dev.1")
				return &lastVer, nil
			}

			ver, err := releaseDir.NextDevVersion("rel1", false)
			Expect(err).ToNot(HaveOccurred())
			Expect(ver.String()).To(Equal(semver.MustNewVersionFromString("1.1+dev.2").String()))
		})

		It("returns timestamp-ed dev version for specific release name", func() {
			devReleases.LastVersionStub = func(name string) (*semver.Version, error) {
				Expect(name).To(Equal("rel1"))
				lastVer := semver.MustNewVersionFromString("1.1+dev.1")
				return &lastVer, nil
			}

			ver, err := releaseDir.NextDevVersion("rel1", true)
			Expect(err).ToNot(HaveOccurred())
			Expect(ver.String()).To(Equal(semver.MustNewVersionFromString("1.1+dev.1257894062").String()))
		})

		It("returns incremented greater dev version compared to final version for specific release name", func() {
			finalReleases.LastVersionStub = func(name string) (*semver.Version, error) {
				Expect(name).To(Equal("rel1"))
				lastVer := semver.MustNewVersionFromString("1.1")
				return &lastVer, nil
			}

			devReleases.LastVersionStub = func(name string) (*semver.Version, error) {
				Expect(name).To(Equal("rel1"))
				lastVer := semver.MustNewVersionFromString("1.1+dev.1")
				return &lastVer, nil
			}

			ver, err := releaseDir.NextDevVersion("rel1", false)
			Expect(err).ToNot(HaveOccurred())
			Expect(ver.String()).To(Equal(semver.MustNewVersionFromString("1.1+dev.2").String()))
		})

		It("returns incremented greater final version compared to dev version for specific release name", func() {
			finalReleases.LastVersionStub = func(name string) (*semver.Version, error) {
				Expect(name).To(Equal("rel1"))
				lastVer := semver.MustNewVersionFromString("1.2")
				return &lastVer, nil
			}

			devReleases.LastVersionStub = func(name string) (*semver.Version, error) {
				Expect(name).To(Equal("rel1"))
				lastVer := semver.MustNewVersionFromString("1.1+dev.1")
				return &lastVer, nil
			}

			ver, err := releaseDir.NextDevVersion("rel1", false)
			Expect(err).ToNot(HaveOccurred())
			Expect(ver.String()).To(Equal(semver.MustNewVersionFromString("1.2+dev.1").String()))
		})

		It("returns '0+dev.1' if there are no dev or final versions", func() {
			ver, err := releaseDir.NextDevVersion("rel1", false)
			Expect(err).ToNot(HaveOccurred())
			Expect(ver.String()).To(Equal(semver.MustNewVersionFromString("0+dev.1").String()))
		})

		It("returns first timestamp-ed dev version if there are no dev or final versions", func() {
			ver, err := releaseDir.NextDevVersion("rel1", true)
			Expect(err).ToNot(HaveOccurred())
			Expect(ver.String()).To(Equal(semver.MustNewVersionFromString("0+dev.1257894062").String()))
		})

		It("returns error if cannot find out last dev version", func() {
			devReleases.LastVersionReturns(nil, errors.New("fake-err"))

			_, err := releaseDir.NextDevVersion("rel1", false)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-err"))
		})

		It("returns error if cannot find out last final version", func() {
			finalReleases.LastVersionReturns(nil, errors.New("fake-err"))

			_, err := releaseDir.NextDevVersion("rel1", false)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-err"))
		})

		It("returns error if incrementing fails", func() {
			lastVer := semver.MustNewVersionFromString("1+a")
			finalReleases.LastVersionReturns(&lastVer, nil)

			_, err := releaseDir.NextDevVersion("rel1", false)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Incrementing last dev version"))
		})
	})

	Describe("FindRelease", func() {
		var (
			expectedRelease *fakerel.FakeRelease
		)

		BeforeEach(func() {
			expectedRelease = &fakerel.FakeRelease{
				NameStub: func() string { return "rel1" },
			}
		})

		Context("when name and version are not specified", func() {
			BeforeEach(func() {
				config.NameReturns("rel1", nil)
			})

			It("returns last final release for specific release name", func() {
				finalReleases.LastVersionStub = func(name string) (*semver.Version, error) {
					Expect(name).To(Equal("rel1"))
					lastVer := semver.MustNewVersionFromString("1.1")
					return &lastVer, nil
				}

				finalReleases.ManifestPathStub = func(name, ver string) string {
					Expect(name).To(Equal("rel1"))
					Expect(ver).To(Equal("1.1"))
					return "manifest-path"
				}

				reader.ReadStub = func(path string) (boshrel.Release, error) {
					Expect(path).To(Equal("manifest-path"))
					return expectedRelease, nil
				}

				release, err := releaseDir.FindRelease("", semver.Version{})
				Expect(err).ToNot(HaveOccurred())
				Expect(release).To(Equal(expectedRelease))
			})

			It("returns last dev release for specific release name", func() {
				devReleases.LastVersionStub = func(name string) (*semver.Version, error) {
					Expect(name).To(Equal("rel1"))
					lastVer := semver.MustNewVersionFromString("1.1+dev.1")
					return &lastVer, nil
				}

				devReleases.ManifestPathStub = func(name, ver string) string {
					Expect(name).To(Equal("rel1"))
					Expect(ver).To(Equal("1.1+dev.1"))
					return "manifest-path"
				}

				reader.ReadStub = func(path string) (boshrel.Release, error) {
					Expect(path).To(Equal("manifest-path"))
					return expectedRelease, nil
				}

				release, err := releaseDir.FindRelease("", semver.Version{})
				Expect(err).ToNot(HaveOccurred())
				Expect(release).To(Equal(expectedRelease))
			})

			It("returns greater dev release compared to final release for specific release name", func() {
				finalReleases.LastVersionStub = func(name string) (*semver.Version, error) {
					Expect(name).To(Equal("rel1"))
					lastVer := semver.MustNewVersionFromString("1.1")
					return &lastVer, nil
				}

				devReleases.LastVersionStub = func(name string) (*semver.Version, error) {
					Expect(name).To(Equal("rel1"))
					lastVer := semver.MustNewVersionFromString("1.1+dev.1")
					return &lastVer, nil
				}

				devReleases.ManifestPathStub = func(name, ver string) string {
					Expect(name).To(Equal("rel1"))
					Expect(ver).To(Equal("1.1+dev.1"))
					return "manifest-path"
				}

				reader.ReadStub = func(path string) (boshrel.Release, error) {
					Expect(path).To(Equal("manifest-path"))
					return expectedRelease, nil
				}

				release, err := releaseDir.FindRelease("", semver.Version{})
				Expect(err).ToNot(HaveOccurred())
				Expect(release).To(Equal(expectedRelease))
			})

			It("returns greater final release compared to dev release for specific release name", func() {
				finalReleases.LastVersionStub = func(name string) (*semver.Version, error) {
					Expect(name).To(Equal("rel1"))
					lastVer := semver.MustNewVersionFromString("1.2")
					return &lastVer, nil
				}

				devReleases.LastVersionStub = func(name string) (*semver.Version, error) {
					Expect(name).To(Equal("rel1"))
					lastVer := semver.MustNewVersionFromString("1.1+dev.1")
					return &lastVer, nil
				}

				finalReleases.ManifestPathStub = func(name, ver string) string {
					Expect(name).To(Equal("rel1"))
					Expect(ver).To(Equal("1.2"))
					return "manifest-path"
				}

				reader.ReadStub = func(path string) (boshrel.Release, error) {
					Expect(path).To(Equal("manifest-path"))
					return expectedRelease, nil
				}

				release, err := releaseDir.FindRelease("", semver.Version{})
				Expect(err).ToNot(HaveOccurred())
				Expect(release).To(Equal(expectedRelease))
			})

			It("returns error if there are no dev or final versions", func() {
				_, err := releaseDir.FindRelease("", semver.Version{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Expected to find at least one dev or final version"))
			})

			It("returns error if cannot find out last dev version", func() {
				devReleases.LastVersionReturns(nil, errors.New("fake-err"))

				_, err := releaseDir.FindRelease("", semver.Version{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-err"))
			})

			It("returns error if cannot find out last final version", func() {
				finalReleases.LastVersionReturns(nil, errors.New("fake-err"))

				_, err := releaseDir.FindRelease("", semver.Version{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-err"))
			})

			It("retuns error if cannot determine final name", func() {
				config.NameReturns("", errors.New("fake-err"))

				_, err := releaseDir.FindRelease("", semver.Version{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-err"))
			})
		})

		Context("when name and version is specified", func() {
			It("returns final release for specific release name and version", func() {
				finalReleases.ManifestPathStub = func(name, ver string) string {
					Expect(name).To(Equal("rel1"))
					Expect(ver).To(Equal("1.1"))
					return "manifest-path"
				}

				reader.ReadStub = func(path string) (boshrel.Release, error) {
					Expect(path).To(Equal("manifest-path"))
					return expectedRelease, nil
				}

				release, err := releaseDir.FindRelease("rel1", semver.MustNewVersionFromString("1.1"))
				Expect(err).ToNot(HaveOccurred())
				Expect(release).To(Equal(expectedRelease))
			})
		})
	})

	Describe("BuildRelease", func() {
		var (
			ver             semver.Version
			expectedRelease *fakerel.FakeRelease
		)

		BeforeEach(func() {
			ver = semver.MustNewVersionFromString("1.1")

			expectedRelease = &fakerel.FakeRelease{
				NameStub: func() string { return "rel1" },
				ManifestStub: func() boshman.Manifest {
					return boshman.Manifest{Name: "rel1"}
				},
			}
		})

		It("builds release", func() {
			var ops []string

			gitRepo.MustNotBeDirtyStub = func(force bool) (bool, error) {
				ops = append(ops, "dirty")
				return true, nil
			}

			gitRepo.LastCommitSHAReturns("commit", nil)

			blobsDir.SyncBlobsStub = func(numOfParallelWorkers int) error {
				ops = append(ops, "blobs")
				return nil
			}

			reader.ReadStub = func(path string) (boshrel.Release, error) {
				Expect(path).To(Equal("/dir"))
				ops = append(ops, "read")
				return expectedRelease, nil
			}

			devReleases.AddStub = func(manifest boshman.Manifest) error {
				Expect(manifest).To(Equal(boshman.Manifest{Name: "rel1"}))
				ops = append(ops, "manifest")
				return nil
			}

			release, err := releaseDir.BuildRelease("rel1", ver, false)
			Expect(err).ToNot(HaveOccurred())
			Expect(release).To(Equal(expectedRelease))

			Expect(expectedRelease.SetNameArgsForCall(0)).To(Equal("rel1"))
			Expect(expectedRelease.SetVersionArgsForCall(0)).To(Equal("1.1"))
			Expect(expectedRelease.SetCommitHashArgsForCall(0)).To(Equal("commit"))
			Expect(expectedRelease.SetUncommittedChangesArgsForCall(0)).To(BeTrue())

			Expect(ops).To(Equal([]string{"dirty", "blobs", "read", "manifest"}))
		})

		It("returns error if git is dirty and force is not set", func() {
			gitRepo.MustNotBeDirtyReturns(true, errors.New("dirty"))

			_, err := releaseDir.BuildRelease("rel1", ver, false)
			Expect(err).To(Equal(errors.New("dirty")))
		})

		It("returns error if last commit cannot be retrieved", func() {
			gitRepo.LastCommitSHAReturns("", errors.New("fake-err"))

			_, err := releaseDir.BuildRelease("rel1", ver, false)
			Expect(err).To(Equal(errors.New("fake-err")))
		})

		It("returns error if reading release", func() {
			reader.ReadReturns(nil, errors.New("fake-err"))

			_, err := releaseDir.BuildRelease("rel1", ver, false)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-err"))
		})

		It("returns error if adding dev release fails", func() {
			reader.ReadReturns(expectedRelease, nil)
			devReleases.AddReturns(errors.New("fake-err"))

			_, err := releaseDir.BuildRelease("rel1", ver, false)
			Expect(err).To(Equal(errors.New("fake-err")))
		})
	})

	Describe("VendorPackage", func() {
		It("finalizes given package and its dependencies and records vendored package manifest locks", func() {
			pkg1Res := &fakeres.FakeResource{
				NameStub:        func() string { return "pkg1-name" },
				FingerprintStub: func() string { return "pkg1-fp" },
			}
			pkg2Res := &fakeres.FakeResource{
				NameStub:        func() string { return "pkg2-name" },
				FingerprintStub: func() string { return "pkg2-fp" },
			}
			pkg3Res := &fakeres.FakeResource{
				NameStub:        func() string { return "pkg3-name" },
				FingerprintStub: func() string { return "pkg3-fp" },
			}
			pkg4Res := &fakeres.FakeResource{
				NameStub:        func() string { return "pkg4-name" },
				FingerprintStub: func() string { return "pkg4-fp" },
			}

			pkg1 := boshpkg.NewPackage(pkg1Res, []string{"pkg2-name", "pkg3-name"})
			pkg2 := boshpkg.NewPackage(pkg2Res, []string{"pkg4-name"})
			pkg3 := boshpkg.NewPackage(pkg3Res, []string{"pkg4-name"})
			pkg4 := boshpkg.NewPackage(pkg4Res, nil)
			pkg1.AttachDependencies([]*boshpkg.Package{pkg2, pkg3})
			pkg2.AttachDependencies([]*boshpkg.Package{pkg4})
			pkg3.AttachDependencies([]*boshpkg.Package{pkg4})

			// previous content will be overwritten
			Expect(fs.WriteFileString("/dir/packages/pkg1-name/spec", "old-spec")).ToNot(HaveOccurred())
			Expect(fs.WriteFileString("/dir/packages/pkg1-name/packaging", "old-packaging")).ToNot(HaveOccurred())
			Expect(fs.WriteFileString("/dir/packages/pkg2-name/spec.lock", "old-spec-lock")).ToNot(HaveOccurred())

			err := releaseDir.VendorPackage(pkg1)
			Expect(err).ToNot(HaveOccurred())

			// recorded files
			Expect(fs.FileExists("/dir/packages/pkg1-name/spec")).To(BeFalse())
			Expect(fs.FileExists("/dir/packages/pkg1-name/packaging")).To(BeFalse())
			Expect(fs.ReadFileString("/dir/packages/pkg1-name/spec.lock")).To(Equal(`name: pkg1-name
fingerprint: pkg1-fp
dependencies:
- pkg2-name
- pkg3-name
`))

			Expect(fs.ReadFileString("/dir/packages/pkg2-name/spec.lock")).To(Equal(`name: pkg2-name
fingerprint: pkg2-fp
dependencies:
- pkg4-name
`))

			Expect(fs.ReadFileString("/dir/packages/pkg3-name/spec.lock")).To(Equal(`name: pkg3-name
fingerprint: pkg3-fp
dependencies:
- pkg4-name
`))

			Expect(fs.ReadFileString("/dir/packages/pkg4-name/spec.lock")).To(Equal(`name: pkg4-name
fingerprint: pkg4-fp
`))

			// Use == for pointer equality
			Expect(finalIndicies.Packages).ToNot(BeNil())

			Expect(pkg1Res.FinalizeCallCount()).To(Equal(1))
			Expect(pkg1Res.FinalizeArgsForCall(0) == finalIndicies.Packages).To(BeTrue())

			Expect(pkg2Res.FinalizeCallCount()).To(Equal(1))
			Expect(pkg2Res.FinalizeArgsForCall(0) == finalIndicies.Packages).To(BeTrue())

			Expect(pkg3Res.FinalizeCallCount()).To(Equal(1))
			Expect(pkg3Res.FinalizeArgsForCall(0) == finalIndicies.Packages).To(BeTrue())

			Expect(pkg4Res.FinalizeCallCount()).To(Equal(1))
			Expect(pkg4Res.FinalizeArgsForCall(0) == finalIndicies.Packages).To(BeTrue())
		})

		It("returns error if package finalize fails", func() {
			pkg1Res := &fakeres.FakeResource{
				NameStub:        func() string { return "pkg1-name" },
				FingerprintStub: func() string { return "pkg1-fp" },
			}
			pkg1 := boshpkg.NewPackage(pkg1Res, nil)

			pkg1Res.FinalizeReturns(errors.New("fake-err"))

			err := releaseDir.VendorPackage(pkg1)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-err"))
		})

		It("returns error if removing old package dir fails", func() {
			pkg1Res := &fakeres.FakeResource{
				NameStub:        func() string { return "pkg1-name" },
				FingerprintStub: func() string { return "pkg1-fp" },
			}
			pkg1 := boshpkg.NewPackage(pkg1Res, nil)

			fs.RemoveAllStub = func(path string) error { return errors.New("fake-err") }

			err := releaseDir.VendorPackage(pkg1)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-err"))
		})

		It("returns error if serializing manifest lock fails", func() {
			pkg1Res := &fakeres.FakeResource{
				NameStub: func() string { return "pkg1-name" },
			}
			pkg1 := boshpkg.NewPackage(pkg1Res, nil)

			err := releaseDir.VendorPackage(pkg1)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Marshaling vendored package 'pkg1-name' spec lock"))
		})

		It("returns error if writing manifest lock fails", func() {
			pkg1Res := &fakeres.FakeResource{
				NameStub:        func() string { return "pkg1-name" },
				FingerprintStub: func() string { return "pkg1-fp" },
			}
			pkg1 := boshpkg.NewPackage(pkg1Res, nil)

			fs.WriteFileErrors["/dir/packages/pkg1-name/spec.lock"] = errors.New("fake-err")

			err := releaseDir.VendorPackage(pkg1)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-err"))
		})
	})

	Describe("FinalizeRelease", func() {
		var (
			release *fakerel.FakeRelease
		)

		BeforeEach(func() {
			release = &fakerel.FakeRelease{
				NameStub:    func() string { return "rel1" },
				VersionStub: func() string { return "ver1" },
				ManifestStub: func() boshman.Manifest {
					return boshman.Manifest{Name: "rel1"}
				},
			}
		})

		It("finalizes release", func() {
			var ops []string

			gitRepo.MustNotBeDirtyStub = func(force bool) (bool, error) {
				ops = append(ops, "dirty")
				return true, nil
			}

			finalReleases.ContainsStub = func(rel boshrel.Release) (bool, error) {
				Expect(rel).To(Equal(release))
				ops = append(ops, "check")
				return false, nil
			}

			release.FinalizeStub = func(indicies boshrel.ArchiveIndicies, parallel int) error {
				Expect(parallel).To(Equal(2))
				Expect(indicies.Jobs).To(Equal(finalIndicies.Jobs)) // unique check
				ops = append(ops, "finalize")
				return nil
			}

			finalReleases.AddStub = func(manifest boshman.Manifest) error {
				Expect(manifest).To(Equal(boshman.Manifest{Name: "rel1"}))
				ops = append(ops, "manifest")
				return nil
			}

			err := releaseDir.FinalizeRelease(release, false)
			Expect(err).ToNot(HaveOccurred())

			Expect(ops).To(Equal([]string{"dirty", "check", "finalize", "manifest"}))
		})

		It("returns error if git is dirty and force is not set", func() {
			gitRepo.MustNotBeDirtyReturns(true, errors.New("dirty"))

			err := releaseDir.FinalizeRelease(release, false)
			Expect(err).To(Equal(errors.New("dirty")))
		})

		It("returns error if checking for a final release fails", func() {
			finalReleases.ContainsReturns(false, errors.New("fake-err"))

			err := releaseDir.FinalizeRelease(release, false)
			Expect(err).To(Equal(errors.New("fake-err")))
		})

		It("returns error if final release index already contains this name/ver", func() {
			finalReleases.ContainsReturns(true, nil)

			err := releaseDir.FinalizeRelease(release, false)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Release 'rel1' version 'ver1' already exists"))
		})

		It("returns error if adding final release fails", func() {
			finalReleases.AddReturns(errors.New("fake-err"))

			err := releaseDir.FinalizeRelease(release, false)
			Expect(err).To(Equal(errors.New("fake-err")))
		})
	})
})
