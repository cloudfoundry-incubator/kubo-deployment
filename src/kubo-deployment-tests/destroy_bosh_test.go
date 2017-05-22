package kubo_deployment_tests_test

import (
	"path"

	. "github.com/jhvhs/gob-mock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Destroy KuBOSH", func() {
	validGcpEnvironment := path.Join(testEnvironmentPath, "test_gcp_with_creds")
	validvSphereEnvironment := path.Join(testEnvironmentPath, "test_vsphere")
	validOpenstackEnvironment := path.Join(testEnvironmentPath, "test_openstack")

	Context("fails", func() {
		JustBeforeEach(func() {
			ApplyMocks(bash, []Gob{Stub("bosh-cli")})
			bash.Source("", func(string) ([]byte, error) {
				return repoDirectoryFunction, nil
			})
		})

		DescribeTable("when wrong number of arguments is used", func(params []string) {
			code, err := bash.Run(pathToScript("deploy_bosh"), params)
			Expect(err).NotTo(HaveOccurred())
			Expect(code).NotTo(Equal(0))

		},
			Entry("has no arguments", []string{}),
			Entry("has one argument", []string{"gcp"}),
			Entry("has three arguments", []string{"gcp", "foo", "bar"}),
		)

		It("is given an invalid environment path", func() {
			code, err := bash.Run(pathToScript("destroy_bosh"), []string{pathFromRoot(""), pathFromRoot("README.md")})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).NotTo(Equal(0))
		})

		It("is given a non-existing file", func() {
			code, err := bash.Run(pathToScript("destroy_bosh"), []string{validGcpEnvironment, pathFromRoot("non-existing.file")})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).NotTo(Equal(0))
		})
	})

	Context("succeeds", func() {
		BeforeEach(func() {
			bash.Source(pathToScript("destroy_bosh"), nil)
			bashMock := MockOrCallThrough("bosh-cli", `echo "bosh-cli $@" >&2`, `[[ "$1" == "int" ]]`)
			ApplyMocks(bash, []Gob{bashMock})
		})

		It("runs with a valid environment and an extra file", func() {
			code, err := bash.Run("main", []string{validGcpEnvironment, pathFromRoot("README.md")})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))
		})

		It("cleans up before deleting the GCP environment", func() {
			bash.Run("main", []string{validGcpEnvironment, pathFromRoot("README.md")})
			Expect(stderr).To(gbytes.Say("bosh-cli clean-up --all"))
			Expect(stderr).To(gbytes.Say("bosh-cli delete-env"))
		})

		It("cleans up before deleting the OpenStack environment", func() {
			bash.Run("main", []string{validOpenstackEnvironment, pathFromRoot("README.md")})
			Expect(stderr).To(gbytes.Say("bosh-cli clean-up --all"))
			Expect(stderr).To(gbytes.Say("bosh-cli delete-env"))
		})

		It("destroys a GCP environment", func() {
			code, err := bash.Run("main", []string{validGcpEnvironment, pathFromRoot("README.md")})
			Expect(stderr).To(gbytes.Say("bosh-cli delete-env"))
			Expect(stderr).To(gbytes.Say("/bosh-deployment/gcp/cpi.yml"))
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))
		})

		It("destroys a vSphere environment", func() {
			code, err := bash.Run("main", []string{validvSphereEnvironment, pathFromRoot("README.md")})
			Expect(stderr).To(gbytes.Say("bosh-cli delete-env"))
			Expect(stderr).To(gbytes.Say("/bosh-deployment/vsphere/cpi.yml"))
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))
		})

		It("destroys an Openstack environment", func() {
			code, err := bash.Run("main", []string{validOpenstackEnvironment, pathFromRoot("README.md")})
			Expect(stderr).To(gbytes.Say("bosh-cli delete-env"))
			Expect(stderr).To(gbytes.Say("/bosh-deployment/openstack/cpi.yml"))
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))
		})

		It("expands the environment path", func() {
			relativePath := testEnvironmentPath + "/../environments/test_gcp_with_creds"
			code, err := bash.Run("main", []string{relativePath, pathFromRoot("README.md")})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(0))
			Expect(stderr).NotTo(gbytes.Say("\\.\\./environments/test_gcp"))
		})
	})
})
