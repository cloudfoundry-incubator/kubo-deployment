package kubo_deployment_tests_test

import (
	boshtpl "github.com/cloudfoundry/bosh-cli/director/template"
	"github.com/cppforlife/go-patch/patch"
	"gopkg.in/yaml.v2"
)

func propertyFromYaml(path string, yamlContents []byte) (string, error) {
	var yamlMap map[string]interface{}
	err := yaml.Unmarshal(yamlContents, &yamlMap)
	if err != nil {
		return "", err
	}

	template := boshtpl.NewTemplate(yamlContents)
	vars := boshtpl.StaticVariables{}
	ops := patch.FindOp{Path: patch.MustNewPointerFromString(path)}

	bytes, err := template.Evaluate(vars, ops, boshtpl.EvaluateOpts{ExpectAllKeys: false})
	return choppedString(bytes), err
}

func choppedString(bytes []byte) string {
	if len(bytes) > 0 {
		return string(bytes[:len(bytes)-1])
	} else {
		return ""
	}
}
