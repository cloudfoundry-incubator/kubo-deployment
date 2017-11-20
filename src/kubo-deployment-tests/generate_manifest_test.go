package kubo_deployment_tests_test

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"io"

	boshtpl "github.com/cloudfoundry/bosh-cli/director/template"
	"github.com/cppforlife/go-patch/patch"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"strings"
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

		DescribeTable("populated properties for CF-based deployment", func(yPath, value string) {
			cfEnv := filepath.Join(testEnvironmentPath, "test_vsphere_with_creds")
			status, _ := bash.Run("main", []string{cfEnv, "klingon", "director_uuid"})

			Expect(status).To(Equal(0))

			pathValue, err := propertyFromManifest(yPath, stdout.Contents())
			Expect(err).NotTo(HaveOccurred())
			Expect(pathValue).To(Equal(value))
		},
			Entry("deployment name", "/name", "klingon"),
			Entry("network name", "/instance_groups/name=master/networks/0/name", "network-name"),
			Entry("Master node has the etcd job", "/instance_groups/name=master/jobs/name=etcd/release", "kubo-etcd"),
			Entry("There is only one master node", "/instance_groups/name=master/instances", "1"),
			Entry("Master Kubeconfig Kubelet Password", "/instance_groups/name=master/jobs/name=kubeconfig/properties/kubelet-password", "((kubelet-password))"),
			Entry("Worker Kubeconfig Kubernetes API URL", "/instance_groups/name=worker/jobs/name=kubeconfig/properties/kubernetes-api-url", "https://a.router.name:101928"),
			Entry("kubernetes external port", "/instance_groups/name=master/jobs/name=kubernetes-api-route-registrar/properties/external_kubo_port", "101928"),
			Entry("CF API URL", "/instance_groups/name=master/jobs/name=kubernetes-api-route-registrar/properties/cloud_foundry/api_url", "cf.api.url"),
			Entry("CF UAA URL", "/instance_groups/name=master/jobs/name=kubernetes-api-route-registrar/properties/cloud_foundry/uaa_url", "cf.uaa.url"),
			Entry("CF Client ID", "/instance_groups/name=master/jobs/name=kubernetes-api-route-registrar/properties/cloud_foundry/uaa_client_id", "cf.client.id"),
			Entry("CF Client Secret", "/instance_groups/name=master/jobs/name=kubernetes-api-route-registrar/properties/cloud_foundry/uaa_client_secret", "cf.client.secret"),
			Entry("Auto-generated kubelet password", "/instance_groups/name=master/jobs/name=kubernetes-api/properties/kubelet-password", "((kubelet-password))"),
			Entry("Auto-generated admin password", "/instance_groups/name=master/jobs/name=kubernetes-api/properties/admin-password", "((kubo-admin-password))"),
		)

		DescribeTable("populated properties for IaaS-based deployment", func(yPath, value string) {
			status, _ := bash.Run("main", []string{kuboEnv, "grinder", "director_uuid"})

			Expect(status).To(Equal(0))

			pathValue, err := propertyFromManifest(yPath, stdout.Contents())
			Expect(err).NotTo(HaveOccurred())
			Expect(pathValue).To(Equal(value))
		},
			Entry("deployment name", "/name", "grinder"),
			Entry("network name", "/instance_groups/name=master/networks/0/name", "network-name"),
			Entry("Auto-generated kubelet password", "/instance_groups/name=master/jobs/name=kubernetes-api/properties/kubelet-password", "((kubelet-password))"),
			Entry("Auto-generated admin password", "/instance_groups/name=master/jobs/name=kubernetes-api/properties/admin-password", "((kubo-admin-password))"),
			Entry("worker node tag", "/instance_groups/name=master/jobs/name=cloud-provider/properties/cloud-provider/gce/worker-node-tag", "TheDirector-grinder-worker"),
			Entry("Kubernetes API URL", "/instance_groups/name=worker/jobs/name=kubeconfig/properties/kubernetes-api-url", "https://12.23.34.45:101928"),
		)

		It("should always use dns addresses", func() {
			status, err := bash.Run("main", []string{kuboEnv, "grinder", "director_uuid"})

			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(0))

			pathValue, err := propertyFromManifest("/features/use_dns_addresses", stdout.Contents())
			Expect(err).NotTo(HaveOccurred())
			Expect(pathValue).To(Equal("true"))
		})

		It("should include a variable section with tls-kubelet, tls-kubernetes", func() {
			status, err := bash.Run("main", []string{kuboEnv, "cucumber", "director_uuid"})
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(0))

			Expect(stdout).To(gbytes.Say("variables:"))
			Expect(stdout).To(gbytes.Say("tls-kubelet"))
			Expect(stdout).To(gbytes.Say("tls-kubernetes"))
		})

		It("should include an alternative name with master.kubo for the tls-kubernetes variable", func() {
			status, err := bash.Run("main", []string{kuboEnv, "cucumber", "director_uuid"})
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(0))

			Expect(stdout).To(gbytes.Say("variables:"))
			Expect(stdout).To(gbytes.Say("tls-kubernetes"))
			Expect(stdout).To(gbytes.Say("alternative_names:"))
			Expect(stdout).To(gbytes.Say("master.kubo"))
		})

		It("should default the authorization mode property to ABAC", func() {
			status, err := bash.Run("main", []string{kuboEnv, "cucumber", "director_uuid"})
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(0))

			Expect(stdout).To(gbytes.Say("authorization-mode: abac"))
		})

		It("should use the abac authorization mode set in the kubo environment", func() {
			abacEnv := filepath.Join(testEnvironmentPath, "test_gcp_abac")
			status, err := bash.Run("main", []string{abacEnv, "cucumber", "director_uuid"})
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(0))

			Expect(stdout).To(gbytes.Say("authorization-mode: abac"))
		})

		It("should use the rbac authorization mode set in the kubo environment", func() {
			rbacEnv := filepath.Join(testEnvironmentPath, "test_gcp_rbac")
			status, err := bash.Run("main", []string{rbacEnv, "cucumber", "director_uuid"})
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(0))

			Expect(stdout).To(gbytes.Say("authorization-mode: rbac"))
		})

		It("should reproduce the same manifest on the second run", func() {
			bash.Run("main", []string{kuboEnv, "fort", "director_uuid"})

			firstRun := make([]byte, len(stdout.Contents()))
			secondRun := make([]byte, len(stdout.Contents()))

			_, err := stdout.Read(firstRun)
			Expect(err).NotTo(HaveOccurred())

			bash.Run("main", []string{kuboEnv, "fort", "director_uuid"})

			_, err = stdout.Read(secondRun)
			Expect(err).NotTo(HaveOccurred())

			Expect(firstRun).To(Equal(secondRun))
		})

		It("does not include aws tags in the gcp manifest", func() {
			bash.Run("main", []string{kuboEnv, "foo", "bar"})

			Expect(stdout).NotTo(gbytes.Say("\ntags:\n  KubernetesCluster:"))
		})

		It("does include aws tags in the aws with iaas routing manifest", func() {
			bash.Run("main", []string{filepath.Join(testEnvironmentPath, "test_aws"), "zing", "x"})

			Expect(stdout).To(gbytes.Say("\ntags:\n  KubernetesCluster:"))
		})

		It("does include aws tags in the aws with cf routing mode manifest", func() {
			bash.Run("main", []string{filepath.Join(testEnvironmentPath, "test_aws_cf"), "zing", "x"})

			Expect(stdout).To(gbytes.Say("\ntags:\n  KubernetesCluster:"))
		})

		It("generates a manifest without the secrets", func() {
			secretlessEnv := filepath.Join(testEnvironmentPath, "secretless")
			status, err := bash.Run("main", []string{secretlessEnv, "sensors", "director_uuid"})
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(0))

			Expect(stdout).To(gbytes.Say("\n        uaa_client_secret: \\(\\(routing-cf-client-secret\\)\\)\n"))
		})

		It("uses ops-files to modify the manifest", func() {
			opsfileEnv := filepath.Join(testEnvironmentPath, "with_ops")
			status, err := bash.Run("main", []string{opsfileEnv, "name", "director_uuid"})

			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(0))
			Expect(stdout).To(gbytes.Say("\n  os: MALARIA0\n"))
		})

		It("applies http proxy settings if they exist", func() {
			opsfileEnv := filepath.Join(testEnvironmentPath, "with_http_proxy")
			status, err := bash.Run("main", []string{opsfileEnv, "name", "director_uuid"})

			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(0))
			Expect(stdout).To(gbytes.Say("\n        http_proxy: my.proxy.com\n"))
		})

		It("applies https proxy settings if they exist", func() {
			opsfileEnv := filepath.Join(testEnvironmentPath, "with_https_proxy")
			status, err := bash.Run("main", []string{opsfileEnv, "name", "director_uuid"})

			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(0))
			Expect(stdout).To(gbytes.Say("\n        https_proxy: my.sslproxy.com\n"))
		})

		It("applies http no_proxy settings if they exist", func() {
			opsfileEnv := filepath.Join(testEnvironmentPath, "with_no_proxy")
			status, err := bash.Run("main", []string{opsfileEnv, "name", "director_uuid"})

			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(0))
			Expect(stdout).To(gbytes.Say("\n        no_proxy: dont.proxy.me\n"))
		})

		It("uses vars-files to modify the manifest", func() {
			opsfileEnv := filepath.Join(testEnvironmentPath, "with_vars")
			status, err := bash.Run("main", []string{opsfileEnv, "name", "director_uuid"})

			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(0))
			Expect(stdout).To(gbytes.Say("\n      kubelet-password: Shields up, ancient life!\n"))
		})

		It("should not embed addons_specs if not specified in the director", func() {
			status, err := bash.Run("main", []string{kuboEnv, "grinder", "director_uuid"})

			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(0))

			pathValue, err := propertyFromManifest("/instance_groups/name=master/jobs/name=kubernetes-system-specs/properties/addons_spec", stdout.Contents())
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Expected to find a map key 'addons_spec'"))
			Expect(pathValue).To(Equal(""))
		})

		It("should embed addons_specs", func() {
			opsfileEnv := filepath.Join(testEnvironmentPath, "with_addons")
			status, err := bash.Run("main", []string{opsfileEnv, "name", "director_uuid"})

			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(0))
			pathValue, err := propertyFromManifest("/instance_groups/name=master/jobs/name=kubernetes-system-specs/properties/addons-spec", stdout.Contents())

			Expect(pathValue).To(Equal("valid:\n  key: value"))
		})
	})

	It("errors out if addons_spec file is missing", func() {
		noAddonsEnv := filepath.Join(testEnvironmentPath, "absent_addons_spec_failing")
		status, err := bash.Run("main", []string{noAddonsEnv, "coaster", "director_uuid"})

		Expect(err).NotTo(HaveOccurred())

		Expect(status).NotTo(Equal(0))
		Expect(strings.Contains(string(stdout.Contents()), "No file exists")).To(BeTrue())
	})

	It("errors out if addons_spec file is not valid yaml", func() {
		invalidAddonsEnv := filepath.Join(testEnvironmentPath, "invalid_addons_spec_failing")
		status, err := bash.Run("main", []string{invalidAddonsEnv, "coaster", "director_uuid"})

		Expect(err).NotTo(HaveOccurred())
		Expect(status).NotTo(Equal(0))
		Expect(strings.Contains(string(stdout.Contents()), "Invalid yaml")).To(BeTrue())
	})

	It("expands the bosh environment path to absolute value", func() {
		command := exec.Command("./generate_kubo_manifest", "../src/kubo-deployment-tests/resources/environments/test_gcp", "name", "director_uuid")
		command.Stdout = bash.Stdout
		command.Stderr = bash.Stderr
		command.Dir = pathToScript("")
		Expect(command.Run()).To(Succeed())
	})

	It("runs from any location", func() {
		command := exec.Command("./bin/generate_kubo_manifest", "src/kubo-deployment-tests/resources/environments/test_gcp", "name", "director_uuid")
		command.Stdout = bash.Stdout
		command.Stderr = bash.Stderr
		command.Dir = pathFromRoot("")
		Expect(command.Run()).To(Succeed())
	})

	It("cloud provider properties are populated for GCP", func() {
		command := exec.Command("./bin/generate_kubo_manifest", "src/kubo-deployment-tests/resources/environments/test_gcp", "name", "director_uuid")
		command.Stdout = bash.Stdout
		command.Stderr = bash.Stderr
		command.Dir = pathFromRoot("")
		Expect(command.Run()).To(Succeed())
		Expect(stdout).To(gbytes.Say("project-id: GCP_Project_ID"))
	})

	It("cloud provider properties are populated for vSphere", func() {
		command := exec.Command("./bin/generate_kubo_manifest", "src/kubo-deployment-tests/resources/environments/test_vsphere", "name", "director_uuid")
		command.Stdout = bash.Stdout
		command.Stderr = bash.Stderr
		command.Dir = pathFromRoot("")
		Expect(command.Run()).To(Succeed())
		Expect(stdout).To(gbytes.Say("working-dir: /big_data_center/vm/big_vms/director_uuid"))
	})

	It("should generate a valid manifest", func() {
		files, _ := filepath.Glob(testEnvironmentPath + "/*")
		for _, env := range files {
			if strings.Contains(env, "_failing"){
				continue
			}
			command := exec.Command("./bin/generate_kubo_manifest", env, "env-name", "director_uuid")
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
			if strings.Contains(env, "_failing"){
				continue
			}
			command := exec.Command("./bin/generate_kubo_manifest", env, "env-name", "director_uuid")
			errBuffer := gbytes.NewBuffer()
			command.Stdout = GinkgoWriter
			command.Stderr = errBuffer
			command.Dir = pathFromRoot("")
			Expect(command.Run()).To(Succeed(), fmt.Sprintf("Failed with environmenrt %s", env))
			Expect(string(errBuffer.Contents())).To(HaveLen(0))
		}
	})

	It("should set the tls-kubernetes common_name to the kubernetes_master_host", func() {
		command := exec.Command("./bin/generate_kubo_manifest", "src/kubo-deployment-tests/resources/environments/test_external", "name", "director_uuid")

		stdoutTemp := gbytes.NewBuffer()
		stderrTemp := gbytes.NewBuffer()

		command.Stdout = io.MultiWriter(stdoutTemp, GinkgoWriter)
		command.Stderr = io.MultiWriter(stderrTemp, GinkgoWriter)
		command.Dir = pathFromRoot("")
		Expect(command.Run()).To(Succeed())

		command2 := exec.Command("bosh-cli", "int", "-", "--path", "/variables/name=tls-kubernetes/options/common_name")
		command2.Stdin = stdoutTemp
		command2.Stdout = bash.Stdout
		command2.Stderr = bash.Stderr

		Expect(command2.Run()).To(Succeed())
		Expect(stdout).To(gbytes.Say("12.23.34.45"))
	})

	It("should add the kubernetes_master_host to tls-kubernetes alternative_names", func() {
		command := exec.Command("./bin/generate_kubo_manifest", "src/kubo-deployment-tests/resources/environments/test_external", "name", "director_uuid")

		stdoutTemp := gbytes.NewBuffer()
		stderrTemp := gbytes.NewBuffer()

		command.Stdout = io.MultiWriter(stdoutTemp, GinkgoWriter)
		command.Stderr = io.MultiWriter(stderrTemp, GinkgoWriter)
		command.Dir = pathFromRoot("")
		Expect(command.Run()).To(Succeed())

		command2 := exec.Command("bosh-cli", "int", "-", "--path", "/variables/name=tls-kubernetes/options/alternative_names")
		command2.Stdin = stdoutTemp
		command2.Stdout = bash.Stdout
		command2.Stderr = bash.Stderr

		Expect(command2.Run()).To(Succeed())
		Expect(stdout).To(gbytes.Say("- 12.23.34.45"))
	})

	It("should assign static ip to the master in vsphere with proxy mode", func() {
		command := exec.Command("./bin/generate_kubo_manifest", "src/kubo-deployment-tests/resources/environments/test_vsphere_with_haproxy", "name", "director_uuid")
		command.Stdout = bash.Stdout
		command.Stderr = bash.Stderr
		command.Dir = pathFromRoot("")
		Expect(command.Run()).To(Succeed())

		value, err := propertyFromManifest("/instance_groups/name=master/networks/0/static_ips", stdout.Contents())
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal("- 1.2.3.4"))
	})

	It("should assign static ip to the master in open stack with proxy mode", func() {
		command := exec.Command("./bin/generate_kubo_manifest", "src/kubo-deployment-tests/resources/environments/test_openstack_with_haproxy", "name", "director_uuid")
		command.Stdout = bash.Stdout
		command.Stderr = bash.Stderr
		command.Dir = pathFromRoot("")
		Expect(command.Run()).To(Succeed())

		value, err := propertyFromManifest("/instance_groups/name=master/networks/1/type", stdout.Contents())
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal("vip"))

		value, err = propertyFromManifest("/instance_groups/name=master/networks/1/static_ips", stdout.Contents())
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal("- 1.2.3.4"))
	})

})

func propertyFromManifest(path string, manifestString []byte) (string, error) {
	var manifest map[string]interface{}
	err := yaml.Unmarshal(manifestString, &manifest)
	if err != nil {
		return "", err
	}

	template := boshtpl.NewTemplate([]byte(stdout.Contents()))
	vars := boshtpl.StaticVariables{}
	ops := patch.FindOp{Path: patch.MustNewPointerFromString(path)}

	bytes, err := template.Evaluate(vars, ops, boshtpl.EvaluateOpts{ExpectAllKeys: false})
	return choppedString(bytes), err
}

func choppedString(bytes []byte) string {
	if len(bytes) > 0 {
		return string(bytes[:len(bytes)-1])
	} else {
		return ""
	}
}
