// Package processor provides a registry and interface for image processing steps.
//
// The package implements a plugin-style architecture that allows for registering
// custom image processing filters that can be configured via TOML files.
package processor

import (
	"fmt"

	"github.com/Elliot727/gocvkit/config"

	"gocv.io/x/gocv"
)

// Processable is the simplified interface for user-defined filters.
// You only need to implement this.
type Processable interface {
	Process(src gocv.Mat, dst *gocv.Mat) error
}

// Step is the internal interface used by the Pipeline.
// It combines the user's logic with the system's need for a Name.
type Step interface {
	Processable
	Name() string
	Close()
}

// Factory is a function that creates a Step from configuration.
type Factory func(config.StepConfig) (Step, error)

// registry is private to ensure thread-safety and prevent external tampering.
var registry = make(map[string]Factory)

// Register adds a new processor. It is smart and accepts two types:
// 1. A struct instance (Processable): Automatically wrapped with AutoConfig.
// 2. A factory function: For complex setup logic.
func Register(name string, item interface{}) {
	switch v := item.(type) {
	case Processable:
		// The Magic: We auto-wrap the struct here!
		registry[name] = AutoConfig(v)
	case func(config.StepConfig) (Step, error):
		registry[name] = v
	default:
		panic(fmt.Sprintf("processor.Register: %q must be a Processable struct or a Factory func", name))
	}
}

// Get looks up a processor factory by name.
func Get(name string) (Factory, bool) {
	f, ok := registry[name]
	return f, ok
}
