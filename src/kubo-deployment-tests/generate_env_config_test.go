package kubo_deployment_tests_test

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"

	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("generate_env_config", func() {
	var tmpDir string

	BeforeEach(func() {
		var err error
		tmpDir, err = ioutil.TempDir("", "generate-env-config")
		Expect(err).NotTo(HaveOccurred())

		bash.Source(pathToScript("generate_env_config"), nil)
		bash.Source("", func(string) ([]byte, error) {
			return []byte(fmt.Sprintf(`repo_directory() { echo "%s"; }`, pathFromRoot("src/kubo-deployment-tests/resources"))), nil
		})
	})

	AfterEach(func() {
		Expect(os.RemoveAll(tmpDir)).To(Succeed())
	})

	DescribeTable("with incorrect parameters", func(params []string) {
		status, err := bash.Run("main", params)

		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal(1))
		Expect(stdout).To(gbytes.Say("Usage:"))
	},
		Entry("no params", []string{}),
		Entry("single parameter", []string{"a"}),
		Entry("two parameters", []string{"a", "b"}),
		Entry("invalid iaas", []string{".", "a", "invalid-iaas"}),
	)

	It("should error out if environment dir doesn't exist on disk", func() {
		status, err := bash.Run("main", []string{"/invalid-dir", "a", "b"})

		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal(1))
		Expect(stdout).To(gbytes.Say("should be an existing directory"))
		Expect(stdout).To(gbytes.Say("Usage:"))
	})

	Context("When target file already exists", func() {
		var (
			envDir string
			err    error
		)

		BeforeEach(func() {
			envDir, err = ioutil.TempDir("", "generateEnvConfig")
			Expect(err).ToNot(HaveOccurred())

			err = os.Mkdir(path.Join(envDir, "bosh_name"), 0777)
			Expect(err).ToNot(HaveOccurred())
		})

		AfterEach(func() {
			defer os.RemoveAll(envDir)
		})

		It("does not generate config if director.yml already exists", func() {
			pathToDirectorYml := path.Join(envDir, "bosh_name", "director.yml")
			err = ioutil.WriteFile(pathToDirectorYml, []byte(""), 0777)
			Expect(err).ToNot(HaveOccurred())

			status, err := bash.Run("main", []string{envDir, "bosh_name", "gcp"})

			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(0))
			Expect(stdout).To(gbytes.Say("Skipping config generation because " + pathToDirectorYml + " already exists"))
		})

		It("does not generate config if director-secrets.yml already exists", func() {
			pathToDirectorYml := path.Join(envDir, "bosh_name", "director-secrets.yml")
			err = ioutil.WriteFile(pathToDirectorYml, []byte(""), 0777)
			Expect(err).ToNot(HaveOccurred())

			status, err := bash.Run("main", []string{envDir, "bosh_name", "gcp"})

			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(0))
			Expect(stdout).To(gbytes.Say("Skipping config generation because " + pathToDirectorYml + " already exists"))
		})
	})

	Context("config contains iaas-specific properties", func() {
		var (
			envDir  string
			err     error
			envName = "Lamas"
		)

		BeforeEach(func() {
			envDir, err = ioutil.TempDir("", "generateEnvConfig")
			Expect(err).ToNot(HaveOccurred())

			err = os.Mkdir(path.Join(envDir, envName), 0777)
			Expect(err).ToNot(HaveOccurred())

			bash.Source("", func(string) ([]byte, error) {
				return []byte(fmt.Sprintf(`repo_directory() { echo "%s"; }`, pathFromRoot(""))), nil
			})
		})

		AfterEach(func() {
			defer os.RemoveAll(envDir)
		})

		expectPropertyToExistForIaaS := func(iaas, configFile, propertyName string) {
			status, err := bash.Run("main", []string{envDir, envName, iaas})
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(0))
			directorPath := path.Join(envDir, envName, configFile)
			director, err := ioutil.ReadFile(directorPath)
			Expect(err).NotTo(HaveOccurred())
			_, err = propertyFromYaml("/"+propertyName, director)
			Expect(err).ToNot(HaveOccurred())
		}

		DescribeTable("checks property is present in director.yml", func(iaas string, propertyName string) {
			expectPropertyToExistForIaaS(iaas, "director.yml", propertyName)
		},
			Entry("AWS", "aws", "master_iam_instance_profile"),
			Entry("GCP master", "gcp", "service_account_master"),
			Entry("GCP worker", "gcp", "service_account_worker"),
			Entry("OpenStack", "openstack", "openstack_domain"),
			Entry("vSphere", "vsphere", "vcenter_ip"),
		)

		DescribeTable("checks property is present in director-secrets.yml", func(iaas string, propertyName string) {
			expectPropertyToExistForIaaS(iaas, "director-secrets.yml", propertyName)
		},
			Entry("AWS", "aws", "access_key_id"),
			Entry("OpenStack", "openstack", "openstack_password"),
			Entry("vSphere", "vsphere", "vcenter_password"),
		)
	})

	Context("default settings", func() {
		JustBeforeEach(func() {
			bash.Source("", func(string) ([]byte, error) {
				return []byte(fmt.Sprintf(`repo_directory() { echo "%s"; }`, pathFromRoot(""))), nil
			})
		})

		It("sets default authorization mode to rbac", func() {
			status, _ := bash.Run("main", []string{tmpDir, "b00t", "gcp"})
			Expect(status).To(Equal(0))

			config, err := ioutil.ReadFile(filepath.Join(tmpDir, "b00t/director.yml"))
			Expect(err).NotTo(HaveOccurred())

			expectPathContent("/authorization_mode", config, "rbac")
		})
	})

	It("gracefully concatenates the templates", func() {
		iaas := "aws"
		status, _ := bash.Run("main", []string{tmpDir, "b00t", iaas})
		Expect(status).To(Equal(0))

		config, err := ioutil.ReadFile(filepath.Join(tmpDir, "b00t/director.yml"))
		Expect(err).NotTo(HaveOccurred())

		expectPathContent("/some-other", config, "setting")
		expectPathContent("/iaas", config, iaas)

		secrets, err := ioutil.ReadFile(filepath.Join(tmpDir, "b00t/director-secrets.yml"))
		Expect(err).NotTo(HaveOccurred())

		expectPathContent("/ssshhh", secrets, "ssshhh")
	})
})

func expectPathContent(yamlPath string, yamlSlice []byte, content string) {
	value, err := propertyFromYaml(yamlPath, yamlSlice)
	Expect(err).ToNot(HaveOccurred())
	Expect(value).To(Equal(content))
}
