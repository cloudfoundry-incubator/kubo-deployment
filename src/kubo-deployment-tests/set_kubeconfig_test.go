package kubo_deployment_tests_test

import (
	. "github.com/jhvhs/gob-mock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("set_kubeconfig", func() {
	BeforeEach(func() {
		bash.Source(pathToScript("set_kubeconfig"), nil)
		bash.Source("", func(string) ([]byte, error) {
			return repoDirectoryFunction, nil
		})
		credMock := Mock("credhub", `set +x; echo '{"value": {"ca": "certiffy cat"}}'; [[ -z "$DEBUG" ]] || set -x`)
		mocks := []Gob{Spy("kubectl"), credMock}
		ApplyMocks(bash, mocks)

	})

	DescribeTable("with incorrect parameters", func(params []string) {
		status, err := bash.Run("set_kubeconfig", params)

		Expect(err).NotTo(HaveOccurred())
		Expect(status).NotTo(Equal(0))
	},
		Entry("no params", []string{}),
		Entry("single parameter", []string{"a"}),
		Entry("three parameters", []string{"a", "b", "c"}),
		Entry("malformed cluster with no director", []string{"/deployment", "a"}),
		Entry("malformed cluster with no deployment", []string{"director/", "a"}),
		Entry("malformed cluster with no slash", []string{"directordeployment", "a"}),
		Entry("malformed cluster with trailing slash", []string{"director/deployment/", "a"}),
	)

	Context("when correct parameters are provided", func() {
		BeforeEach(func() {
			status, err := bash.Run("set_kubeconfig", []string{"director/deployment", "https://kubernetes.io:8443"})
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(0))
		})

		It("should set cluster config on kubectl", func() {
			Expect(stderr).To(gbytes.Say("kubectl config set-cluster cfcr/director/deployment --server=https://kubernetes.io:8443"))
			Expect(stderr).To(gbytes.Say("kubectl config set-credentials cfcr/director/deployment/cfcr-admin --token=\\w+"))
			Expect(stderr).To(gbytes.Say("kubectl config set-context cfcr/director/deployment --cluster=cfcr/director/deployment --user=cfcr/director/deployment/cfcr-admin"))
			Expect(stderr).To(gbytes.Say("kubectl config use-context cfcr/director/deployment"))
		})
	})

	It("Should hide secrets even when debug flag is set", func() {
		code, err := bash.Run("set -x; set_kubeconfig", []string{"director/deployment", "https://kubernetes.io:8443"})
		Expect(err).NotTo(HaveOccurred())
		Expect(code).To(Equal(0))
		Expect(stderr.Contents()).NotTo(matchDebugOutput("certiffy cat"))
	})
})
