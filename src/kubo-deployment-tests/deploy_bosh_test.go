package kubo_deployment_tests_test

import (
	"path"

	"fmt"
	. "github.com/jhvhs/gob-mock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Deploy KuBOSH", func() {
	validGcpEnvironment := path.Join(testEnvironmentPath, "test_gcp_with_creds")
	validvSphereEnvironment := path.Join(testEnvironmentPath, "test_vsphere_with_creds")
	validOpenstackEnvironment := path.Join(testEnvironmentPath, "test_openstack_with_creds")
	validAwsEnvironment := path.Join(testEnvironmentPath, "test_aws_with_creds")

	JustBeforeEach(func() {
		bash.Source("", func(string) ([]byte, error) {
			return repoDirectoryFunction, nil
		})
	})

	mockKeyFile := pathFromRoot("README.md")
	Context("fails", func() {
		BeforeEach(func() {
			boshCli := SpyAndConditionallyCallThrough("bosh-cli", "[[ \"$1\" =~ ^int ]]")
			ApplyMocks(bash, []Gob{boshCli})
		})

		DescribeTable("when wrong number of arguments is used", func(params []string) {
			script := pathToScript("deploy_bosh")
			bash.Source(script, nil)
			code, err := bash.Run("main", params)
			Expect(err).NotTo(HaveOccurred())
			Expect(code).NotTo(Equal(0))
			Expect(stderr).To(gbytes.Say("Usage: "))
		},
			Entry("has no arguments", []string{}),
			Entry("has one argument", []string{"gcp"}),
			Entry("has three arguments", []string{"gcp", "foo", "bar"}),
		)

		It("is given an invalid environment path", func() {
			code, err := bash.Run(pathToScript("deploy_bosh"), []string{pathFromRoot(""), mockKeyFile})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).NotTo(Equal(0))
		})

		It("is given a non-existing file", func() {
			code, err := bash.Run(pathToScript("deploy_bosh"), []string{validGcpEnvironment, pathFromRoot("non-existing.file")})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).NotTo(Equal(0))
		})
	})

	Context("succeeds", func() {
		BeforeEach(func() {
			bash.Source(pathToScript("deploy_bosh"), nil)
			boshMock := MockOrCallThrough("bosh-cli", `echo "bosh-cli $@" >&2`, "[ $1 == 'int' ]")
			ApplyMocks(bash, []Gob{boshMock})
		})

		It("deploys to a GCP environment", func() {
			code, err := bash.Run("main", []string{validGcpEnvironment, mockKeyFile})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))
			Expect(stderr).To(gbytes.Say("/bosh-deployment/gcp/cpi.yml"))
		})

		It("deploys to a vSphere environment", func() {
			code, err := bash.Run("main", []string{validvSphereEnvironment, mockKeyFile})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))
			Expect(stderr).To(gbytes.Say("/bosh-deployment/vsphere/cpi.yml"))
		})

		It("deploys to an Openstack environment", func() {
			code, err := bash.Run("main", []string{validOpenstackEnvironment, mockKeyFile})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))
			Expect(stderr).To(gbytes.Say("/bosh-deployment/openstack/cpi.yml"))
		})

		It("deploys to an AWS environment", func() {
			code, err := bash.Run("main", []string{validAwsEnvironment, mockKeyFile})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))
			Expect(stderr).To(gbytes.Say("/bosh-deployment/aws/cpi.yml"))
		})

		It("expands the environment path", func() {
			relativePath := testEnvironmentPath + "/../environments/test_gcp_with_creds"
			code, err := bash.Run("main", []string{relativePath, mockKeyFile})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))
			Expect(stderr).NotTo(gbytes.Say("\\.\\./environments/test_gcp"))
		})

		Context("To enable BOSH DNS", func() {
			It("enables local DNS", func() {
				code, err := bash.Run("generate_manifest_generic", []string{validGcpEnvironment})
				Expect(err).NotTo(HaveOccurred())
				Expect(code).To(Equal(0))
				Expect(stderr).To(gbytes.Say("/bosh-deployment/local-dns.yml"))
				Expect(stdout).To(gbytes.Say("local_dns:\n[ ]+enabled: true"))
			})

			It("adds BOSH DNS runtime config", func() {
				code, err := bash.Run("main", []string{validGcpEnvironment, mockKeyFile})
				Expect(err).NotTo(HaveOccurred())
				Expect(code).To(Equal(0))
				Expect(stderr).To(gbytes.Say(fmt.Sprintf("update-runtime-config -n %s", pathFromRoot("bosh-deployment/runtime-configs/dns.yml"))))
			})

		})
	})
})
