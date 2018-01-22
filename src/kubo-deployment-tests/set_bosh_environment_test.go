package kubo_deployment_tests_test

import (
	"fmt"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Setup BOSH environment", func() {

	var kuboEnv = filepath.Join(testEnvironmentPath, "test_gcp_with_creds")

	It("Sets the BOSH environment", func() {
		bash.Export("BOSH_ENV", kuboEnv)
		Expect(bash.Source(pathToScript("set_bosh_environment"), nil)).To(Succeed())

		expectedVariables := map[string]string{
			"BOSH_ENVIRONMENT":   "internal.ip",
			"BOSH_CLIENT":        "bosh_admin",
			"BOSH_CLIENT_SECRET": "test-bosh-admin-client-secret",
			"BOSH_CA_CERT":       "I-am-a-bosh-ca-cert",
		}

		for key, value := range expectedVariables {
			exitCode, err := bash.Run(fmt.Sprintf("echo %s=$%s", key, key), nil)
			Expect(exitCode).To(Equal(0))
			Expect(err).NotTo(HaveOccurred())
			Expect(string(stdout.Contents())).To(ContainSubstring(fmt.Sprintf("%s=%s", key, value)))
		}
	})

	It("Succeeds in a strict context", func() {
		bash.Export("BOSH_ENV", kuboEnv)
		bash.Source("", func(_ string) ([]byte, error) {
			return []byte("set -eu"), nil
		})
		Expect(bash.Source(pathToScript("set_bosh_environment"), nil)).To(Succeed())
		exitCode, err := bash.Run("echo", nil)
		Expect(exitCode).To(Equal(0))
		Expect(err).NotTo(HaveOccurred())
	})
})
