package main_test

import (
	"github.com/progrium/go-basher"
	"os"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("set_kubeconfig", func(){
	It("fails without arguments", func() {
		bash, _ := basher.NewContext("/bin/bash", false)
		bash.Stdout = GinkgoWriter
		bash.Stderr = GinkgoWriter
		bash.Source("../../bin/set_kubeconfig", nil)
		status, err := bash.Run("main", os.Args[1:])
		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal(1))

	})
})