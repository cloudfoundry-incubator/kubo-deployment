package kubo_deployment_tests_test

import (
	"fmt"

	. "github.com/jhvhs/gob-mock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"path"
)

var _ = Describe("DeployUtils", func() {
	Describe("set_cloud_config", func() {
		It("Should generate cloud config", func() {
			bash.Source(pathToScript("lib/deploy_utils"), nil)

			bash.Export("BOSH_ENV", "kubo-env")
			bash.Source("", func(string) ([]byte, error) {
				return []byte(fmt.Sprintf("export PATH=%s:$PATH", pathFromRoot("bin"))), nil
			})

			generateCloudConfigMock := Mock("generate_cloud_config", `echo -n "cc"`)
			boshCliMock := Mock("bosh-cli", `echo -n "$@"`)
			ApplyMocks(bash, []Gob{generateCloudConfigMock, boshCliMock})

			code, err := bash.Run("set_cloud_config", []string{})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))
			Expect(stderr).To(gbytes.Say("generate_cloud_config kubo-env"))
		})

		It("Should update cloud config", func() {
			bash.Source(pathToScript("lib/deploy_utils"), nil)

			bash.Export("BOSH_ENV", "kubo-env")
			bash.Export("BOSH_NAME", "env-name")
			bash.Source("", func(string) ([]byte, error) {
				return []byte(fmt.Sprintf("export PATH=%s:$PATH", pathFromRoot("bin"))), nil
			})

			generateCloudConfigMock := Mock("generate_cloud_config", `echo -n "cc"`)
			boshCliMock := Mock("bosh-cli", `echo -n "$@"`)
			ApplyMocks(bash, []Gob{generateCloudConfigMock, boshCliMock})

			code, err := bash.Run("set_cloud_config", []string{})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))
			Expect(stderr).To(gbytes.Say("bosh-cli -n -e env-name update-cloud-config -"))
		})
	})

	Describe("export_bosh_environment", func() {
		It("should set BOSH_ENV and BOSH_NAME", func() {
			bash.Source(pathToScript("lib/deploy_utils"), nil)
			code, err := bash.Run("export_bosh_environment /envs/foo && echo $BOSH_ENV && echo $BOSH_NAME", []string{})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))
			Expect(stdout).To(gbytes.Say("/envs/foo"))
			Expect(stdout).To(gbytes.Say("foo"))
		})
	})

	Describe("deploy_to_bosh", func() {
		It("should deploy to bosh", func() {
			bash.Source(pathToScript("lib/deploy_utils"), nil)

			getBoshSecretMock := Mock("get_bosh_secret", `echo "the-secret"`)
			boshCliMock := Mock("bosh-cli", `echo -n "$@"`)
			ApplyMocks(bash, []Gob{ getBoshSecretMock, boshCliMock })

			bash.Export("BOSH_ENV", "kubo-env")
			bash.Export("BOSH_NAME", "env-name")

			code, err := bash.Run("deploy_to_bosh", []string{"manifest", "deployment-name"})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))
			Expect(stderr).To(gbytes.Say("bosh-cli -d deployment-name -n deploy --no-redact -"))
		})
	})

	Describe("get_bosh_secret", func() {
		It("should get bosh_admin_client_secret setting", func() {
			bash.Source(pathToScript("lib/deploy_utils"), nil)

			getSettingMock := Mock("get_setting", `echo "the-secret"`)
			ApplyMocks(bash, []Gob{ getSettingMock })

			code, err := bash.Run("get_bosh_secret", []string{})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))
			Expect(stdout).To(gbytes.Say("the-secret"))
		})
	})

	Describe("get_setting", func() {
		It("should call bosh interpolate", func() {
			bash.Source(pathToScript("lib/deploy_utils"), nil)

			boshCliMock := Mock("bosh-cli", `echo -n "$@"`)
			ApplyMocks(bash, []Gob{ boshCliMock })

			bash.Export("BOSH_ENV", "kubo-env")

			code, err := bash.Run("get_setting", []string{"fileToQuery.yml", "path/subpath"})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))
			Expect(stderr).To(gbytes.Say("bosh-cli int kubo-env/fileToQuery.yml --path path/subpath"))
		})

		It("should return the setting value", func() {
			bash.Source(pathToScript("lib/deploy_utils"), nil)

			boshCliMock := Mock("bosh-cli", `echo "value-at-path"`)
			ApplyMocks(bash, []Gob{ boshCliMock })

			code, err := bash.Run("get_setting", []string{"fileToQuery.yml", "path/subpath"})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))
			Expect(stdout).To(gbytes.Say("value-at-path"))
		})
	})

	Describe("create_and_upload_release", func(){
		Context("is called with a valid directory path", func(){
			It("should create a release and upload it", func(){
				bash.Source(pathToScript("lib/deploy_utils"), nil)

				getBoshSecretMock := Mock("get_bosh_secret", `echo "the-secret"`)
				boshCliMock := Mock("bosh-cli", `echo -n "$@"`)
				uploadReleaseMock := Mock("upload_release", `echo`)
				ApplyMocks(bash, []Gob{ getBoshSecretMock, boshCliMock, uploadReleaseMock})


				releasePath := path.Join(resourcesPath, "releases/mock-release")

				code, err := bash.Run("create_and_upload_release", []string{releasePath})
				Expect(err).NotTo(HaveOccurred())
				Expect(code).To(Equal(0))
				Expect(stderr).To(gbytes.Say("bosh-cli create-release --force --name mock"))
				Expect(stderr).To(gbytes.Say("upload_release --name=mock"))
			})
		})

		Context("is called with an invalid argument", func(){
			It("should exit", func(){
				bash.Source(pathToScript("lib/deploy_utils"), nil)
				code, err := bash.Run("create_and_upload_release", []string{"path_does_not_exist"})
				Expect(err).NotTo(HaveOccurred())
				Expect(code).To(Equal(1))
			})
		})
	})

	Describe("upload_release", func(){
		It("should upload release", func() {
			bash.Source(pathToScript("lib/deploy_utils"), nil)

			getBoshSecretMock := Mock("get_bosh_secret", `echo "the-secret"`)
			boshCliMock := Mock("bosh-cli", `echo -n "$@"`)
			ApplyMocks(bash, []Gob{getBoshSecretMock, boshCliMock})

			bash.Export("BOSH_ENV", "kubo-env")
			bash.Export("BOSH_NAME", "env-name")

			code, err := bash.Run("upload_release", []string{"release-name"})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))
			Expect(stderr).To(gbytes.Say("bosh-cli upload-release release-name"))
		})
	})
})
