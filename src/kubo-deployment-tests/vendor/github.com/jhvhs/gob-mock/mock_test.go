package gobmock

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Mock", func() {

	It("includes the stub function", func() {
		mock := Mock("jimjam", "").MockContents()
		Expect(mock).To(ContainSubstring("jimjam() {"))
	})

	It("includes the pipe handling", func() {
		mock := Mock("helicopter", "").MockContents()
		Expect(mock).To(ContainSubstring("while read -r -t0.1; do"))

	})

	It("includes the custom mock script", func() {
		mock := Mock("visitor", "cakes and coffee").MockContents()
		Expect(mock).To(ContainSubstring("cakes and coffee"))
	})

	It("can conditionally call through", func() {
		mock := MockOrCallThrough("printf", "echo 'starships in bottles'", "[ $1 == 'monkey' ]")
		Expect(mock.MockContents()).To(ContainSubstring("monkey"))
	})

})
