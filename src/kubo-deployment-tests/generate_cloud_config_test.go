package kubo_deployment_tests_test

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Generate cloud config", func() {

	var kuboEnv = filepath.Join(testEnvironmentPath, "test_gcp")

	BeforeEach(func() {
		bash.Source(pathToScript("lib/deploy_utils"), nil)
		bash.Source(pathToScript("generate_cloud_config"), nil)
		bash.ExportFunc("bosh-cli", emptyCallback)
		bash.ExportFunc("popd", emptyCallback)
		bash.ExportFunc("pushd", emptyCallback)
		bash.SelfPath = "/bin/echo"
	})

	It("calls bosh-cli with appropriate arguments", func() {
		bash.Source("__", func(string) ([]byte, error) {
			return []byte(`bosh-cli() {
				[ "$4" == "/iaas" ] && echo "gcp";
				[ "$4" != "/iaas" ] && echo "bosh-cli $@";
				return 0;
			}`), nil
		})

		status, err := bash.Run("main", []string{kuboEnv})

		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal(0))
		lines := strings.Split(string(stdout.Contents()), "\n")
		Expect(lines).To(ContainElement("bosh-cli int configurations/gcp/cloud-config.yml --vars-file " + kuboEnv + "/director.yml"))
	})

	It("fails with no arguments", func() {
		status, err := bash.Run("main", []string{})

		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal(1))
	})

	It("expands the bosh environment path to absolute value", func() {
		command := exec.Command("./generate_cloud_config", "../src/kubo-deployment-tests/resources/environments/test_gcp")
		command.Stdout = bash.Stdout
		command.Stderr = bash.Stderr
		command.Dir = pathToScript("")
		Expect(command.Run()).To(Succeed())
	})

	It("should temporarily step into an upper level directory", func() {
		bash.SelfPath = "invocationRecorder"

		status, err := bash.Run("main", []string{kuboEnv})
		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal(0))

		errOutput := string(stderr.Contents())

		// Our test executable is ~/.basher/bash, so the path should be one level up
		targetPath := strings.Replace(bashPath, "/bash", "/../", 1)
		Expect(errOutput).To(ContainSubstring(fmt.Sprintf("[1] ::: pushd %s", targetPath)))
		Expect(errOutput).To(ContainSubstring("[2] ::: bosh-cli"))
		Expect(errOutput).To(ContainSubstring("[3] ::: popd"))
	})
})
