package kubo_deployment_tests_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	basher "github.com/progrium/go-basher"
	"os"
	"github.com/onsi/gomega/gbytes"
	"io"
	"strings"
)

var _ = Describe("Generate cloud config", func() {

	var (
		bash    *basher.Context
		kuboEnv = pathFromRoot("src/kubo-deployment-tests/resources/test_gcp")
	)

	BeforeEach(func() {
		bash, _ = basher.NewContext("/bin/bash", true)
		bash.CopyEnv()
		bash.Source(pathToScript("lib/deploy_utils"), nil)
		bash.Source(pathToScript("generate_cloud_config"), nil)
		if bash.HandleFuncs(os.Args) {
			os.Exit(0)
		}
		bash.ExportFunc("bosh-cli", EmptyCallback)
		bash.ExportFunc("popd", EmptyCallback)
		bash.ExportFunc("pushd", EmptyCallback)
	})

	It("calls bosh-cli with appropriate arguments", func() {
		stdout := gbytes.NewBuffer()
		stderr := gbytes.NewBuffer()
		bash.Stdout = io.MultiWriter(GinkgoWriter, stdout)
		bash.Stderr = io.MultiWriter(GinkgoWriter, stderr)


		bash.SelfPath = "/bin/echo"

		status, err := bash.Run("main", []string{kuboEnv})

		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal(0))
		lines := strings.Split(string(stdout.Contents()), "\n")
		Expect(lines).To(ContainElement("::: bosh-cli int configurations/gcp/cloud-config.yml --vars-file " + kuboEnv + "/director.yml"))
	})

	It("fails with no arguments", func() {
		stdout := gbytes.NewBuffer()
		stderr := gbytes.NewBuffer()
		bash.Stdout = io.MultiWriter(GinkgoWriter, stdout)
		bash.Stderr = io.MultiWriter(GinkgoWriter, stderr)


		bash.SelfPath = "/bin/echo"

		status, err := bash.Run("main", []string{})

		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal(1))
	})

	It("should temporarily step into an upper level directory", func() {
		stdout := gbytes.NewBuffer()
		stderr := gbytes.NewBuffer()
		bash.Stdout = io.MultiWriter(GinkgoWriter, stdout)
		bash.Stderr = io.MultiWriter(GinkgoWriter, stderr)

		bash.Source("_", func(string) ([]byte, error) {
			return []byte(`
				callCounter=0
				invocationRecorder() {
					callCounter=$(expr $callCounter + 1)
					echo "[$callCounter] $@" | tee /dev/fd/2
				}
			`), nil
		})
		bash.SelfPath = "invocationRecorder"

		status, err := bash.Run("main", []string{kuboEnv})
		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal(0))

		errOutput := string(stderr.Contents())

		// Our test executable is /bin/bash, so the path should be one level up: /bin/../
		Expect(errOutput).To(ContainSubstring("[1] ::: pushd /bin/../"))
		Expect(errOutput).To(ContainSubstring("[2] ::: bosh-cli"))
		Expect(errOutput).To(ContainSubstring("[3] ::: popd"))
	})
})
