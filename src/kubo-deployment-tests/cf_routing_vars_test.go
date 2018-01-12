package kubo_deployment_tests_test

import (
	"encoding/json"
	"fmt"

	. "github.com/jhvhs/gob-mock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var (
	cfManifest = fmt.Sprintf("cat %s", pathFromRoot("src/kubo-deployment-tests/resources/misc/cf-deployment.yml"))
)

const (
	domainOutput = `echo "Getting domains in org some-org as some-admin...
name                 status   type
some-app-hostname    shared
some-tcp-hostname    shared   tcp"`

	boshInstancesJson = `echo '{
    "Tables": [
        {
            "Content": "instances",
            "Header": {
                "az": "AZ",
                "instance": "Instance",
                "ips": "IPs",
                "process_state": "Process State"
            },
            "Rows": [
                {
                    "az": "z1",
                    "instance": "nats/uuid1",
                    "ips": "10.42.42.42",
                    "process_state": "running"
                },
                {
                    "az": "z2",
                    "instance": "nats/uuid2",
                    "ips": "10.42.42.43",
                    "process_state": "running"
                },
                {
                    "az": "z3",
                    "instance": "other/uuid3",
                    "ips": "10.42.42.44",
                    "process_state": "running"
                }
            ],
            "Notes": null
        }
    ],
    "Blocks": null,
    "Lines": [
        "Using environment '10.42.42.41' as user 'admin' (openid, bosh.admin)",
        "Task 11",
        ". Done",
        "Succeeded"
    ]
	}'`
)

var _ = Describe("CfRoutingVars", func() {
	BeforeEach(func() {
		bash.Source(pathFromRoot("manifests/helper/cf-routing-vars.sh"), nil)
		bash.Export("BOSH_ENVIRONMENT", "bosh-env")
	})

	It("generates routing data from cf deployment", func() {
		boshMock := MockOrCallThrough("bosh", cfManifest, `[[ "$1" == "int" ]]`)
		cfMock := Mock("cf", domainOutput)
		routingClientSecretMock := Mock("get_routing_client_secret", "echo routing-secret")
		natsPasswordMock := Mock("get_nats_password", "echo nats-password")
		natsIPsMock := Mock("get_nats_ips_json", "echo '[10.23.23.10, 10.23.23.11]'")
		ApplyMocks(bash, []Gob{boshMock, cfMock, routingClientSecretMock, natsPasswordMock, natsIPsMock})

		exitCode, err := bash.Run("main", []string{""})
		Expect(err).ToNot(HaveOccurred())
		Expect(exitCode).To(Equal(0))

		Expect(stdout).To(gbytes.Say("kubernetes_master_host: some-tcp-hostname"))
		Expect(stdout).To(gbytes.Say("kubernetes_master_port: 8443"))
		Expect(stdout).To(gbytes.Say("routing_cf_api_url: https://api.system.domain"))
		Expect(stdout).To(gbytes.Say("routing_cf_uaa_url: https://uaa.system.domain"))
		Expect(stdout).To(gbytes.Say("routing_cf_app_domain_name: app.domain"))
		Expect(stdout).To(gbytes.Say("routing_cf_client_id: routing_api_client"))
		Expect(stdout).To(gbytes.Say(`routing_cf_client_secret: "routing-secret"`))
		Expect(stdout).To(gbytes.Say("routing_cf_nats_port: 4222"))
		Expect(stdout).To(gbytes.Say("routing_cf_nats_username: nat"))
		Expect(stdout).To(gbytes.Say(`routing_cf_nats_password: "nats-password"`))
		Expect(stdout).To(gbytes.Say(`routing_cf_nats_internal_ips: \[10.23.23.10, 10.23.23.11\]`))
	})

	Describe("get_routing_client_secret", func() {
		It("gets routing client secret from credhub", func() {
			credhubMock := Mock("credhub", `echo '{"value": "routing-secret"}'`)
			ApplyMocks(bash, []Gob{credhubMock})
			exitCode, err := bash.Run("get_routing_client_secret", []string{})
			Expect(err).ToNot(HaveOccurred())
			Expect(exitCode).To(Equal(0))
			Expect(stdout).To(gbytes.Say("^routing-secret\n$"))
			Expect(stderr).To(gbytes.Say("<1> credhub get -n bosh-env/cf/uaa_clients_routing_api_client_secret --output-json"))
		})
	})

	Describe("get_nats_password", func() {
		It("get nats password from credhub", func() {
			credhubMock := Mock("credhub", `echo '{"value": "nats-password"}'`)
			ApplyMocks(bash, []Gob{credhubMock})
			exitCode, err := bash.Run("get_nats_password", []string{})
			Expect(err).ToNot(HaveOccurred())
			Expect(exitCode).To(Equal(0))
			Expect(stdout).To(gbytes.Say("^nats-password\n$"))
			Expect(stderr).To(gbytes.Say("<1> credhub get -n bosh-env/cf/nats_password --output-json"))
		})
	})

	Describe("get_nats_ips_json", func() {
		It("get nats ips from bosh", func() {
			boshMock := Mock("bosh", boshInstancesJson)
			ApplyMocks(bash, []Gob{boshMock})
			exitCode, err := bash.Run("get_nats_ips_json", []string{})
			Expect(err).ToNot(HaveOccurred())
			Expect(exitCode).To(Equal(0))
			ips := []string{}
			jsonString := stdout.Contents()
			json.Unmarshal(jsonString, &ips)
			Expect(ips).To(Equal([]string{"10.42.42.42", "10.42.42.43"}))
			Expect(stderr).To(gbytes.Say("<1> bosh instances -d cf --json"))
		})
	})
})
