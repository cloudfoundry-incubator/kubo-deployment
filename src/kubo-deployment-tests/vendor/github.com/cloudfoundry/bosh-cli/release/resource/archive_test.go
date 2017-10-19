package resource_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	boshcrypto "github.com/cloudfoundry/bosh-utils/crypto"
	boshcmd "github.com/cloudfoundry/bosh-utils/fileutil"
	fakecmd "github.com/cloudfoundry/bosh-utils/fileutil/fakes"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	fakesys "github.com/cloudfoundry/bosh-utils/system/fakes"
	boshuuid "github.com/cloudfoundry/bosh-utils/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	bicrypto "github.com/cloudfoundry/bosh-cli/crypto"
	fakecrypto "github.com/cloudfoundry/bosh-cli/crypto/fakes"
	. "github.com/cloudfoundry/bosh-cli/release/resource"
	fakeres "github.com/cloudfoundry/bosh-cli/release/resource/resourcefakes"
)

var _ = Describe("Archive", func() {
	var (
		archive Archive
	)

	BeforeEach(func() {
		archive = nil
	})

	Describe("Fingerprint", func() {
		var (
			fingerprinter    *fakeres.FakeFingerprinter
			digestCalculator *fakecrypto.FakeDigestCalculator
			compressor       *fakecmd.FakeCompressor
			cmdRunner        *fakesys.FakeCmdRunner
			fs               *fakesys.FakeFileSystem
		)

		BeforeEach(func() {
			releaseDirPath := filepath.Join("/", "tmp", "release")
			fingerprinter = &fakeres.FakeFingerprinter{}
			digestCalculator = fakecrypto.NewFakeDigestCalculator()
			compressor = fakecmd.NewFakeCompressor()
			cmdRunner = fakesys.NewFakeCmdRunner()
			fs = fakesys.NewFakeFileSystem()
			archiveFactoryArgs := ArchiveFactoryArgs{
				Files:          []File{NewFile(filepath.Join("/", "tmp", "file"), filepath.Join("/", "tmp"))},
				PrepFiles:      []File{NewFile(filepath.Join("/", "tmp", "prep-file"), filepath.Join("/", "tmp"))},
				Chunks:         []string{"chunk"},
				FollowSymlinks: false,
			}

			archive = NewArchiveImpl(
				archiveFactoryArgs,
				releaseDirPath,
				fingerprinter,
				compressor,
				digestCalculator,
				cmdRunner,
				fs,
			)
		})

		It("returns fingerprint", func() {
			fingerprinter.CalculateReturns("fp", nil)

			fp, err := archive.Fingerprint()
			Expect(err).ToNot(HaveOccurred())
			Expect(fp).To(Equal("fp"))

			files, chunks := fingerprinter.CalculateArgsForCall(0)
			Expect(files).To(Equal([]File{NewFile(filepath.Join("/", "tmp", "file"), filepath.Join("/", "tmp"))}))
			Expect(chunks).To(Equal([]string{"chunk"}))
		})

		It("returns error", func() {
			fingerprinter.CalculateReturns("", errors.New("fake-err"))

			_, err := archive.Fingerprint()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-err"))
		})
	})

	Describe("Build", func() {
		var (
			uniqueDir string
			fs        boshsys.FileSystem

			compressor       boshcmd.Compressor
			digestCalculator bicrypto.DigestCalculator
			followSymlinks   bool
			existingUmask    int
		)

		BeforeEach(func() {
			followSymlinks = false
			existingUmask = syscall.Umask(0)
		})

		JustBeforeEach(func() {
			releaseDirPath := filepath.Join("/", "tmp", "release")

			suffix, err := boshuuid.NewGenerator().Generate()
			Expect(err).ToNot(HaveOccurred())

			uniqueDir = filepath.Join("/", "tmp", suffix)

			/*
				Create a file system with the following structure:
					/
				 ├── file1
				 ├── dir (directory)
				 │   ├── file2
				 │   ├── file3
				 │   ├── symlink-dir-target (directory)
				 │   │   └── file4
				 │   ├── symlink-dir -> symlink-dir-target
				 │   ├── symlink-file -> ../file1
				 │   └── symlink-file-missing -> nonexistant-file
				 ├── run-build-dir
				 ├── run-release-dir
				 └── run-file3
			*/

			logger := boshlog.NewLogger(boshlog.LevelNone)
			fs = boshsys.NewOsFileSystemWithStrictTempRoot(logger)

			err = fs.ChangeTempRoot(uniqueDir)
			Expect(err).ToNot(HaveOccurred())

			err = fs.WriteFileString(filepath.Join(uniqueDir, "file1"), "file1")
			Expect(err).ToNot(HaveOccurred())

			err = fs.Chmod(filepath.Join(uniqueDir, "file1"), os.FileMode(0600))
			Expect(err).ToNot(HaveOccurred())

			err = fs.MkdirAll(filepath.Join(uniqueDir, "dir"), os.FileMode(0755))
			Expect(err).ToNot(HaveOccurred())

			err = fs.WriteFileString(filepath.Join(uniqueDir, "dir", "file2"), "file2")
			Expect(err).ToNot(HaveOccurred())

			err = fs.Chmod(filepath.Join(uniqueDir, "dir", "file2"), os.FileMode(0744))
			Expect(err).ToNot(HaveOccurred())

			err = fs.WriteFileString(filepath.Join(uniqueDir, "dir", "file3"), "file3")
			Expect(err).ToNot(HaveOccurred())

			err = fs.MkdirAll(filepath.Join(uniqueDir, "dir", "symlink-dir-target"), os.FileMode(0744))
			Expect(err).ToNot(HaveOccurred())

			err = fs.WriteFileString(filepath.Join(uniqueDir, "dir", "symlink-dir-target", "file4"), "file4")
			Expect(err).ToNot(HaveOccurred())

			err = fs.Symlink("symlink-dir-target", filepath.Join(uniqueDir, "dir", "symlink-dir"))
			Expect(err).ToNot(HaveOccurred())

			err = fs.Symlink("../file1", filepath.Join(uniqueDir, "dir", "symlink-file"))
			Expect(err).ToNot(HaveOccurred())

			if !followSymlinks {
				err = fs.Symlink("nonexistant-file", filepath.Join(uniqueDir, "dir", "symlink-file-missing"))
				Expect(err).ToNot(HaveOccurred())
			}

			err = fs.WriteFileString(filepath.Join(uniqueDir, "run-build-dir"), "echo -n $BUILD_DIR > build-dir")
			Expect(err).ToNot(HaveOccurred())

			err = fs.WriteFileString(filepath.Join(uniqueDir, "run-release-dir"), "echo -n $RELEASE_DIR > release-dir")
			Expect(err).ToNot(HaveOccurred())

			err = fs.WriteFileString(filepath.Join(uniqueDir, "run-file3"), "rm dir/file3")
			Expect(err).ToNot(HaveOccurred())

			digestCalculator = bicrypto.NewDigestCalculator(fs, []boshcrypto.Algorithm{boshcrypto.DigestAlgorithmSHA1})
			fingerprinter := NewFingerprinterImpl(digestCalculator, fs, followSymlinks)
			cmdRunner := boshsys.NewExecCmdRunner(logger)
			compressor = boshcmd.NewTarballCompressor(cmdRunner, fs)

			files :=
				[]File{
					NewFile(filepath.Join(uniqueDir, "file1"), uniqueDir),
					NewFile(filepath.Join(uniqueDir, "dir", "file2"), uniqueDir),
					NewFile(filepath.Join(uniqueDir, "dir", "file3"), uniqueDir),
					NewFile(filepath.Join(uniqueDir, "dir", "symlink-file"), uniqueDir),
				}

			if followSymlinks {
				files = append(files, NewFile(filepath.Join(uniqueDir, "dir", "symlink-dir", "file4"), uniqueDir))
			} else {
				files = append(files, NewFile(filepath.Join(uniqueDir, "dir", "symlink-dir"), uniqueDir))
				files = append(files, NewFile(filepath.Join(uniqueDir, "dir", "symlink-file-missing"), uniqueDir))
			}

			archiveFactoryArgs := ArchiveFactoryArgs{
				Files: files,
				PrepFiles: []File{
					NewFile(filepath.Join(uniqueDir, "run-build-dir"), uniqueDir),
					NewFile(filepath.Join(uniqueDir, "run-release-dir"), uniqueDir),
					NewFile(filepath.Join(uniqueDir, "run-file3"), uniqueDir),
				},
				Chunks:         []string{"chunk"},
				FollowSymlinks: followSymlinks,
			}

			archive = NewArchiveImpl(
				archiveFactoryArgs,
				releaseDirPath,
				fingerprinter,
				compressor,
				digestCalculator,
				cmdRunner,
				fs,
			)
		})

		AfterEach(func() {
			if fs != nil {
				_ = fs.RemoveAll(uniqueDir)
			}
			syscall.Umask(existingUmask)
		})

		modeAsStr := func(m os.FileMode) string {
			return fmt.Sprintf("%#o", m)
		}

		It("returns archive, sha1 when built successfully", func() {
			archivePath, archiveSHA1, err := archive.Build("31a86e1b2b76e47ca5455645bb35018fe7f73e5d")
			Expect(err).ToNot(HaveOccurred())

			actualArchiveDigest, err := digestCalculator.Calculate(archivePath)
			Expect(err).ToNot(HaveOccurred())
			Expect(actualArchiveDigest).To(Equal(archiveSHA1))

			decompPath, err := fs.TempDir("test-resource")
			Expect(err).ToNot(HaveOccurred())

			err = compressor.DecompressFileToDir(archivePath, decompPath, boshcmd.CompressorOptions{})
			Expect(err).ToNot(HaveOccurred())

			{
				// Copies specified files
				Expect(fs.ReadFileString(filepath.Join(decompPath, "file1"))).To(Equal("file1"))
				Expect(fs.ReadFileString(filepath.Join(decompPath, "dir", "file2"))).To(Equal("file2"))

				// Copies specified symlinks
				stat, err := fs.Lstat(filepath.Join(decompPath, "dir", "symlink-file"))
				Expect(err).ToNot(HaveOccurred())
				Expect(stat.Mode()&os.ModeSymlink != 0).To(BeTrue())
				Expect(fs.Readlink(filepath.Join(decompPath, "dir", "symlink-file"))).To(Equal("../file1"))

				stat, err = fs.Lstat(filepath.Join(decompPath, "dir", "symlink-file-missing"))
				Expect(err).ToNot(HaveOccurred())
				Expect(stat.Mode()&os.ModeSymlink != 0).To(BeTrue())
				Expect(fs.Readlink(filepath.Join(decompPath, "dir", "symlink-file-missing"))).To(Equal("nonexistant-file"))

				stat, err = fs.Lstat(filepath.Join(decompPath, "dir", "symlink-dir"))
				Expect(err).ToNot(HaveOccurred())
				Expect(stat.Mode()&os.ModeSymlink != 0).To(BeTrue())
				Expect(fs.Readlink(filepath.Join(decompPath, "dir", "symlink-dir"))).To(Equal("symlink-dir-target"))
				Expect(fs.FileExists(filepath.Join(decompPath, "dir", "simlink-dir-target"))).To(BeFalse())

				// Dir permissions
				stat, err = fs.Stat(filepath.Join(decompPath, "dir"))
				Expect(err).ToNot(HaveOccurred())
				Expect(modeAsStr(stat.Mode())).To(Equal("020000000755")) // 02... is for directory

				// File permissions
				stat, err = fs.Stat(filepath.Join(decompPath, "file1"))
				Expect(err).ToNot(HaveOccurred())
				Expect(modeAsStr(stat.Mode())).To(Equal("0644"))
				stat, err = fs.Stat(filepath.Join(decompPath, "dir"))
				Expect(err).ToNot(HaveOccurred())
				Expect(modeAsStr(stat.Mode())).To(Equal("020000000755"))
				stat, err = fs.Stat(filepath.Join(decompPath, "dir", "file2"))
				Expect(err).ToNot(HaveOccurred())
				Expect(modeAsStr(stat.Mode())).To(Equal("0755"))
			}

			{
				// Runs scripts
				Expect(fs.ReadFileString(filepath.Join(decompPath, "build-dir"))).ToNot(BeEmpty())
				Expect(fs.ReadFileString(filepath.Join(decompPath, "release-dir"))).To(Equal(filepath.Join("/", "tmp", "release")))
				Expect(fs.FileExists(filepath.Join(decompPath, "dir", "file3"))).To(BeFalse())
			}

			{
				// Deletes scripts
				Expect(fs.FileExists(filepath.Join(decompPath, "run-build-dir"))).To(BeFalse())
				Expect(fs.FileExists(filepath.Join(decompPath, "run-release-dir"))).To(BeFalse())
			}
		})

		Context("when following symlinks", func() {
			BeforeEach(func() {
				followSymlinks = true
			})

			It("copies the contents of the symlink", func() {
				archivePath, archiveSHA1, err := archive.Build("1e18c219903f57abe1d28730660fe387e077f378")
				Expect(err).ToNot(HaveOccurred())

				actualArchiveDigest, err := digestCalculator.Calculate(archivePath)
				Expect(err).ToNot(HaveOccurred())
				Expect(actualArchiveDigest).To(Equal(archiveSHA1))

				decompPath, err := fs.TempDir("test-resource")
				Expect(err).ToNot(HaveOccurred())

				err = compressor.DecompressFileToDir(archivePath, decompPath, boshcmd.CompressorOptions{})
				Expect(err).ToNot(HaveOccurred())

				// Copies the contents of the specified symlinks
				Expect(fs.ReadFileString(filepath.Join(decompPath, "dir", "symlink-file"))).To(Equal("file1"))

				stat, err := fs.Stat(filepath.Join(decompPath, "dir", "symlink-dir"))
				Expect(err).ToNot(HaveOccurred())
				Expect(modeAsStr(stat.Mode())).To(Equal("020000000744")) // 02... is for directory
				Expect(fs.ReadFileString(filepath.Join(decompPath, "dir", "symlink-dir", "file4"))).To(Equal("file4"))
			})
		})
	})
})
