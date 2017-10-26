package kubo_deployment_tests_test

import (
	"path/filepath"

	. "github.com/jhvhs/gob-mock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("set_kubeconfig", func() {

	var kuboEnv = filepath.Join(testEnvironmentPath, "test_gcp")

	BeforeEach(func() {
		bash.Source(pathToScript("set_kubeconfig"), nil)
		bash.Source("", func(string) ([]byte, error) {
			return repoDirectoryFunction, nil
		})
		boshMock := MockOrCallThrough("bosh-cli", `echo "Secret data"`, `[[ "$1" =~ ^int ]] && ! [[ "$2" =~ creds.yml$ ]]`)
		credMock := Mock("credhub", `echo '{"value": {"ca": "certiffy cat"}}'`)
		mocks := []Gob{Spy("kubectl"), boshMock, credMock}
		ApplyMocks(bash, mocks)

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
			status, err := bash.Run("main", []string{kuboEnv, "deployment-name"})

			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(0))
		})

		It("should set cluster config on kubectl", func() {
			Expect(stderr).To(gbytes.Say("kubectl config set-cluster deployment-name --server=https://12.23.34.45:101928"))
		})

		It("should set credentials on kubectl", func() {
			Expect(stderr).To(gbytes.Say("kubectl config set-credentials deployment-name-admin --token=\\w+"))
		})

		It("should set context on kubectl", func() {
			Expect(stderr).To(gbytes.Say("kubectl config set-context kubo-deployment-name --cluster=deployment-name --user=deployment-name-admin"))
		})

		It("should use context on kubectl", func() {
			Expect(stderr).To(gbytes.Say("kubectl config use-context kubo-deployment-name"))
		})
	})
})
