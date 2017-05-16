package gobmock

import "fmt"

// Produces a bash function with a given name.
func Stub(name string) Gob {
	return &stub{name: name}
}

type stub struct {
	name string
}

func (s *stub) MockContents() string {
	return s.stubFunction() + s.stubExport()
}

func (s *stub) stubExport() string {
	return fmt.Sprintf(exportDefinition, s.name)
}

func (s *stub) stubFunction() string {
	script := scriptStart + stubDefinition + scriptEnd
	return fmt.Sprintf(script, s.name)
}
