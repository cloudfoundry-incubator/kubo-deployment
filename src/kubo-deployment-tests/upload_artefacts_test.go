package kubo_deployment_tests_test

import (
	"path"

	. "github.com/jhvhs/gob-mock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("upload_artefacts", func() {
	validGcpEnvironment := path.Join(testEnvironmentPath, "test_gcp")
	validvSphereEnvironment := path.Join(testEnvironmentPath, "test_vsphere")
	validOpenstackEnvironment := path.Join(testEnvironmentPath, "test_openstack")
	validAWSEnvironment := path.Join(testEnvironmentPath, "test_aws")

	BeforeEach(func() {
		bash.Source(pathToScript("upload_artefacts"), nil)
		boshMock := MockOrCallThrough("bosh", `echo -n "3124.12"`, `[ "$1" == 'int' ]`)
		ApplyMocks(bash, []Gob{boshMock})
	})

	JustBeforeEach(func() {
		bash.Source("", func(string) ([]byte, error) {
			return repoDirectoryFunction, nil
		})
	})

	Context("When release source is local", func() {
		It("uploads the local release", func() {
			code, err := bash.Run("main", []string{validGcpEnvironment, "local"})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))
			Expect(stderr).To(gbytes.Say("bosh upload-release kubo-release.tgz"))
		})

		It("uploads the stemcell for GCP", func() {
			boshMock := MockOrCallThrough("bosh", `echo -n "3124.12"`, `[ "$1" == 'int' ] && [ ! "$4" == '/stemcells/0/version' ] `)
			ApplyMocks(bash, []Gob{boshMock})

			code, err := bash.Run("main", []string{validGcpEnvironment, "local"})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))

			Expect(stderr).To(gbytes.Say("bosh upload-stemcell https://s3.amazonaws.com/bosh-gce-light-stemcells/light-bosh-stemcell-3124.12-google-kvm-ubuntu-trusty-go_agent.tgz"))
		})

		It("uploads the stemcell for vSphere", func() {
			boshMock := MockOrCallThrough("bosh", `echo -n "3124.12"`, `[ "$1" == 'int' ] && [ ! "$4" == '/stemcells/0/version' ] `)
			ApplyMocks(bash, []Gob{boshMock})
			code, err := bash.Run("main", []string{validvSphereEnvironment, "local"})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))

			Expect(stderr).To(gbytes.Say("bosh upload-stemcell https://s3.amazonaws.com/bosh-core-stemcells/vsphere/bosh-stemcell-3124.12-vsphere-esxi-ubuntu-trusty-go_agent.tgz"))
		})

		It("uploads the stemcell for OpenStack", func() {
			boshMock := MockOrCallThrough("bosh", `echo -n "3124.12"`, `[ "$1" == 'int' ] && [ ! "$4" == '/stemcells/0/version' ] `)
			ApplyMocks(bash, []Gob{boshMock})
			code, err := bash.Run("main", []string{validOpenstackEnvironment, "local"})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))

			Expect(stderr).To(gbytes.Say("bosh upload-stemcell https://s3.amazonaws.com/bosh-core-stemcells/openstack/bosh-stemcell-3124.12-openstack-kvm-ubuntu-trusty-go_agent.tgz"))
		})

		It("uploads the stemcell for AWS", func() {
			boshMock := MockOrCallThrough("bosh", `echo -n "3124.12"`, `[ "$1" == 'int' ] && [ ! "$4" == '/stemcells/0/version' ] `)
			ApplyMocks(bash, []Gob{boshMock})
			code, err := bash.Run("main", []string{validAWSEnvironment, "local"})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))

			Expect(stderr).To(gbytes.Say("bosh upload-stemcell https://s3.amazonaws.com/bosh-core-stemcells/aws/bosh-stemcell-3124.12-aws-xen-hvm-ubuntu-trusty-go_agent.tgz"))
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

			code, err := bash.Run("main", []string{validGcpEnvironment, "dev"})
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

			code, err := bash.Run("main", []string{validGcpEnvironment, "public"})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))
			Expect(stderr).To(gbytes.Say("upload_release url-to-public-release"))
		})
	})
})

var _ = Describe("upload_stemcell", func() {
	BeforeEach(func() {
		bash.Source(pathToScript("upload_artefacts"), nil)
		boshMock := Mock("bosh", "echo")
		ApplyMocks(bash, []Gob{boshMock})
	})

	Context("GCP", func() {
		It("should upload a google kvm ubuntu trusty stemcell", func() {
			code, err := bash.Run("upload_stemcell", []string{"gcp"})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))
			Expect(stderr).To(gbytes.Say("bosh upload-stemcell"))
			Expect(stderr).To(gbytes.Say("google-kvm-ubuntu-trusty-go_agent"))
		})
	})

	Context("vSphere", func() {
		It("should upload a vSphere esxi ubuntu trusty stemcell", func() {
			code, err := bash.Run("upload_stemcell", []string{"vsphere"})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))
			Expect(stderr).To(gbytes.Say("bosh upload-stemcell"))
			Expect(stderr).To(gbytes.Say("vsphere-esxi-ubuntu-trusty-go_agent"))
		})
	})

	Context("OpenStack", func() {
		It("should upload an OpenStack kvm ubuntu trusty stemcell", func() {
			code, err := bash.Run("upload_stemcell", []string{"openstack"})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))
			Expect(stderr).To(gbytes.Say("bosh upload-stemcell"))
			Expect(stderr).To(gbytes.Say("openstack-kvm-ubuntu-trusty-go_agent"))
		})
	})

	Context("AWS", func() {
		It("should upload an AWS xen hvm ubuntu trusty full stemcell", func() {
			code, err := bash.Run("upload_stemcell", []string{"aws"})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))
			Expect(stderr).To(gbytes.Say("bosh-core-stemcells"))
			Expect(stderr).To(gbytes.Say("aws-xen-hvm-ubuntu-trusty-go_agent"))
		})
	})

})
