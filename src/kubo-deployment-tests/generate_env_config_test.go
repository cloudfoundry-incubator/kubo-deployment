package kubo_deployment_tests_test

import (
	"fmt"
	"io/ioutil"

	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("generate_env_config", func() {
	BeforeEach(func() {
		bash.Source(pathToScript("generate_env_config"), nil)
		bash.Source("", func(string) ([]byte, error) {
			return []byte(fmt.Sprintf(`repo_directory() { echo "%s"; }`, pathFromRoot("src/kubo-deployment-tests/resources"))), nil
		})
	})

	AfterEach(func() {
		os.RemoveAll("/tmp/b00t")
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

	It("should error out if environment dir doesn't exist on disk", func(){
		status, err := bash.Run("main", []string{"/invalid-dir", "a", "b"})

		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal(1))
		Expect(stdout).To(gbytes.Say("should be an existing directory"))
		Expect(stdout).To(gbytes.Say("Usage:"))
	})

	It("does not generate config if file already exists", func() {
		tmpDir, err := ioutil.TempDir("", "generateEnvConfig")
		Expect(err).ToNot(HaveOccurred())
		_, err = ioutil.TempFile(tmpDir, "director.yml")
		Expect(err).ToNot(HaveOccurred())

		status, err := bash.Run("main", []string{tmpDir, "bosh_name", "gcp"})

		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal(0))
		Expect(stdout).To(gbytes.Say("Skipping config generation"))
	})

	It("gracefully concatenates the templates", func() {
		iaas := "aws"
		status, _ := bash.Run("main", []string{"/tmp", "b00t", iaas})
		Expect(status).To(Equal(0))

		config, err := ioutil.ReadFile("/tmp/b00t/director.yml")
		Expect(err).NotTo(HaveOccurred())
		configString := string(config)
		Expect(configString).To(ContainSubstring(fmt.Sprintf("\niaas: %s", iaas)))
		Expect(configString).To(ContainSubstring("\nsome-other: setting"))

		secrets, err := ioutil.ReadFile("/tmp/b00t/director-secrets.yml")
		Expect(err).NotTo(HaveOccurred())
		secretsString := string(secrets)
		Expect(secretsString).To(ContainSubstring("\nssshhh: ssshhh"))
	})
})
