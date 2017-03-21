package kubo_deployment_tests_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestKuboDeploymentTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "KuboDeploymentTests Suite")
}
