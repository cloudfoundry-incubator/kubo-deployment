package kubo_deployment_tests_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/ginkgo/extensions/table"
	"path"
	"fmt"
	"path/filepath"
)

var _ = Describe("Deploy KuBOSH", func() {
	validEnvironment := path.Join(testEnvironmentPath, "test_gcp")
	BeforeEach(func() {
		bash.ExportFunc("bosh-cli", emptyCallback)
		bash.ExportFunc("credhub", emptyCallback)
	})

	Context("fails", func() {

		DescribeTable("when wrong number of arguments is used", func(params []string) {
			code, err := bash.Run(pathToScript("deploy_bosh"), params)
			Expect(err).NotTo(HaveOccurred())
			Expect(code).NotTo(Equal(0))

		},
			Entry("has no arguments", []string{}),
			Entry("has one argument", []string{"gcp"}),
			Entry("has three arguments", []string{"gcp", "foo", "bar"}),
		)

		It("is given an invalid environment path", func() {
			code, err := bash.Run(pathToScript("deploy_bosh"), []string{pathFromRoot(""), pathFromRoot("README.md")})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).NotTo(Equal(0))
		})

		It("is given a non-existing file", func() {
			code, err := bash.Run(pathToScript("deploy_bosh"), []string{validEnvironment, pathFromRoot("non-existing.file")})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).NotTo(Equal(0))
		})
	})

	Context("succeeds", func() {
		BeforeEach(func() {
			bash.SelfPath = "invocationRecorder"
			bash.Source(pathToScript("deploy_bosh"), nil)
			boshCliMock := filepath.Join(resourcesPath, "lib", "bosh_cli_mock.sh")
			bash.Source(boshCliMock, nil)
			bash.Source("_", func(string) ([]byte, error) {
				repoDirectory := fmt.Sprintf(`
				repo_directory() { echo "%s"; }
				export -f bosh-cli
				`, pathFromRoot(""))
				return []byte(repoDirectory), nil
			})
		})

		It("runs with a valid environment and an extra file", func() {
			bash.Debug = true
			code, err := bash.Run("main", []string{validEnvironment, pathFromRoot("README.md")})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))
		})
	})
})
