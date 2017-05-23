package gobmock

import (
	"fmt"
	"github.com/tonnerre/golang-text"
)

// Produces a bash function with a given name.
// The function will report it's arguments
// as well as any data passed into it via STDIN.
// All reporting messages are sent to STDERR.
// The `mockScript` string will be inserted at the end.
func Mock(name string, mockScript string) Gob {
	return &mock{name: name, script: mockScript, condition: ""}
}

func MockOrCallThrough(name string, mockScript string, callThroughCondition string) Gob {
	return &mock{name: name, script: mockScript, condition: callThroughCondition}
}

type mock struct {
	name      string
	script    string
	condition string
}

func (m *mock) MockContents() string {
	return m.mockFunction() + m.mockExport()
}

func (m *mock) mockExport() string {
	return fmt.Sprintf(exportDefinition, m.name)
}

func (m *mock) mockFunction() string {
	script := scriptStart + spyDefinition + m.mockBody() + scriptEnd
	return fmt.Sprintf(script, m.name, m.script)
}

func (m *mock) mockBody() string {
	if m.condition != "" {
		return m.callThrough()
	}
	return mockDefinition
}

func (m *mock)callThrough() string {
	return text.Indent("\nif " + m.condition + "; then\n" + callThroughDefinition + "else\n" + mockDefinition + "fi\n", "  ")
}
