package kubo_deployment_tests_test

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	. "github.com/jhvhs/gob-mock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Generate cloud config", func() {

	var kuboEnv = filepath.Join(testEnvironmentPath, "test_gcp")

	BeforeEach(func() {
		bash.Source(pathToScript("lib/deploy_utils"), nil)
		bash.Source(pathToScript("generate_cloud_config"), nil)
		mocks := []Gob{Spy("pushd"), Spy("popd"), Spy("bosh-cli")}
		ApplyMocks(bash, mocks)
	})

	It("calls bosh-cli with appropriate arguments", func() {
		boshMock := Mock("bosh-cli", `[ "$4" == "/iaas" ] && echo "gcp"`)
		ApplyMocks(bash, []Gob{boshMock})
		status, err := bash.Run("main", []string{kuboEnv})

		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal(0))
		Expect(stderr).To(gbytes.Say("bosh-cli int configurations/gcp/cloud-config.yml --vars-file " + kuboEnv + "/director.yml"))
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

		status, err := bash.Run("main", []string{kuboEnv})
		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal(0))


		// Our test executable is ~/.basher/bash, so the path should be one level up
		targetPath := strings.Replace(bashPath, "/bash", "/../", 1)
		Expect(stderr).To(gbytes.Say(fmt.Sprintf("<1> pushd %s", targetPath)))
		Expect(stderr).To(gbytes.Say("<2> bosh-cli"))
		Expect(stderr).To(gbytes.Say("<3> popd"))
	})
})
