package gobmock_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGobmock(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gobmock Suite")
}
