package gobmock

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Spy", func() {

	It("is created", func() {
		Expect(Spy("glass")).NotTo(BeNil())
	})

	It("includes a basic shell spy", func() {
		spy := Spy("chicken-with-a-pulley").MockContents()
		Expect(spy).To(MatchRegexp("^\n# Gob\nchicken-with-a-pulley\\(\\) {\\s+"))
		Expect(spy).To(MatchRegexp("echo.* > /dev/fd/2"))
	})

	It("includes a pipe consumer", func() {
		spy := Spy("grass").MockContents()
		Expect(spy).To(ContainSubstring("while read -r -t0.1; do"))
	})

	It("does not call through by default", func() {
		spy := Spy("bee").MockContents()
		Expect(spy).NotTo(ContainSubstring("$(which bee)"))
	})

	It("can call through to the executable", func() {
		spy := SpyAndCallThrough("squash").MockContents()
		Expect(spy).To(ContainSubstring("$(which squash)"))
	})

	It("can include a condition for a call through", func() {
		spy := SpyAndConditionallyCallThrough("raspbery", "[[ 3 -ne 4 ]]")
		Expect(spy.MockContents()).To(ContainSubstring("if [[ 3 -ne 4 ]]; then"))
	})

	It("is exported by default", func() {
		spy := Spy("on-me").MockContents()
		Expect(spy).To(ContainSubstring("export -f on-me"))
	})
	Context("shallow spy", func() {
		It("does not export a regular spy", func() {
			spy := ShallowSpy("on-you").MockContents()
			Expect(spy).NotTo(ContainSubstring("export -f on-you"))
		})

		It("doesn't export a call through spy", func() {
			spy := ShallowSpyAndCallThrough("vixen").MockContents()
			Expect(spy).NotTo(ContainSubstring("export -f vixen"))
		})

		It("doesn't export a conditionally call through spy", func() {
			spy := ShallowSpyAndConditionallyCallThrough("shutter", "[[ 1 -eq 2 ]]").MockContents()
			Expect(spy).NotTo(ContainSubstring("export -f shutter"))
		})
	})
})
