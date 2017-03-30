package kubo_deployment_tests_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
	"runtime"
	"path/filepath"
)

func pathToScript(name string) string {
	return pathFromRoot("bin/" + name)
}

func pathFromRoot(relativePath string) string {
	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)
	return filepath.Join(currentDir, "..", "..", relativePath)
}

var EmptyCallback = func([]string) {}

func TestKuboDeploymentTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "KuboDeploymentTests Suite")
}
