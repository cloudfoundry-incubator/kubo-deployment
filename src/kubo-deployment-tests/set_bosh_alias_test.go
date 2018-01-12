package kubo_deployment_tests_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	yaml "gopkg.in/yaml.v2"
)

var _ = Describe("Set bosh alias", func() {

	const testEnvName = "test_gcp_with_creds"

	var setBoshAliasPath string
	var kuboEnv string
	var session *gexec.Session
	var cmd *exec.Cmd
	var err error

	BeforeEach(func() {
		setBoshAliasPath = pathToScript("set_bosh_alias")
		kuboEnv = filepath.Join(testEnvironmentPath, testEnvName)
	})

	Context("with incorrect parameters", func() {
		BeforeEach(func() {
			cmd = exec.Command(setBoshAliasPath)
			session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
		})

		It("gives a usage statement", func() {
			Eventually(session).Should(gexec.Exit(1))
			Eventually(session.Out).Should(gbytes.Say("Usage:"))
		})
	})

	Context("with correct parameters", func() {
		BeforeEach(func() {
			cmd = exec.Command("/usr/bin/env", "bash", setBoshAliasPath, kuboEnv)
			cmd.Env = append(os.Environ(), "DEBUG=true")
			session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session, "60s").Should(gexec.Exit())
		})

		It("sets the required environment variables", func() {
			expectedBoshEnv := fmt.Sprintf("export BOSH_ENV=%s", kuboEnv)
			Eventually(session.Err).Should(gbytes.Say(expectedBoshEnv))

			expectedBoshName := fmt.Sprintf("export BOSH_NAME=%s", testEnvName)
			Eventually(session.Err).Should(gbytes.Say(expectedBoshName))
		})

		It("sets the bosh alias", func() {
			var contents []byte
			var directorYaml struct {
				InternalIp string `yaml:"internal_ip"`
				DefaultCA  struct {
					Ca string `yaml:"ca"`
				} `yaml:"default_ca"`
			}

			var credsYaml struct {
				BoshAdminClientSecret string `yaml:"bosh_admin_client_secret"`
			}

			directorYamlPath := filepath.Join(kuboEnv, "director.yml")
			contents, err = ioutil.ReadFile(directorYamlPath)
			Expect(err).NotTo(HaveOccurred())
			err = yaml.Unmarshal(contents, &directorYaml)
			Expect(err).NotTo(HaveOccurred())

			credsYamlPath := filepath.Join(kuboEnv, "creds.yml")
			contents, err = ioutil.ReadFile(credsYamlPath)
			Expect(err).NotTo(HaveOccurred())
			err = yaml.Unmarshal(contents, &credsYaml)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session.Err).Should(gbytes.Say(fmt.Sprintf("bosh alias-env %s -e %s", testEnvName, directorYaml.InternalIp)))
		})
	})
})
