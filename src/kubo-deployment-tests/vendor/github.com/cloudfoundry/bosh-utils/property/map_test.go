package property_test

import (
	. "github.com/cloudfoundry/bosh-utils/property"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
)

var _ = Describe("Map", func() {
	It("can be unmarshaled from yaml", func() {
		expectedMap := Map{"foo": Map{"bar": List{"baz", "asdf"}}}

		inputYAML := `
foo:
  bar:
    - baz
    - asdf
`
		target := Map{}
		err := yaml.Unmarshal([]byte(inputYAML), &target)
		Expect(err).ToNot(HaveOccurred())
		Expect(target).To(Equal(expectedMap))
	})
})
