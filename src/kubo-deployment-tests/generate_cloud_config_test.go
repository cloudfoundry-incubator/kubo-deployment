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

	It("Does not include load balancer config for cf-based environment", func() {
		bash.Run("main", []string{filepath.Join(testEnvironmentPath, "test_vsphere")})

		Expect(stdout).NotTo(gbytes.Say("    target_pool: \\(\\(master_target_pool\\)\\)"))
	})

	It("includes load balancer configuration for iaas-based environment", func() {
		bash.Run("main", []string{kuboEnv})

		Expect(stdout).To(gbytes.Say("    target_pool: \\(\\(master_target_pool\\)\\)"))
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

	It("applies extra cloud config ops files when CLOUD_CONFIG_OPS_FILES variable is set", func() {
		opsFiles := []string{
			filepath.Join(resourcesPath, "ops-files", "cloud-config-plus.yml"),
			filepath.Join(resourcesPath, "ops-files", "cloud-config-plus-plus.yml"),
		}

		bash.Export("CLOUD_CONFIG_OPS_FILES", strings.Join(opsFiles, ":"))
		status, err := bash.Run("main", []string{kuboEnv})
		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal(0))

		Expect(stdout).To(gbytes.Say("machine_type: foo"))
		Expect(stdout).To(gbytes.Say("tags: supertest"))
	})

	It("applies extra cloud config ops files when CLOUD_CONFIG_OPS_FILES variable is set to one file", func() {
		opsFiles := []string{
			filepath.Join(resourcesPath, "ops-files", "cloud-config-plus.yml"),
		}

		bash.Export("CLOUD_CONFIG_OPS_FILES", strings.Join(opsFiles, ":"))
		status, err := bash.Run("main", []string{kuboEnv})
		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal(0))

		Expect(stdout).To(gbytes.Say("machine_type: foo"))
	})

	It("hides cloud-config from output with debug flag", func(){
		bash.Export("DEBUG", "1")
		status, err := bash.Run("main", []string{kuboEnv})
		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal(0))

		Expect(stderr).NotTo(gbytes.Say("azs"))
	})

	Context("On GCP", func() {
		Context("With IAAS Routing", func() {
			Context("When service_accounts for master and worker are provided", func() {
				It("adds service account to master and worker", func() {
					command := exec.Command("./generate_cloud_config", "../src/kubo-deployment-tests/resources/environments/test_gcp")
					command.Stdout = bash.Stdout
					command.Stderr = bash.Stderr
					command.Dir = pathToScript("")
					Expect(command.Run()).To(Succeed())

					masterServiceAccount, err := propertyFromYaml("/vm_types/name=master/cloud_properties/service_account", stdout.Contents())
					Expect(err).ToNot(HaveOccurred())
					Expect(masterServiceAccount).To(Equal("master-service-account@google.com"))

					workerServiceAccount, err := propertyFromYaml("/vm_types/name=worker/cloud_properties/service_account", stdout.Contents())
					Expect(err).ToNot(HaveOccurred())
					Expect(workerServiceAccount).To(Equal("worker-service-account@google.com"))
				})
			})

			Context("When service_accounts for master and worker are not provided", func() {
				It("doesn't add service account to master and worker", func() {
					command := exec.Command("./generate_cloud_config", "../src/kubo-deployment-tests/resources/environments/test_gcp_with_service_key")
					command.Stdout = bash.Stdout
					command.Stderr = bash.Stderr
					command.Dir = pathToScript("")
					Expect(command.Run()).To(Succeed())

					_, err := propertyFromYaml("/vm_types/name=master/cloud_properties/service_account", stdout.Contents())
					Expect(err).To(HaveOccurred(), "The master service account should have been removed")
					_, err = propertyFromYaml("/vm_types/name=worker/cloud_properties/service_account", stdout.Contents())
					Expect(err).To(HaveOccurred(), "The worker service account should have been removed")
				})
			})
		})
		Context("With CF Routing", func() {
			Context("When service_accounts for master and worker are provided", func() {
				It("adds service account to master and worker", func() {
					command := exec.Command("./generate_cloud_config", "../src/kubo-deployment-tests/resources/environments/test_gcp_with_cf")
					command.Stdout = bash.Stdout
					command.Stderr = bash.Stderr
					command.Dir = pathToScript("")
					Expect(command.Run()).To(Succeed())

					masterServiceAccount, err := propertyFromYaml("/vm_types/name=master/cloud_properties/service_account", stdout.Contents())
					Expect(err).ToNot(HaveOccurred())
					Expect(masterServiceAccount).To(Equal("master-service-account@google.com"))

					workerServiceAccount, err := propertyFromYaml("/vm_types/name=worker/cloud_properties/service_account", stdout.Contents())
					Expect(err).ToNot(HaveOccurred())
					Expect(workerServiceAccount).To(Equal("worker-service-account@google.com"))
				})
			})

			Context("When service_accounts for master and worker are not provided", func() {
				It("doesn't add service account to master and worker", func() {
					command := exec.Command("./generate_cloud_config", "../src/kubo-deployment-tests/resources/environments/test_gcp_with_cf_and_service_key")
					command.Stdout = bash.Stdout
					command.Stderr = bash.Stderr
					command.Dir = pathToScript("")
					Expect(command.Run()).To(Succeed())

					_, err := propertyFromYaml("/vm_types/name=master/cloud_properties/service_account", stdout.Contents())
					Expect(err).To(HaveOccurred(), "The master service account should have been removed")
					_, err = propertyFromYaml("/vm_types/name=worker/cloud_properties/service_account", stdout.Contents())
					Expect(err).To(HaveOccurred(), "The worker service account should have been removed")
				})
			})
		})
	})
})
