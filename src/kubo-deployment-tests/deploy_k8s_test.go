package kubo_deployment_tests_test

import (
	"fmt"
	"path"

	. "github.com/jhvhs/gob-mock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Deploy K8s", func() {
	validGcpEnvironment := path.Join(testEnvironmentPath, "test_gcp_with_creds")

	BeforeEach(func() {
		bash.Source(pathToScript("deploy_k8s"), nil)
		boshMock := MockOrCallThrough("bosh", `echo -n "3124.12"`, `[ "$1" == 'int' ]`)
		getDirectorUUIDMock := Mock("get_director_uuid", `echo -n "director-uuid"`)
		ApplyMocks(bash, []Gob{boshMock, getDirectorUUIDMock})
	})

	JustBeforeEach(func() {
		bash.Source("", func(string) ([]byte, error) {
			return repoDirectoryFunction, nil
		})
	})

	Context("When release source is empty", func() {
		It("deploys with local release", func() {
			code, err := bash.Run("main", []string{validGcpEnvironment, "deployment"})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))
			Expect(stderr).To(gbytes.Say("bosh upload-release kubo-release.tgz"))
		})
	})

	Context("Artefact upload", func() {
		BeforeEach(func() {
			uploadMock := Spy("upload_artefacts")
			ApplyMocks(bash, []Gob{uploadMock})
		})

		DescribeTable("valid upload sources", func(source string) {
			code, err := bash.Run("main", []string{validGcpEnvironment, "deployment", source})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))
			uploadInvocation := fmt.Sprintf("upload_artefacts %s %s", validGcpEnvironment, source)
			Expect(stderr).To(gbytes.Say(uploadInvocation))
		},
			Entry("local", "local"),
			Entry("dev", "dev"),
			Entry("public", "public"),
		)

		Context("when skip is given", func() {
			It("skips the upload", func() {
				code, err := bash.Run("main", []string{validGcpEnvironment, "deployment", "skip"})
				Expect(err).NotTo(HaveOccurred())
				Expect(code).To(Equal(0))
				Expect(stderr).NotTo(gbytes.Say("upload_artefacts"))

			})
		})

	})

	It("Should export bosh env, set cloud config and deploy", func() {
		cloudConfigMock := Mock("set_cloud_config", "echo")
		exportBoshEnvironmentMock := Mock("export_bosh_environment", "echo")
		deployToBoshMock := Mock("deploy_to_bosh", "echo")
		ApplyMocks(bash, []Gob{cloudConfigMock, exportBoshEnvironmentMock, deployToBoshMock})

		depsMock := Mock("get_deps", "echo")
		ApplyMocks(bash, []Gob{depsMock})

		code, err := bash.Run("main", []string{validGcpEnvironment, "deployment", "skip"})
		Expect(err).NotTo(HaveOccurred())
		Expect(code).To(Equal(0))
		Expect(stderr).To(gbytes.Say("export_bosh_environment"))
		Expect(stderr).To(gbytes.Say("set_cloud_config"))
		Expect(stderr).To(gbytes.Say("deploy_to_bosh"))
	})

	Context("When apply-specs is present in the manifest", func() {
		It("should run apply-specs errand", func() {
			cloudConfigMock := Mock("set_cloud_config", "echo")
			exportBoshEnvironmentMock := Mock("export_bosh_environment", "echo")
			deployToBoshMock := Mock("deploy_to_bosh", "echo")
			ApplyMocks(bash, []Gob{cloudConfigMock, exportBoshEnvironmentMock, deployToBoshMock})

			depsMock := Mock("get_deps", "echo")
			ApplyMocks(bash, []Gob{depsMock})
			boshMock := MockOrCallThrough("bosh", `echo -n "0"`, `! [[ "$4" =~ 'apply-specs' ]]`)
			ApplyMocks(bash, []Gob{boshMock})

			code, err := bash.Run("main", []string{validGcpEnvironment, "deployment", "skip"})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))
			Expect(stderr).To(gbytes.Say("run-errand apply-specs"))
		})
	})
})

var _ = Describe("get_director_uuid", func() {
	It("should return UUID from bosh env command", func() {
		bash.Source(pathToScript("deploy_k8s"), nil)
		boshMock := MockOrCallThrough("bosh", `echo -n \
'{
  "Tables": [
			{
					"Content": "",
					"Header": {
							"cpi": "CPI",
							"features": "Features",
							"name": "Name",
							"user": "User",
							"uuid": "UUID",
							"version": "Version"
					},
					"Rows": [
							{
									"cpi": "google_cpi",
									"features": "compiled_package_cache: disabled\nconfig_server: enabled\ndns: enabled\nsnapshots: disabled",
									"name": "I AM CI",
									"user": "admin",
									"uuid": "director-uuid",
									"version": "262.0.0 (00000000)"
							}
					],
					"Notes": null
			}
	],
	"Blocks": null,
	"Lines": [
			"Using environment '10.0.250.252' as user 'admin' (openid, bosh.admin)",
			"Succeeded"
	]
}
'`, `[ "$1" != 'environment' ]`)
		ApplyMocks(bash, []Gob{boshMock})

		code, err := bash.Run("get_director_uuid", []string{})

		Expect(err).NotTo(HaveOccurred())
		Expect(code).To(Equal(0))
		Expect(stdout).To(gbytes.Say("director-uuid"))
	})
})
