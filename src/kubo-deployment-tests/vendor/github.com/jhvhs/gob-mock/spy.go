package gobmock

import "fmt"

// Produces a bash function with a given name.
// The function will report it's arguments
// as well as any data passed into it via STDIN.
// All reporting messages are sent to STDERR.
func Spy(name string) Gob {
	return &spy{name: name, withCallThrough: false}
}

func SpyAndCallThrough(name string) Gob {
	return &spy{name: name, withCallThrough: true}
}

type spy struct {
	name            string
	withCallThrough bool
}

func (s *spy) MockContents() string {
	return s.spyFunction() + s.spyExport()
}

func (s *spy) spyExport() string {
	return fmt.Sprintf(exportDefinition, s.name)
}

func (s *spy) spyFunction() string {
	script := scriptStart + spyDefinition
	if s.withCallThrough {
		script = script + callThroughDefinition
	}
	return fmt.Sprintf(script + scriptEnd, s.name)
}
