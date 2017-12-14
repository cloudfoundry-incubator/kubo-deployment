package kubo_deployment_tests_test

import (
	"fmt"

	. "github.com/jhvhs/gob-mock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/ginkgo/extensions/table"
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

	Describe("generate_manifest", func() {
		BeforeEach(func() {
			bash.Source(pathToScript("lib/deploy_utils"), nil)

		})

		It("applies dev, bootstrap and use-runtime-config-bosh-dns ops files", func() {

			boshMock := Mock("bosh-cli", `
			if [[ "$3" =~ "addons_spec_path" \
				|| "$3" =~ "http_proxy" \
				|| "$3" =~ "https_proxy" \
				|| "$3" =~ "no_proxy" ]]; then
				return 1
			elif [[ "$3" =~ "routing_mode" ]]; then
				echo "the-routing-mode"
			elif [[ "$3" =~ "iaas" ]]; then
				echo "the-iaas"
			else
				echo
			fi`)
			ApplyMocks(bash, []Gob{boshMock})

			code, err := bash.Run("generate_manifest", []string{"environment-path", "deployment-name", "non-existent-manifest-path", "director-uuid"})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))
			Expect(stderr).To(gbytes.Say("routing_mode"))
			Expect(stderr).To(gbytes.Say("iaas"))
			Expect(stderr).To(gbytes.Say("misc/dev.yml"))
			Expect(stderr).To(gbytes.Say("misc/bootstrap.yml"))
			Expect(stderr).To(gbytes.Say("use-runtime-config-bosh-dns.yml"))
			Expect(stderr).To(gbytes.Say("--var=deployment_name=\"deployment-name\""))
		})

		Context("when http_proxy is set", func() {
			It("applies add-http-proxy ops file", func() {
				boshMock := Mock("bosh-cli", `
				if [[ "$3" =~ "addons_spec_path" \
					|| "$3" =~ "https_proxy" \
					|| "$3" =~ "no_proxy" ]]; then
					return 1
				elif [[ "$3" =~ "routing_mode" ]]; then
					echo "the-routing-mode"
				elif [[ "$3" =~ "iaas" ]]; then
					echo "the-iaas"
				else
					echo
				fi`)
				ApplyMocks(bash, []Gob{boshMock})

				code, err := bash.Run("generate_manifest", []string{"environment-path", "deployment-name", "non-existent-manifest-path", "director-uuid"})
				Expect(err).NotTo(HaveOccurred())
				Expect(code).To(Equal(0))
				Expect(stderr).To(gbytes.Say("add-http-proxy.yml"))
			})
		})


		Context("when https_proxy is set", func() {
			It("applies add-https-proxy ops file", func() {
				boshMock := Mock("bosh-cli", `
				if [[ "$3" =~ "addons_spec_path" \
					|| "$3" =~ "http_proxy" \
					|| "$3" =~ "no_proxy" ]]; then
					return 1
				elif [[ "$3" =~ "routing_mode" ]]; then
					echo "the-routing-mode"
				elif [[ "$3" =~ "iaas" ]]; then
					echo "the-iaas"
				else
					echo
				fi`)
				ApplyMocks(bash, []Gob{boshMock})

				code, err := bash.Run("generate_manifest", []string{"environment-path", "deployment-name", "non-existent-manifest-path", "director-uuid"})
				Expect(err).NotTo(HaveOccurred())
				Expect(code).To(Equal(0))
				Expect(stderr).To(gbytes.Say("add-https-proxy.yml"))
			})
		})

		Context("when no_proxy is set", func() {
			It("applies no-proxy ops-file", func() {
				boshMock := Mock("bosh-cli", `
				if [[ "$3" =~ "addons_spec_path" \
					|| "$3" =~ "http_proxy" \
					|| "$3" =~ "https_proxy" ]]; then
					return 1
				elif [[ "$3" =~ "routing_mode" ]]; then
					echo "the-routing-mode"
				elif [[ "$3" =~ "iaas" ]]; then
					echo "the-iaas"
				else
					echo
				fi`)
				ApplyMocks(bash, []Gob{boshMock})

				code, err := bash.Run("generate_manifest", []string{"environment-path", "deployment-name", "non-existent-manifest-path", "director-uuid"})
				Expect(err).NotTo(HaveOccurred())
				Expect(code).To(Equal(0))
				Expect(stderr).To(gbytes.Say("add-no-proxy.yml"))
			})
		})

		Context("when routing_mode is cf", func() {
			It("applies cf-routing ops-file", func() {
				boshMock := Mock("bosh-cli", `
				if [[ "$3" =~ "addons_spec_path" \
					|| "$3" =~ "http_proxy" \
					|| "$3" =~ "https_proxy" \
					|| "$3" =~ "no_proxy" ]]; then
					return 1
				elif [[ "$3" =~ "routing_mode" ]]; then
					echo "cf"
				elif [[ "$3" =~ "iaas" ]]; then
					echo "the-iaas"
				else
					echo
				fi`)
				ApplyMocks(bash, []Gob{boshMock})

				code, err := bash.Run("generate_manifest", []string{"environment-path", "deployment-name", "non-existent-manifest-path", "director-uuid"})
				Expect(err).NotTo(HaveOccurred())
				Expect(code).To(Equal(0))
				Expect(stderr).To(gbytes.Say("cf-routing.yml"))
			})
		})

		Context("when iaas is aws", func() {
			It("applies aws lb ops-file", func() {
				boshMock := Mock("bosh-cli", `
				if [[ "$3" =~ "addons_spec_path" \
					|| "$3" =~ "http_proxy" \
					|| "$3" =~ "https_proxy" \
					|| "$3" =~ "no_proxy" ]]; then
					return 1
				elif [[ "$3" =~ "routing_mode" ]]; then
					echo "the-routing-mode"
				elif [[ "$3" =~ "iaas" ]]; then
					echo "aws"
				else
					echo
				fi`)
				ApplyMocks(bash, []Gob{boshMock})

				code, err := bash.Run("generate_manifest", []string{"environment-path", "deployment-name", "non-existent-manifest-path", "director-uuid"})
				Expect(err).NotTo(HaveOccurred())
				Expect(code).To(Equal(0))
				Expect(stderr).To(gbytes.Say("aws/lb.yml"))
			})
		})

		DescribeTable("applies the cloud_provider ops-file", func(iaas string) {
			boshMock := Mock("bosh-cli", fmt.Sprintf(`
			if [[ "$3" =~ "addons_spec_path" \
				|| "$3" =~ "http_proxy" \
				|| "$3" =~ "https_proxy" \
				|| "$3" =~ "no_proxy" ]]; then
				return 1
			elif [[ "$3" =~ "routing_mode" ]]; then
				echo "the-routing-mode"
			elif [[ "$3" =~ "iaas" ]]; then
				echo "%s"
			else
				echo
			fi`, iaas))
			ApplyMocks(bash, []Gob{boshMock})

			code, err := bash.Run("generate_manifest", []string{"environment-path", "deployment-name", pathFromRoot("manifests/cfcr.yml"), "director-uuid"})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))
			Expect(stderr).To(gbytes.Say(fmt.Sprintf("%s/cloud-provider.yml", iaas)))
		},
			Entry("when the iaas is aws", "aws"),
			Entry("when the iaas is gcp", "gcp"),
			Entry("when the iaas is vsphere", "vsphere"),
		)

		Context("when authorization_mode is not set", func() {
			It("sets authorization_mode variable to abac", func(){
				boshMock := Mock("bosh-cli", `
				if [[ "$3" =~ "addons_spec_path" \
					|| "$3" =~ "http_proxy" \
					|| "$3" =~ "https_proxy" \
					|| "$3" =~ "no_proxy" \
					|| "$3" =~ "authorization_mode" ]]; then
					return 1
				elif [[ "$3" =~ "routing_mode" ]]; then
					echo "the-routing-mode"
				elif [[ "$3" =~ "iaas" ]]; then
					echo "the-iaas"
				else
					echo
				fi`)
				ApplyMocks(bash, []Gob{boshMock})

				code, err := bash.Run("generate_manifest", []string{"environment-path", "deployment-name", "non-existent-manifest-path", "director-uuid"})
				Expect(err).NotTo(HaveOccurred())
				Expect(code).To(Equal(0))
				Expect(stderr).To(gbytes.Say("var authorization_mode=abac"))
			})
		})

		Context("when worker_count is not set", func() {
			It("sets worker_count variable to 3", func() {
				boshMock := Mock("bosh-cli", `
				if [[ "$3" =~ "addons_spec_path" \
					|| "$3" =~ "http_proxy" \
					|| "$3" =~ "https_proxy" \
					|| "$3" =~ "no_proxy" \
					|| "$3" =~ "worker_count" ]]; then
					return 1
				elif [[ "$3" =~ "routing_mode" ]]; then
					echo "the-routing-mode"
				elif [[ "$3" =~ "iaas" ]]; then
					echo "the-iaas"
				else
					echo
				fi`)
				ApplyMocks(bash, []Gob{boshMock})

				code, err := bash.Run("generate_manifest", []string{"environment-path", "deployment-name", "non-existent-manifest-path", "director-uuid"})
				Expect(err).NotTo(HaveOccurred())
				Expect(code).To(Equal(0))
				Expect(stderr).To(gbytes.Say("var worker_count=3"))
			})
		})

		Context("when iaas is gcp", func() {
			Context("when service_account_worker is not set", func(){
				It("applies the add-service-key-worker ops-file", func() {
					boshMock := Mock("bosh-cli", `
						if [[ "$3" =~ "addons_spec_path" \
							|| "$3" =~ "http_proxy" \
							|| "$3" =~ "https_proxy" \
							|| "$3" =~ "no_proxy" \
							|| "$3" =~ "service_account_worker" ]]; then
							return 1
						elif [[ "$3" =~ "routing_mode" ]]; then
							echo "the-routing-mode"
						elif [[ "$3" =~ "iaas" ]]; then
							echo "gcp"
						else
							echo
						fi`)
					ApplyMocks(bash, []Gob{boshMock})

					code, err := bash.Run("generate_manifest", []string{"environment-path", "deployment-name", "non-existent-manifest-path", "director-uuid"})
					Expect(err).NotTo(HaveOccurred())
					Expect(code).To(Equal(0))
					Expect(stderr).To(gbytes.Say("gcp/add-service-key-worker.yml"))
				})
			})

			Context("when service_account_master is not set", func(){
				It("applies the add-service-key-master ops-file", func() {
					boshMock := Mock("bosh-cli", `
						if [[ "$3" =~ "addons_spec_path" \
							|| "$3" =~ "http_proxy" \
							|| "$3" =~ "https_proxy" \
							|| "$3" =~ "no_proxy" \
							|| "$3" =~ "service_account_master" ]]; then
							return 1
						elif [[ "$3" =~ "routing_mode" ]]; then
							echo "the-routing-mode"
						elif [[ "$3" =~ "iaas" ]]; then
							echo "gcp"
						else
							echo
						fi`)
					ApplyMocks(bash, []Gob{boshMock})

					code, err := bash.Run("generate_manifest", []string{"environment-path", "deployment-name", "non-existent-manifest-path", "director-uuid"})
					Expect(err).NotTo(HaveOccurred())
					Expect(code).To(Equal(0))
					Expect(stderr).To(gbytes.Say("gcp/add-service-key-master.yml"))
				})
			})
		})

		Context("when creds.yml exists", func() {
			It("applies the creds.yml vars-file", func() {
				boshMock := Mock("bosh-cli", `
					if [[ "$3" =~ "addons_spec_path" \
						|| "$3" =~ "http_proxy" \
						|| "$3" =~ "https_proxy" \
						|| "$3" =~ "no_proxy" ]]; then
						return 1
					elif [[ "$3" =~ "routing_mode" ]]; then
						echo "the-routing-mode"
					elif [[ "$3" =~ "iaas" ]]; then
						echo "the-iaas"
					else
						echo
					fi`)
				ApplyMocks(bash, []Gob{boshMock})

				code, err := bash.Run("generate_manifest", []string{path.Join(testEnvironmentPath, "test_gcp_with_creds"), "deployment-name", "non-existent-manifest-path", "director-uuid"})
				Expect(err).NotTo(HaveOccurred())
				Expect(code).To(Equal(0))
				Expect(stderr).To(gbytes.Say("creds.yml"))
			})
		})

		Context("when director-secrets.yml exists", func() {
			It("applies the director-secrets.yml vars-file", func() {
				boshMock := Mock("bosh-cli", `
					if [[ "$3" =~ "addons_spec_path" \
						|| "$3" =~ "http_proxy" \
						|| "$3" =~ "https_proxy" \
						|| "$3" =~ "no_proxy" ]]; then
						return 1
					elif [[ "$3" =~ "routing_mode" ]]; then
						echo "the-routing-mode"
					elif [[ "$3" =~ "iaas" ]]; then
						echo "the-iaas"
					else
						echo
					fi`)
				ApplyMocks(bash, []Gob{boshMock})

				code, err := bash.Run("generate_manifest", []string{path.Join(testEnvironmentPath, "test_gcp_with_creds"), "deployment-name", "non-existent-manifest-path", "director-uuid"})
				Expect(err).NotTo(HaveOccurred())
				Expect(code).To(Equal(0))
				Expect(stderr).To(gbytes.Say("director-secrets.yml"))
			})
		})

		Context("when addons_spec_path exists", func() {
			It("applies addons-spec.yml ops-file and addon_path as vars file", func(){
				boshMock := Mock("bosh-cli", `
					if [[ "$3" =~ "http_proxy" \
						|| "$3" =~ "https_proxy" \
						|| "$3" =~ "no_proxy" ]]; then
						return 1
					elif [[ "$3" =~ "routing_mode" ]]; then
						echo "the-routing-mode"
					elif [[ "$3" =~ "iaas" ]]; then
						echo "the-iaas"
					elif [[ "$3" =~ "addons_spec_path" ]]; then
						echo "addon.yml"
					else
						echo
					fi`)
				ApplyMocks(bash, []Gob{boshMock})

				code, err := bash.Run("generate_manifest", []string{path.Join(testEnvironmentPath, "with_addons"), "deployment-name", pathFromRoot("manifests/cfcr.yml"), "director-uuid"})
				Expect(err).NotTo(HaveOccurred())
				Expect(code).To(Equal(0))
				Expect(stderr).To(gbytes.Say("addons-spec.yml"))
				Expect(stderr).To(gbytes.Say("var-file=\"addons-spec"))
				Expect(stderr).To(gbytes.Say("addon.yml"))
			})
		})
	})
})
