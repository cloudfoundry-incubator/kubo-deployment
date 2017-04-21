package kubo_deployment_tests_test

import (
	"io/ioutil"
	"os"
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"path/filepath"
)

var _ = Describe("set_kubeconfig", func() {

	var kuboEnv =  filepath.Join(testEnvironmentPath, "test_gcp")

	BeforeEach(func() {
		bash.Source(pathToScript("lib/deploy_utils"), nil)
		bash.Source(pathToScript("set_kubeconfig"), nil)
		bash.Source("__", func(string) ([]byte, error) {
			return []byte(`get_setting() {
				[ "$2" == "/external-kubo-port" ] && echo "some-port";
				[ "$2" == "/cf-tcp-router-name" ] && echo "some-url";
				[ "$2" == "/kubo-admin-password" ] && echo "sekret";
				return 0;
			}`), nil
		})
		bash.ExportFunc("kubectl", emptyCallback)
		bash.ExportFunc("bosh-cli", emptyCallback)
		bash.ExportFunc("credhub", emptyCallback)
		bash.SelfPath = "invocationRecorder"

		tmpdir := os.TempDir()
		deployUtilContent := []byte("\n")

		os.MkdirAll(path.Join(tmpdir, "lib"), os.FileMode(0755))
		ioutil.WriteFile(path.Join(tmpdir, "lib/deploy_utils"), deployUtilContent, 0755)
	})

	DescribeTable("with incorrect parameters", func(params []string) {
		status, err := bash.Run("main", params)

		Expect(err).NotTo(HaveOccurred())
		Expect(status).NotTo(Equal(0))
	},
		Entry("no params", []string{}),
		Entry("single parameter", []string{"a"}),
		Entry("three parameters", []string{"a", "b", "c"}),
		Entry("with missing environemnt", []string{"/missing", "a"}),
	)

	Context("when correct parameters are provided", func() {
		BeforeEach(func() {
			bash.Source("__", func(string) ([]byte, error) {
				return []byte(`bosh-cli() {
					[ "$4" == "/iaas" ] && echo "gcp";
					return 0;
				}`), nil
			})
			status, err := bash.Run("main", []string{kuboEnv, "deployment-name"})

			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(0))
		})

		It("should set cluser config on kubectl", func() {
			Expect(stderr).To(gbytes.Say("kubectl config set-cluster deployment-name --server=https://some-url:some-port"))
		})

		It("should set credentials on kubectl", func() {
			Expect(stderr).To(gbytes.Say("kubectl config set-credentials deployment-name-admin --token=sekret"))
		})

		It("should set context on kubectl", func() {
			Expect(stderr).To(gbytes.Say("kubectl config set-context kubo-deployment-name --cluster=deployment-name --user=deployment-name-admin"))
		})

		It("should use context on kubectl", func() {
			Expect(stderr).To(gbytes.Say("kubectl config use-context kubo-deployment-name"))
		})
	})
})
