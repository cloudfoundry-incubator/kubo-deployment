package gobmock

import (
	"fmt"

	"github.com/tonnerre/golang-text"
)

const unconditionalCallthrough = "📣"

// Produces a bash function with a given name.
// The function will report it's arguments
// as well as any data passed into it via STDIN.
// All reporting messages are sent to STDERR.
// The function is exported for use in child processes
func Spy(name string) Gob {
	return &spy{name: name, callThroughCondition: "", shouldExport: true}
}

// Produces a bash function with a given name.
// The function will report it's arguments
// as well as any data passed into it via STDIN.
// All reporting messages are sent to STDERR.
func ShallowSpy(name string) Gob {
	return &spy{name: name, callThroughCondition: "", shouldExport: false}
}

// Produces a bash function with a given name.
// The function will report it's arguments
// as well as any data passed into it via STDIN.
// All reporting messages are sent to STDERR.
// This function will also call through to the
// original executable
// The function is exported for use in child processes
func SpyAndCallThrough(name string) Gob {
	return &spy{name: name, callThroughCondition: unconditionalCallthrough, shouldExport: true}
}

// Produces a bash function with a given name.
// The function will report it's arguments
// as well as any data passed into it via STDIN.
// All reporting messages are sent to STDERR.
// This function will also call through to the
// original executable
func ShallowSpyAndCallThrough(name string) Gob {
	return &spy{name: name, callThroughCondition: unconditionalCallthrough, shouldExport: false}
}

// Produces a bash function with a given name.
// The function will report it's arguments
// as well as any data passed into it via STDIN.
// All reporting messages are sent to STDERR.
// This function will also call through to the
// original executable when a supplied condition is met
// The function is exported for use in child processes
func SpyAndConditionallyCallThrough(name string, callThroughCondition string) Gob {
	return &spy{name: name, callThroughCondition: callThroughCondition, shouldExport: true}
}

// Produces a bash function with a given name.
// The function will report it's arguments
// as well as any data passed into it via STDIN.
// All reporting messages are sent to STDERR.
// This function will also call through to the
// original executable when a supplied condition is met
// The function is exported for use in child processes
func ShallowSpyAndConditionallyCallThrough(name string, callThroughCondition string) Gob {
	return &spy{name: name, callThroughCondition: callThroughCondition, shouldExport: false}
}

type spy struct {
	name                 string
	callThroughCondition string
	shouldExport         bool
}

func (s *spy) MockContents() string {
	if s.shouldExport {
		return s.spyFunction() + s.spyExport()
	}
	return s.spyFunction()
}

func (s *spy) spyExport() string {
	return fmt.Sprintf(exportDefinition, s.name)
}

func (s *spy) spyFunction() string {
	script := scriptStart + spyDefinition
	if s.callThroughCondition == unconditionalCallthrough {
		script = script + callThroughDefinition
	} else if s.callThroughCondition != "" {
		script = script + s.conditionalCallThrough()
	}
	return fmt.Sprintf(script+scriptEnd, s.name)
}

func (s *spy) conditionalCallThrough() string {
	return text.Indent("\nif "+s.callThroughCondition+"; then\n"+callThroughDefinition+"fi\n", "  ")
}
