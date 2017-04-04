package kubo_deployment_tests_test

import (
	"io/ioutil"
	"os"
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("End 2 end run", func() {

	var kuboEnv = pathFromRoot("src/kubo-deployment-tests/resources/test_gcp")

	It("should work now", func() {
		bash.Source(pathToScript("lib/deploy_utils"), nil)
		bash.Source(pathToScript("set_kubeconfig"), nil)
		bash.ExportFunc("kubectl", emptyCallback)
		bash.ExportFunc("bosh-cli", emptyCallback)
		bash.ExportFunc("credhub", emptyCallback)
		bash.SelfPath = "invocationRecorder"

		tmpdir := os.TempDir()
		deployUtilContent := []byte("\n")

		os.MkdirAll(path.Join(tmpdir, "lib"), os.FileMode(0755))
		ioutil.WriteFile(path.Join(tmpdir, "lib/deploy_utils"), deployUtilContent, 0755)

		status, err := bash.Run("main", []string{kuboEnv, "two"})

		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal(0))
	})
})
