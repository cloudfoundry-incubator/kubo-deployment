package kubo_deployment_tests_test

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	. "github.com/jhvhs/gob-mock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("UpdateStemcell", func() {

	var mockManifest, manifest string

	BeforeEach(func() {
		bash.Source(pathToScript("update_stemcell"), nil)
		bash.Source("", func(string) ([]byte, error) {
			return repoDirectoryFunction, nil
		})

		manifest = pathFromRoot("manifests/cfcr.yml")
		mockManifest = "/tmp/mock-cfcr.yml"
		cpCmd := exec.Command("cp", "-f", manifest, mockManifest)
		err := cpCmd.Run()
		Expect(err).ToNot(HaveOccurred())

		manifestFileFunctionMock := Mock("manifest_file", fmt.Sprintf("echo %s", mockManifest))
		ApplyMocks(bash, []Gob{manifestFileFunctionMock})
	})

	It("should update the manifest with the given version when there's a different version", func() {

		exitCode, err := bash.Run("main", []string{"new-stemcell-version"})
		Expect(err).ToNot(HaveOccurred())
		Expect(exitCode).To(Equal(0))

		cmd := exec.Command("bosh", "int", mockManifest, "--path=/stemcells/0/version")
		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))
		Eventually(session).Should(gbytes.Say("^new-stemcell-version\n$"))
	})

	It("should not update the manifest when the version is the same", func() {

		fileInfo, err := os.Stat(mockManifest)
		lastModTimeBefore := fileInfo.ModTime()

		cmd := exec.Command("bosh", "int", mockManifest, "--path=/stemcells/0/version")
		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))
		currentVersion := strings.TrimSuffix(string(session.Out.Contents()), "\n")

		exitCode, err := bash.Run("main", []string{currentVersion})
		Expect(err).ToNot(HaveOccurred())
		Expect(exitCode).To(Equal(0))

		fileInfo, err = os.Stat(mockManifest)
		lastModTimeAfter := fileInfo.ModTime()
		Expect(lastModTimeAfter).To(Equal(lastModTimeBefore))
	})

	It("should keep the order of the manifest the same", func() {

		exitCode, err := bash.Run("main", []string{"new-stemcell-version"})
		Expect(err).ToNot(HaveOccurred())
		Expect(exitCode).To(Equal(0))

		// diff should only have 2 lines of change: the old version and the new version
		cmd := exec.Command("bash", "-c", fmt.Sprintf("diff -U 0 %s %s | grep -v '^@' | grep -v '^---' | grep -v '^+++' | wc -l", manifest, mockManifest))
		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))
		Eventually(session).Should(gbytes.Say("^\\s*2\n$"))
	})

})
