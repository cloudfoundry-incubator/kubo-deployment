package kubo_deployment_tests_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/ginkgo/extensions/table"
	"path/filepath"
	"fmt"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Generate manifest", func() {
	BeforeEach(func() {
		bash.Source(pathToScript("generate_service_manifest"), nil)
		bash.Source("_", func(string) ([]byte, error) {
			return []byte(fmt.Sprintf(`repo_directory() { echo -n "%s"; }`, pathFromRoot(""))), nil
		})
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

	Context("successful manifest generation", func() {
		kuboEnv := filepath.Join(resourcesPath, "test_gcp")
		DescribeTable("populated properties", func(params []string) {
			status, err := bash.Run("main", []string{kuboEnv, "klingon"})

			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(0))

			Expect(stdout).To(gbytes.Say(params[0]))
		},
			Entry("deployment name", []string{"\nname: klingon\n"}),
			Entry("stemcell version", []string{"\n  version: stemcell.version\n"}),
			Entry("network name", []string{"\n  networks:\n  - name: network-name\n"}),
			Entry("kubernetes API URL", []string{"\n      kubernetes-api-url: https://a.router.name:101928\n"}),
			Entry("kubernetes external port", []string{"\n      external_kubo_port: 101928\n"}),
			Entry("CF API URL", []string{"\n        api_url: cf.api.url\n"}),
			Entry("CF UAA URL", []string{"\n        uaa_url: cf.uaa.url\n"}),
			Entry("CF Client ID", []string{"\n        uaa_client_id: cf.client.id\n"}),
			Entry("CF Client Secret", []string{"\n        uaa_client_secret: cf.client.secret\n"}),
		)
	})
})
