package kubo_deployment_tests_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"

	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Generate manifest", func() {
	BeforeEach(func() {
		bash.Source(pathToScript("generate_kubo_manifest"), nil)
		bash.Source("_", func(string) ([]byte, error) {
			return repoDirectoryFunction, nil
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
		kuboEnv := filepath.Join(testEnvironmentPath, "test_gcp")
		AfterEach(func() {
			files, _ := filepath.Glob(testEnvironmentPath + "/**/*creds.yml")
			for _, f := range files {
				if !strings.Contains(f, "with_creds/creds.yml") {
					os.Remove(f)
				}
			}
		})

		DescribeTable("populated properties for CF-based deployment", func(line string) {
			cfEnv := filepath.Join(testEnvironmentPath, "test_vsphere_with_creds")
			status, err := bash.Run("main", []string{cfEnv, "klingon"})

			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(0))

			Expect(stdout).To(gbytes.Say(line))
		},
			Entry("deployment name", "\nname: klingon\n"),
			Entry("stemcell version", "\n  version: stemcell.version\n"),
			Entry("network name", "\n  networks:\n  - name: network-name\n"),
			Entry("kubernetes API URL", "\n      kubernetes-api-url: https://a.router.name:101928\n"),
			Entry("kubernetes external port", "\n      external_kubo_port: 101928\n"),
			Entry("CF API URL", "\n        api_url: cf.api.url\n"),
			Entry("CF UAA URL", "\n        uaa_url: cf.uaa.url\n"),
			Entry("CF Client ID", "\n        uaa_client_id: cf.client.id\n"),
			Entry("CF Client Secret", "\n        uaa_client_secret: cf.client.secret\n"),
			Entry("Auto-generated kubelet password", "\n      kubelet-password: \\(\\(kubelet-password\\)\\)\n"),
			Entry("Auto-generated admin password", "\n      admin-password: \\(\\(kubo-admin-password\\)\\)\n"),
		)


		DescribeTable("populated properties for IaaS-based deployment", func(line string) {
			status, err := bash.Run("main", []string{kuboEnv, "grinder"})

			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(0))

			Expect(stdout).To(gbytes.Say(line))
		},
			Entry("deployment name", "\nname: grinder\n"),
			Entry("stemcell version", "\n  version: stemcell\\.version\\.gcp\n"),
			Entry("network name", "\n  networks:\n  - name: network-name\n"),
			Entry("kubernetes API URL", "\n      kubernetes-api-url: https://12\\.23\\.34\\.45:101928\n"),
			Entry("Auto-generated kubelet password", "\n      kubelet-password: \\(\\(kubelet-password\\)\\)\n"),
			Entry("Auto-generated admin password", "\n      admin-password: \\(\\(kubo-admin-password\\)\\)\n"),
		)

		It("should include a variable section with tls-kubelet, tls-kubernetes", func() {
			status, err := bash.Run("main", []string{kuboEnv, "cucumber"})
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(0))

			Expect(stdout).To(gbytes.Say("variables:"))
		})

		It("should reproduce the same manifest on the second run", func() {
			bash.Run("main", []string{kuboEnv, "fort"})

			firstRun := make([]byte, len(stdout.Contents()))
			secondRun := make([]byte, len(stdout.Contents()))

			_, err := stdout.Read(firstRun)
			Expect(err).NotTo(HaveOccurred())

			bash.Run("main", []string{kuboEnv, "fort"})

			_, err = stdout.Read(secondRun)
			Expect(err).NotTo(HaveOccurred())

			Expect(firstRun).To(Equal(secondRun))
		})

		It("generates a manifest without the secrets", func() {
			secretlessEnv := filepath.Join(testEnvironmentPath, "secretless")
			status, err := bash.Run("main", []string{secretlessEnv, "sensors"})
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(0))

			Expect(stdout).To(gbytes.Say("\n        uaa_client_secret: \\(\\(routing-cf-client-secret\\)\\)\n"))
		})

		It("uses ops-files to modify the manifest", func() {
			opsfileEnv := filepath.Join(testEnvironmentPath, "with_ops")
			status, err := bash.Run("main", []string{opsfileEnv, "name"})

			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(0))
			Expect(stdout).To(gbytes.Say("\n  os: MALARIA0\n"))
		})

		It("applies http proxy settings if they exist", func() {
			opsfileEnv := filepath.Join(testEnvironmentPath, "with_http_proxy")
			status, err := bash.Run("main", []string{opsfileEnv, "name"})

			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(0))
			Expect(stdout).To(gbytes.Say("\n        http_proxy: my.proxy.com\n"))
		})

		It("applies https proxy settings if they exist", func() {
			opsfileEnv := filepath.Join(testEnvironmentPath, "with_https_proxy")
			status, err := bash.Run("main", []string{opsfileEnv, "name"})

			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(0))
			Expect(stdout).To(gbytes.Say("\n        https_proxy: my.sslproxy.com\n"))
		})

		It("applies http no_proxy settings if they exist", func() {
			opsfileEnv := filepath.Join(testEnvironmentPath, "with_no_proxy")
			status, err := bash.Run("main", []string{opsfileEnv, "name"})

			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(0))
			Expect(stdout).To(gbytes.Say("\n        no_proxy: dont.proxy.me\n"))
		})

		It("uses vars-files to modify the manifest", func() {
			opsfileEnv := filepath.Join(testEnvironmentPath, "with_vars")
			status, err := bash.Run("main", []string{opsfileEnv, "name"})

			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(0))
			Expect(stdout).To(gbytes.Say("\n      kubelet-password: Shields up, ancient life!\n"))
		})
	})

	It("expands the bosh environment path to absolute value", func() {
		command := exec.Command("./generate_kubo_manifest", "../src/kubo-deployment-tests/resources/environments/test_gcp", "name")
		command.Stdout = bash.Stdout
		command.Stderr = bash.Stderr
		command.Dir = pathToScript("")
		Expect(command.Run()).To(Succeed())
	})

	It("runs from any location", func() {
		command := exec.Command("./bin/generate_kubo_manifest", "src/kubo-deployment-tests/resources/environments/test_gcp", "name")
		command.Stdout = bash.Stdout
		command.Stderr = bash.Stderr
		command.Dir = pathFromRoot("")
		Expect(command.Run()).To(Succeed())
	})

	It("should generate a valid manifest", func() {
		files, _ := filepath.Glob(testEnvironmentPath + "/*")
		for _, env := range files {
			command := exec.Command("./bin/generate_kubo_manifest", env, "env-name")
			out := gbytes.NewBuffer()
			command.Stdout = out
			command.Dir = pathFromRoot("")
			Expect(command.Run()).To(Succeed())

			var output map[string]interface{}
			Expect(yaml.Unmarshal(out.Contents(), &output)).To(
				Succeed(), fmt.Sprintf("Could not generate manifest for %s %s", env, string(stdout.Contents())))
		}
	})

	It("should not write anything to stderr", func() {
		files, _ := filepath.Glob(testEnvironmentPath + "/*")
		for _, env := range files {
			command := exec.Command("./bin/generate_kubo_manifest", env, "env-name")
			errBuffer := gbytes.NewBuffer()
			command.Stdout = GinkgoWriter
			command.Stderr = errBuffer
			command.Dir = pathFromRoot("")
			Expect(command.Run()).To(Succeed())
			Expect(string(errBuffer.Contents())).To(HaveLen(0))
		}
	})
})
