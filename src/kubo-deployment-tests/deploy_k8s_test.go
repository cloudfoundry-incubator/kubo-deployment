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

	Context("skip", func() {
		It("deploys with skip upload successfully", func() {
			code, err := bash.Run("main", []string{validGcpEnvironment, "deployment", "skip"})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))
		})

		It("does not upload the stemcell", func() {
			code, err := bash.Run("main", []string{validGcpEnvironment, "deployment", "skip"})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))

			Expect(stderr).NotTo(gbytes.Say("bosh-cli upload-stemcell"))
		})
	})

	Context("local", func() {
		It("deploys with local upload successfully", func() {
			code, err := bash.Run("main", []string{validGcpEnvironment, "deployment", "local"})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))
		})

		It("uploads the stemcell successfully for GCP", func() {
			boshMock := MockOrCallThrough("bosh-cli", `echo -n "3124.12"`, `[ "$1" == 'int' ] && [ ! "$4" == '/stemcells/0/version' ] `)
			ApplyMocks(bash, []Gob{boshMock})

			code, err := bash.Run("main", []string{validGcpEnvironment, "deployment", "local"})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))

			Expect(stderr).To(gbytes.Say("bosh-cli upload-stemcell https://s3.amazonaws.com/bosh-core-stemcells/google/bosh-stemcell-3124.12-google-kvm-ubuntu-trusty-go_agent.tgz"))
		})

		It("uploads the stemcell successfully for vSphere", func() {
			boshMock := MockOrCallThrough("bosh-cli", `echo -n "3124.12"`, `[ "$1" == 'int' ] && [ ! "$4" == '/stemcells/0/version' ] `)
			ApplyMocks(bash, []Gob{boshMock})
			code, err := bash.Run("main", []string{validvSphereEnvironment, "deployment", "local"})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))

			Expect(stderr).To(gbytes.Say("bosh-cli upload-stemcell https://s3.amazonaws.com/bosh-core-stemcells/vsphere/bosh-stemcell-3124.12-vsphere-esxi-ubuntu-trusty-go_agent.tgz"))
		})

		It("uploads the stemcell successfully for OpenStack", func() {
			boshMock := MockOrCallThrough("bosh-cli", `echo -n "3124.12"`, `[ "$1" == 'int' ] && [ ! "$4" == '/stemcells/0/version' ] `)
			ApplyMocks(bash, []Gob{boshMock})
			code, err := bash.Run("main", []string{validOpenstackEnvironment, "deployment", "local"})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))

			Expect(stderr).To(gbytes.Say("bosh-cli upload-stemcell https://s3.amazonaws.com/bosh-core-stemcells/openstack/bosh-stemcell-3124.12-openstack-kvm-ubuntu-trusty-go_agent.tgz"))
		})

		It("uploads the stemcell successfully for AWS", func() {
			boshMock := MockOrCallThrough("bosh-cli", `echo -n "3124.12"`, `[ "$1" == 'int' ] && [ ! "$4" == '/stemcells/0/version' ] `)
			ApplyMocks(bash, []Gob{boshMock})
			code, err := bash.Run("main", []string{validAWSEnvironment, "deployment", "local"})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))

			Expect(stderr).To(gbytes.Say("bosh-cli upload-stemcell https://s3.amazonaws.com/bosh-aws-light-stemcells/light-bosh-stemcell-3124.12-aws-xen-hvm-ubuntu-trusty-go_agent.tgz"))
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
}'`, `[ "$1" != 'environment' ]`)
		ApplyMocks(bash, []Gob{boshMock})

		code, err := bash.Run("get_director_uuid", []string{})

		Expect(err).NotTo(HaveOccurred())
		Expect(code).To(Equal(0))
		Expect(stdout).To(gbytes.Say("director-uuid"))
	})
})
