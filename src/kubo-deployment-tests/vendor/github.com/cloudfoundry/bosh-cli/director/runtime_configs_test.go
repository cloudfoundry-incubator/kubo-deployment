package director_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	. "github.com/cloudfoundry/bosh-cli/director"
)

var _ = Describe("Director", func() {
	var (
		director Director
		server   *ghttp.Server
	)

	BeforeEach(func() {
		director, server = BuildServer()
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("LatestRuntimeConfig", func() {
		It("returns latest default runtime config if there is at least one", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/runtime_configs", "name=&limit=1"),
					ghttp.VerifyBasicAuth("username", "password"),
					ghttp.RespondWith(http.StatusOK, `[
	{"properties": "first"},
	{"properties": "second"}
]`),
				),
			)

			cc, err := director.LatestRuntimeConfig("")
			Expect(err).ToNot(HaveOccurred())
			Expect(cc).To(Equal(RuntimeConfig{Properties: "first"}))
		})

		It("returns named runtime config if there is at least one and name is specified", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/runtime_configs", "name=foo-name&limit=1"),
					ghttp.VerifyBasicAuth("username", "password"),
					ghttp.RespondWith(http.StatusOK, `[
	{"properties": "first"},
	{"properties": "second"}
]`),
				),
			)

			cc, err := director.LatestRuntimeConfig("foo-name")
			Expect(err).ToNot(HaveOccurred())
			Expect(cc).To(Equal(RuntimeConfig{Properties: "first"}))
		})

		It("returns error if there is no runtime config", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/runtime_configs", "name=&limit=1"),
					ghttp.VerifyBasicAuth("username", "password"),
					ghttp.RespondWith(http.StatusOK, `[]`),
				),
			)

			_, err := director.LatestRuntimeConfig("")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("No runtime config"))
		})

		It("returns error if info response in non-200", func() {
			AppendBadRequest(ghttp.VerifyRequest("GET", "/runtime_configs"), server)

			_, err := director.LatestRuntimeConfig("")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(
				"Finding runtime configs: Director responded with non-successful status code"))
		})

		It("returns error if info cannot be unmarshalled", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/runtime_configs"),
					ghttp.RespondWith(http.StatusOK, ``),
				),
			)

			_, err := director.LatestRuntimeConfig("")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(
				"Finding runtime configs: Unmarshaling Director response"))
		})
	})

	Describe("UpdateRuntimeConfig", func() {
		It("updates runtime config", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/runtime_configs", "name=smurf-runtime-config"),
					ghttp.VerifyBasicAuth("username", "password"),
					ghttp.VerifyHeader(http.Header{
						"Content-Type": []string{"text/yaml"},
					}),
					ghttp.RespondWith(http.StatusOK, `{}`),
				),
			)

			err := director.UpdateRuntimeConfig("smurf-runtime-config", []byte("config"))
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns error if info response in non-200", func() {
			AppendBadRequest(ghttp.VerifyRequest("POST", "/runtime_configs", "name="), server)

			err := director.UpdateRuntimeConfig("", nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(
				"Updating runtime config: Director responded with non-successful status code"))
		})
	})

	Describe("DiffRuntimeConfig", func() {
		var expectedDiffResponse ConfigDiff

		expectedDiffResponse = ConfigDiff{
			Diff: [][]interface{}{
				[]interface{}{"release:", nil},
				[]interface{}{"  version: 0.0.1", "removed"},
				[]interface{}{"  version: 0.0.2", "added"},
			},
		}

		It("diffs the runtime config with the given name", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/runtime_configs/diff", "name=rc1"),
					ghttp.VerifyBasicAuth("username", "password"),
					ghttp.VerifyHeader(http.Header{
						"Content-Type": []string{"text/yaml"},
					}),
					ghttp.RespondWith(http.StatusOK, `{"diff":[["release:",null],["  version: 0.0.1","removed"],["  version: 0.0.2","added"]]}`),
				),
			)

			diff, err := director.DiffRuntimeConfig("rc1", []byte("config"), false)
			Expect(err).ToNot(HaveOccurred())
			Expect(diff).To(Equal(expectedDiffResponse))
		})

		It("returns error if info response in non-200", func() {
			AppendBadRequest(ghttp.VerifyRequest("POST", "/runtime_configs/diff", "name=smurf"), server)

			_, err := director.DiffRuntimeConfig("smurf", nil, false)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(
				"Fetching diff result: Director responded with non-successful status code"))
		})

		It("is backwards compatible with directors without the `/diff` endpoint", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/runtime_configs/diff", "name=rc1"),
					ghttp.VerifyBasicAuth("username", "password"),
					ghttp.VerifyHeader(http.Header{
						"Content-Type": []string{"text/yaml"},
					}),
					ghttp.RespondWith(http.StatusNotFound, ""),
				),
			)

			diff, err := director.DiffRuntimeConfig("rc1", []byte("config"), false)
			Expect(err).ToNot(HaveOccurred())
			Expect(diff).To(Equal(ConfigDiff{}))
		})

		Context("when 'noRedact' is true", func() {
			It("does pass redact parameter to director", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", "/runtime_configs/diff", "name=rc1&redact=false"),
						ghttp.VerifyBasicAuth("username", "password"),
						ghttp.VerifyHeader(http.Header{
							"Content-Type": []string{"text/yaml"},
						}),
						ghttp.RespondWith(http.StatusOK, `{"diff":[["fake-release:",null]]}`),
					),
				)

				_, err := director.DiffRuntimeConfig("rc1", []byte("config"), true)
				Expect(err).ToNot(HaveOccurred())
			})
		})

	})
})
