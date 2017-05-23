package gobmock

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Stub", func() {
	It("is created", func() {
		Expect(Stub("foo")).NotTo(BeNil())

	})

	It("includes a basic shell stub", func() {
		stub := Stub("devatio-crederes").MockContents()
		Expect(stub).To(MatchRegexp("^\n# Gob\ndevatio-crederes\\(\\)\\s*{"))
		Expect(stub).To(MatchRegexp("while read -r -t0.1; do\\s+:\\s+done"))
	})

})
