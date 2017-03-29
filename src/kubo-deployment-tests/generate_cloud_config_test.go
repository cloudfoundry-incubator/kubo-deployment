package kubo_deployment_tests_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	basher "github.com/progrium/go-basher"
)

var _ = Describe("Generate cloud config", func() {
	It("should work now", func() {
		bash, _ := basher.NewContext("/bin/bash", true)
		bash.Stdout = GinkgoWriter
		bash.Stderr = GinkgoWriter

		bash.CopyEnv()

		bash.Source(pathToScript("lib/deploy_utils"), nil)
		bash.Source(pathToScript("generate_cloud_config"), nil)
		status, err := bash.Run("main", []string{pathFromRoot("src/kubo-deployment-tests/resources/test_gcp")})

		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal(0))
	})
})
