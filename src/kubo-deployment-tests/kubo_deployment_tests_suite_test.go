package kubo_deployment_tests_test

import (
	"io"
	"log"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	basher "github.com/progrium/go-basher"

	"path/filepath"
	"runtime"
	"testing"
)

var (
	bash   *basher.Context
	stdout *gbytes.Buffer
	stderr *gbytes.Buffer

	resourcesPath   string
	environmentPath string

	emptyCallback = func([]string) {}
	bashPath      string
)

func pathToScript(name string) string {
	return pathFromRoot("bin/" + name)
}

func pathFromRoot(relativePath string) string {
	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)
	return filepath.Join(currentDir, "..", "..", relativePath)
}

func TestKuboDeploymentTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "KuboDeploymentTests Suite")
}

var _ = BeforeSuite(func() {
	extractBash()
})

var _ = BeforeEach(func() {
	resourcesPath = filepath.Join(pathFromRoot("src"), "odb-deployment", "resources")
	environmentPath = filepath.Join(resourcesPath, "environment")

	bash, _ = basher.NewContext(bashPath, true)

	stdout = gbytes.NewBuffer()
	stderr = gbytes.NewBuffer()
	bash.Stdout = io.MultiWriter(GinkgoWriter, stdout)
	bash.Stderr = io.MultiWriter(GinkgoWriter, stderr)
	bash.Source("_", func(string) ([]byte, error) {
		return []byte(`
				callCounter=0
				invocationRecorder() {
				  local in_line_count=0
				  declare -a in_lines
				  while read -t0.05; do
				    in_lines[in_line_count]="$REPLY"
				    in_line_count=$(expr ${in_line_count} + 1)
				  done
				  callCounter=$(expr ${callCounter} + 1)
				  (>&2 echo "[$callCounter] $@")
				  if [ ${in_line_count} -gt 0 ]; then
				    (>&2 echo "[$callCounter received] input:")
				    (>&2 printf '%s\n' "${in_lines[@]}" )
				    (>&2 echo "[$callCounter end received]")
				  fi
				}
			`), nil
	})

	bash.CopyEnv()
})

var _ = AfterSuite(func() {
	os.Remove(bashPath)
})

func extractBash() {
	bashDir, err := homedir.Expand("~/.basher")
	if err != nil {
		log.Fatal(err, "1")
	}

	bashPath = bashDir + "/bash"
	if _, err := os.Stat(bashPath); os.IsNotExist(err) {
		err = basher.RestoreAsset(bashDir, "bash")
		if err != nil {
			log.Fatal(err, "1")
		}
	}
}
