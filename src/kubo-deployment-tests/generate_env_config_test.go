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
