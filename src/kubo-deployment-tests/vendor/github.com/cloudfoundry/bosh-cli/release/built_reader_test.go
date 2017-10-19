package release_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry/bosh-cli/release"
	fakerel "github.com/cloudfoundry/bosh-cli/release/releasefakes"
	fakeres "github.com/cloudfoundry/bosh-cli/release/resource/resourcefakes"
)

var _ = Describe("BuiltReader", func() {
	Describe("Read", func() {
		var (
			innerReader   *fakerel.FakeReader
			devIndicies   ArchiveIndicies
			finalIndicies ArchiveIndicies
			reader        BuiltReader
		)

		BeforeEach(func() {
			innerReader = &fakerel.FakeReader{}

			devIndicies = ArchiveIndicies{
				Jobs: &fakeres.FakeArchiveIndex{
					FindStub: func(_, _ string) (string, string, error) { return "dev", "", nil },
				},
				Packages: &fakeres.FakeArchiveIndex{},
				Licenses: &fakeres.FakeArchiveIndex{},
			}

			finalIndicies = ArchiveIndicies{
				Jobs: &fakeres.FakeArchiveIndex{
					FindStub: func(_, _ string) (string, string, error) { return "final", "", nil },
				},
				Packages: &fakeres.FakeArchiveIndex{},
				Licenses: &fakeres.FakeArchiveIndex{},
			}

			reader = NewBuiltReader(innerReader, devIndicies, finalIndicies, 3)
		})

		It("reads and builds release", func() {
			readRelease := &fakerel.FakeRelease{}
			innerReader.ReadReturns(readRelease, nil)

			release, err := reader.Read("/release.tgz")
			Expect(err).ToNot(HaveOccurred())
			Expect(release).To(Equal(readRelease))

			Expect(readRelease.BuildCallCount()).To(Equal(1))

			dev, final, parallel := readRelease.BuildArgsForCall(0)
			Expect(dev).To(Equal(devIndicies))
			Expect(final).To(Equal(finalIndicies))
			Expect(parallel).To(Equal(3))
		})

		It("returns error if read fails", func() {
			innerReader.ReadReturns(nil, errors.New("fake-err"))

			_, err := reader.Read("/release.tgz")
			Expect(err).To(Equal(errors.New("fake-err")))
		})

		It("returns error if building fails", func() {
			readRelease := &fakerel.FakeRelease{}
			innerReader.ReadReturns(readRelease, nil)

			readRelease.BuildReturns(errors.New("fake-err"))

			_, err := reader.Read("/release.tgz")
			Expect(err).To(Equal(errors.New("fake-err")))
		})

		Context("parallel is less than 1", func() {
			BeforeEach(func() {
				reader = NewBuiltReader(innerReader, devIndicies, finalIndicies, 0)
			})

			It("builds with parallel of 1", func() {
				readRelease := &fakerel.FakeRelease{}
				innerReader.ReadReturns(readRelease, nil)

				_, err := reader.Read("/release.tgz")
				Expect(err).NotTo(HaveOccurred())

				_, _, parallel := readRelease.BuildArgsForCall(0)
				Expect(parallel).To(Equal(1))
			})
		})
	})
})
