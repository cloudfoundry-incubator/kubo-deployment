package kubo_deployment_tests_test

import (
	"fmt"
	"os/exec"
	"path/filepath"

	. "github.com/jhvhs/gob-mock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Generate cloud config", func() {

	var kuboEnv = filepath.Join(testEnvironmentPath, "test_gcp")

	BeforeEach(func() {
		bash.Source(pathToScript("generate_cloud_config"), nil)
		bash.Source("", func(string) ([]byte, error) {
			return repoDirectoryFunction, nil
		})
	})

	It("calls bosh-cli with appropriate arguments", func() {
		boshMock := SpyAndCallThrough("bosh-cli")
		ApplyMocks(bash, []Gob{boshMock})
		status, err := bash.Run("main", []string{kuboEnv})

		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal(0))
		cloudConfig := pathFromRoot("configurations/gcp/cloud-config.yml")
		boshCmd := fmt.Sprintf("bosh-cli int %s --vars-file %s/director.yml", cloudConfig, kuboEnv)
		Expect(stderr).To(gbytes.Say(boshCmd))
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
})
