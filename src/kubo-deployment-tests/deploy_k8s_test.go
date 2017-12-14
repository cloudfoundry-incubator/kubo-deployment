package kubo_deployment_tests_test

import (
	"path"

	. "github.com/jhvhs/gob-mock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Deploy K8s", func() {
	validGcpEnvironment := path.Join(testEnvironmentPath, "test_gcp")
	validvSphereEnvironment := path.Join(testEnvironmentPath, "test_vsphere")
	validOpenstackEnvironment := path.Join(testEnvironmentPath, "test_openstack")
	validAWSEnvironment := path.Join(testEnvironmentPath, "test_aws")

	BeforeEach(func() {
		bash.Source(pathToScript("deploy_k8s"), nil)
		boshMock := MockOrCallThrough("bosh-cli", `echo -n "3124.12"`, `[ "$1" == 'int' ]`)
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
			Expect(stderr).To(gbytes.Say("bosh-cli upload-release kubo-release.tgz"))
		})
	})

	Context("When release source is skip", func() {
		It("doesn't upload any release", func() {
			code, err := bash.Run("main", []string{validGcpEnvironment, "deployment", "skip"})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))
			Expect(stderr).NotTo(gbytes.Say("bosh-cli upload-release"))
		})

		It("does not upload the stemcell", func() {
			code, err := bash.Run("main", []string{validGcpEnvironment, "deployment", "skip"})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))

			Expect(stderr).NotTo(gbytes.Say("bosh-cli upload-stemcell"))
		})
	})

	Context("When release source is local", func() {
		It("uploads the local release", func() {
			code, err := bash.Run("main", []string{validGcpEnvironment, "deployment", "local"})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))
			Expect(stderr).To(gbytes.Say("bosh-cli upload-release kubo-release.tgz"))
		})

		It("uploads the stemcell for GCP", func() {
			boshMock := MockOrCallThrough("bosh-cli", `echo -n "3124.12"`, `[ "$1" == 'int' ] && [ ! "$4" == '/stemcells/0/version' ] `)
			ApplyMocks(bash, []Gob{boshMock})

			code, err := bash.Run("main", []string{validGcpEnvironment, "deployment", "local"})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))

			Expect(stderr).To(gbytes.Say("bosh-cli upload-stemcell https://s3.amazonaws.com/bosh-gce-light-stemcells/light-bosh-stemcell-3124.12-google-kvm-ubuntu-trusty-go_agent.tgz"))
		})

		It("uploads the stemcell for vSphere", func() {
			boshMock := MockOrCallThrough("bosh-cli", `echo -n "3124.12"`, `[ "$1" == 'int' ] && [ ! "$4" == '/stemcells/0/version' ] `)
			ApplyMocks(bash, []Gob{boshMock})
			code, err := bash.Run("main", []string{validvSphereEnvironment, "deployment", "local"})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))

			Expect(stderr).To(gbytes.Say("bosh-cli upload-stemcell https://s3.amazonaws.com/bosh-core-stemcells/vsphere/bosh-stemcell-3124.12-vsphere-esxi-ubuntu-trusty-go_agent.tgz"))
		})

		It("uploads the stemcell for OpenStack", func() {
			boshMock := MockOrCallThrough("bosh-cli", `echo -n "3124.12"`, `[ "$1" == 'int' ] && [ ! "$4" == '/stemcells/0/version' ] `)
			ApplyMocks(bash, []Gob{boshMock})
			code, err := bash.Run("main", []string{validOpenstackEnvironment, "deployment", "local"})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))

			Expect(stderr).To(gbytes.Say("bosh-cli upload-stemcell https://s3.amazonaws.com/bosh-core-stemcells/openstack/bosh-stemcell-3124.12-openstack-kvm-ubuntu-trusty-go_agent.tgz"))
		})

		It("uploads the stemcell for AWS", func() {
			boshMock := MockOrCallThrough("bosh-cli", `echo -n "3124.12"`, `[ "$1" == 'int' ] && [ ! "$4" == '/stemcells/0/version' ] `)
			ApplyMocks(bash, []Gob{boshMock})
			code, err := bash.Run("main", []string{validAWSEnvironment, "deployment", "local"})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))

			Expect(stderr).To(gbytes.Say("bosh-cli upload-stemcell https://s3.amazonaws.com/bosh-aws-light-stemcells/light-bosh-stemcell-3124.12-aws-xen-hvm-ubuntu-trusty-go_agent.tgz"))
		})
	})

	Context("When release source is dev", func() {
		It("should create and upload a release", func() {
			createAndUploadReleaseMock := Mock("create_and_upload_release", `echo`)
			exportBoshEnvironmentMock := Mock("export_bosh_environment", "echo")
			cloudConfigMock := Mock("set_cloud_config", "echo")
			deployToBoshMock := Mock("deploy_to_bosh", "echo")
			depsMock := Mock("get_deps", "echo")
			ApplyMocks(bash, []Gob{createAndUploadReleaseMock, exportBoshEnvironmentMock, cloudConfigMock, deployToBoshMock, depsMock})

			code, err := bash.Run("main", []string{validGcpEnvironment, "deployment", "dev"})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))
			Expect(stderr).To(gbytes.Say("create_and_upload_release"))
		})
	})

	Context("When release source is public", func() {
		It("should upload a release", func() {
			uploadReleaseMock := Mock("upload_release", `echo -n $@`)
			getSettingMock := Mock("get_setting", `echo url-to-public-release`)
			exportBoshEnvironmentMock := Mock("export_bosh_environment", "echo")
			cloudConfigMock := Mock("set_cloud_config", "echo")
			deployToBoshMock := Mock("deploy_to_bosh", "echo")
			depsMock := Mock("get_deps", "echo")
			ApplyMocks(bash, []Gob{uploadReleaseMock, getSettingMock, exportBoshEnvironmentMock, cloudConfigMock, deployToBoshMock, depsMock})

			code, err := bash.Run("main", []string{validGcpEnvironment, "deployment", "public"})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))
			Expect(stderr).To(gbytes.Say("upload_release url-to-public-release"))
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
			boshMock := MockOrCallThrough("bosh-cli", `echo -n "0"`, `! [[ "$4" =~ 'apply-specs' ]]`)
			ApplyMocks(bash, []Gob{boshMock})

			code, err := bash.Run("main", []string{validGcpEnvironment, "deployment", "skip"})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))
			Expect(stderr).To(gbytes.Say("run-errand apply-specs"))
		})
	})
})

