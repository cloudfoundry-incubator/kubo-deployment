package manifest_test

import (
	. "github.com/cloudfoundry/bosh-cli/deployment/manifest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	biproperty "github.com/cloudfoundry/bosh-utils/property"
)

var _ = Describe("Network", func() {
	var network Network

	Describe("Interface", func() {
		Describe("network defaults", func() {
			BeforeEach(func() {
				network = Network{
					Name:            "fake-manual-network-name",
					Type:            "whatever",
					CloudProperties: biproperty.Map{"cp_key": "cp_value"},
				}
			})

			It("includes default information when defined", func() {
				iface, err := network.Interface([]string{}, []NetworkDefault{"foo", "bar"})
				Expect(err).ToNot(HaveOccurred())
				Expect(iface).To(Equal(biproperty.Map{
					"cloud_properties": biproperty.Map{"cp_key": "cp_value"},
					"default":          []NetworkDefault{"foo", "bar"},
					"type":             "whatever",
				}))
			})
		})

		Context("when network type is manual", func() {
			Context("when network is ipv4", func() {
				BeforeEach(func() {
					network = Network{
						Name: "fake-manual-network-name",
						Type: "manual",
						Subnets: []Subnet{
							{
								Range:   "1.2.3.0/22",
								Gateway: "1.1.1.1",
								DNS:     []string{"1.2.3.4"},
								CloudProperties: biproperty.Map{
									"cp_key": "cp_value",
								},
							},
						},
					}
				})

				It("includes gateway, dns, ip from the job and netmask calculated from range", func() {
					iface, err := network.Interface([]string{"5.6.7.9"}, []NetworkDefault{})
					Expect(err).ToNot(HaveOccurred())
					Expect(iface).To(Equal(biproperty.Map{
						"type":    "manual",
						"ip":      "5.6.7.9",
						"gateway": "1.1.1.1",
						"netmask": "255.255.252.0",
						"dns":     []string{"1.2.3.4"},
						"cloud_properties": biproperty.Map{
							"cp_key": "cp_value",
						},
					}))
				})

				Context("when range is invalid", func() {
					BeforeEach(func() {
						network.Subnets[0].Range = "invalid-range"
					})

					It("returns an error", func() {
						_, err := network.Interface([]string{"5.6.7.9"}, []NetworkDefault{})
						Expect(err).To(HaveOccurred())
						Expect(err.Error()).To(ContainSubstring("Failed to parse subnet range"))
					})
				})
			})

			Context("when network is ipv6", func() {
				BeforeEach(func() {
					network = Network{
						Name: "fake-manual-network-name",
						Type: "manual",
						Subnets: []Subnet{
							{
								Range:   "fd7a:eeed:e696:969f:0000:0000:0000:0000/64",
								Gateway: "fd7a:eeed:e696:969f:0000:0000:0000:0001",
								DNS:     []string{"fd7a:eeed:e696:969f:0000:0000:0000:0008"},
								CloudProperties: biproperty.Map{
									"cp_key": "cp_value",
								},
							},
						},
					}
				})

				It("includes gateway, dns, ip from the job and netmask calculated from range", func() {
					iface, err := network.Interface([]string{"fd7a:eeed:e696:969f:0000:0000:0000:0004"}, []NetworkDefault{})
					Expect(err).ToNot(HaveOccurred())
					Expect(iface).To(Equal(biproperty.Map{
						"type":    "manual",
						"ip":      "fd7a:eeed:e696:969f:0000:0000:0000:0004",
						"gateway": "fd7a:eeed:e696:969f:0000:0000:0000:0001",
						"netmask": "ffff:ffff:ffff:ffff:0000:0000:0000:0000",
						"dns":     []string{"fd7a:eeed:e696:969f:0000:0000:0000:0008"},
						"cloud_properties": biproperty.Map{
							"cp_key": "cp_value",
						},
					}))
				})

				Context("when range is invalid", func() {
					BeforeEach(func() {
						network.Subnets[0].Range = "invalid-range"
					})

					It("returns an error", func() {
						_, err := network.Interface([]string{"fd7a:eeed:e696:969f:0000:0000:0000:0004"}, []NetworkDefault{})
						Expect(err).To(HaveOccurred())
						Expect(err.Error()).To(ContainSubstring("Failed to parse subnet range"))
					})
				})
			})
		})

		Context("when network type is dynamic", func() {
			BeforeEach(func() {
				network = Network{
					Name: "fake-dynamic-network-name",
					Type: "dynamic",
					CloudProperties: biproperty.Map{
						"cp_key": "cp_value",
					},
					DNS: []string{"2.2.2.2"},
				}
			})

			It("includes dns and cloud_properties", func() {
				iface, err := network.Interface([]string{}, []NetworkDefault{})
				Expect(err).ToNot(HaveOccurred())
				Expect(iface).To(Equal(biproperty.Map{
					"type": "dynamic",
					"dns":  []string{"2.2.2.2"},
					"cloud_properties": biproperty.Map{
						"cp_key": "cp_value",
					},
				}))
			})
		})
	})
})
