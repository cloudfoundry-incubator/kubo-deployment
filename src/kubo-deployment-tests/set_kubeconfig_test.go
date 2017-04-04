package kubo_deployment_tests_test

import (
	"os"
	"path"
	"io/ioutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/progrium/go-basher"
)

var _ = Describe("End 2 end run", func() {

	var kuboEnv = pathFromRoot("src/kubo-deployment-tests/resources/test_gcp")

	It("should work now", func() {
		bash, _ := basher.NewContext("/bin/bash", true)

		bash.Stdout = GinkgoWriter
		bash.Stderr = GinkgoWriter

		bash.CopyEnv()

		bash.Source(pathToScript("lib/deploy_utils"), nil)
		bash.Source(pathToScript("set_kubeconfig"), nil)

		if bash.HandleFuncs(os.Args) {
			os.Exit(0)
		}

		tmpdir := os.TempDir()
		deployUtilContent := []byte("\n")

		os.MkdirAll(path.Join(tmpdir, "lib"), os.FileMode(0755))
   	ioutil.WriteFile(path.Join(tmpdir, "lib/deploy_utils"), deployUtilContent, 0755)

		status, err := bash.Run("main", []string{kuboEnv, "two"})

		Expect(err).NotTo(HaveOccurred())
	  Expect(status).To(Equal(0))
	})
})