var _ = Describe("upload_stemcell", func() {
	BeforeEach(func() {
		bash.Source(pathToScript("deploy_k8s"), nil)
		boshMock := Mock("bosh-cli", "echo")
		ApplyMocks(bash, []Gob{boshMock})
	})

	Context("GCP", func() {
		It("should upload a google kvm ubuntu trusty stemcell", func() {
			code, err := bash.Run("upload_stemcell", []string{"public", "gcp"})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))
			Expect(stderr).To(gbytes.Say("bosh-cli upload-stemcell"))
			Expect(stderr).To(gbytes.Say("google-kvm-ubuntu-trusty-go_agent"))
		})
	})

	Context("vSphere", func() {
		It("should upload a vSphere esxi ubuntu trusty stemcell", func() {
			code, err := bash.Run("upload_stemcell", []string{"public", "vsphere"})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))
			Expect(stderr).To(gbytes.Say("bosh-cli upload-stemcell"))
			Expect(stderr).To(gbytes.Say("vsphere-esxi-ubuntu-trusty-go_agent"))
		})
	})

	Context("OpenStack", func() {
		It("should upload an OpenStack kvm ubuntu trusty stemcell", func() {
			code, err := bash.Run("upload_stemcell", []string{"public", "openstack"})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))
			Expect(stderr).To(gbytes.Say("bosh-cli upload-stemcell"))
			Expect(stderr).To(gbytes.Say("openstack-kvm-ubuntu-trusty-go_agent"))
		})
	})

	Context("AWS", func() {
		It("should upload an AWS xen hvm ubuntu trusty light stemcell", func() {
			code, err := bash.Run("upload_stemcell", []string{"public", "aws"})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))
			Expect(stderr).To(gbytes.Say("bosh-aws-light-stemcells"))
			Expect(stderr).To(gbytes.Say("aws-xen-hvm-ubuntu-trusty-go_agent"))
		})
	})

	Context("skip", func() {
		It("should not upload a stemcell", func() {
			code, err := bash.Run("upload_stemcell", []string{"skip", "aws"})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))
			Expect(stderr).NotTo(gbytes.Say("bosh-cli upload-stemcell"))
		})
	})
})

var _ = Describe("get_director_uuid", func() {
	It("should return UUID from bosh env command", func() {
		bash.Source(pathToScript("deploy_k8s"), nil)
		boshMock := MockOrCallThrough("bosh-cli", `echo -n \
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
