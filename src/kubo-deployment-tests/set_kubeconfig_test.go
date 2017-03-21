package kubo_deployment_tests_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"path/filepath"
	"runtime"
	"github.com/progrium/go-basher"
	"os"
)

var credhubArgs [][]string

func credhub(args []string) {
	println("I can haz cheesybytez")
	credhubArgs = append(credhubArgs, args)
}

func pathToScript(name string) string {
	return pathFromRoot("bin/" + name)
}

func pathFromRoot(relativePath string) string {
	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)
	return filepath.Join(currentDir, "..", "..", relativePath)
}

var _ = Describe("End 2 end run", func() {
	BeforeEach(func() {
		credhubArgs = [][]string{{}}
	})

	It("should work now", func() {
		bash, _ := basher.NewContext("/bin/bash", true)
		bash.ExportFunc("credhub", credhub)

		bash.Stdout = GinkgoWriter
		bash.Stderr = GinkgoWriter

		bash.CopyEnv()

		bash.Source(pathToScript("lib/deploy_utils"), nil)
		bash.Source(pathToScript("set_kubeconfig"), nil)

		if bash.HandleFuncs(os.Args) {
			os.Exit(0)
		}

		status, err :=
			bash.Run("main", []string{"one", "two"})

		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal(0))
		Expect(credhubArgs).To(HaveLen(1))
	})
})
