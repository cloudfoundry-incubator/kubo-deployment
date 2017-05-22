package kubo_deployment_tests_test

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	. "github.com/jhvhs/gob-mock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("set_kubeconfig", func() {

	var kuboEnv =  filepath.Join(testEnvironmentPath, "test_gcp")

	BeforeEach(func() {
		bash.Source(pathToScript("lib/deploy_utils"), nil)
		bash.Source(pathToScript("set_kubeconfig"), nil)
		getSetting := Mock("get_setting",
			`[ "$2" == "/external-kubo-port" ] && echo "some-port"
			[ "$2" == "/cf-tcp-router-name" ] && echo "some-url"
			[ "$2" == "/kubo-admin-password" ] && echo "sekret"`)
		mocks := []Gob{Spy("kubectl"), Spy("bosh-cli"), Spy("credhub"), getSetting}
		ApplyMocks(bash, mocks)

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
		Entry("with missing environment", []string{"/missing", "a"}),
	)

	Context("when correct parameters are provided", func() {
		BeforeEach(func() {
			mocks := []Gob{Mock("bosh-cli", `[ "$4" == "/iaas" ] && echo "gcp"`)}
			ApplyMocks(bash, mocks)

			status, err := bash.Run("main", []string{kuboEnv, "deployment-name"})

			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(0))
		})

		It("should set cluster config on kubectl", func() {
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
