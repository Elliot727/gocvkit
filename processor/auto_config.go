// Package processor provides utilities for automatic configuration of processing steps.
//
// The auto_config module implements the reflection-based mechanism that allows
// users to define simple Processable structs with TOML tags, which are then
// automatically configured based on the parameters in the TOML configuration file.
package processor

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/Elliot727/gocvkit/config"

	"github.com/BurntSushi/toml"
	"gocv.io/x/gocv"
)

type Validator interface {
	Validate() error
}

// autoWrapper wraps a user's Processable to add the Name() method required by Step.
// This allows users to implement just Process() without worrying about names.
type autoWrapper struct {
	name string
	impl Processable
}

// Name returns the name of the processing step.
func (a *autoWrapper) Name() string { return a.name }

// Process executes the processing step on the provided source and destination matrices.
func (a *autoWrapper) Process(src gocv.Mat, dst *gocv.Mat) error {
	return a.impl.Process(src, dst)
}

func (a *autoWrapper) Close() {
	// Check if the underlying struct has a Close() method
	if c, ok := a.impl.(interface{ Close() }); ok {
		c.Close()
	}
}

// AutoConfig generates a Factory that creates configured instances of the provided default struct.
func AutoConfig(defaults Processable) Factory {
	return func(cfg config.StepConfig) (Step, error) {
		// 1. Reflection: Create a fresh pointer to a copy of the defaults.
		// We dereference the pointer to get the struct value, then make a new pointer to that type.
		val := reflect.ValueOf(defaults)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		// Create new instance of the specific processor type (e.g. *Canny)
		newStepPtr := reflect.New(val.Type())

		// Copy the default values (e.g. Low=50) into the new instance
		newStepPtr.Elem().Set(val)

		// 2. Decode TOML params into the new struct
		// We re-encode the map to TOML and decode it back into the struct
		// so that standard `toml:"tag"` works perfectly on the user's struct.

		if len(cfg.Params) > 0 {
			var buf bytes.Buffer
			enc := toml.NewEncoder(&buf)

			// Encode the generic map back to TOML format
			if err := enc.Encode(cfg.Params); err != nil {
				return nil, fmt.Errorf("failed to process params for %s: %w", cfg.Name, err)
			}

			// Decode that TOML into the specific struct fields
			if _, err := toml.Decode(buf.String(), newStepPtr.Interface()); err != nil {
				return nil, fmt.Errorf("invalid parameters for processor %q: %w", cfg.Name, err)
			}
		}
		step := newStepPtr.Interface()

		if v, ok := step.(Validator); ok {
			if err := v.Validate(); err != nil {
				return nil, fmt.Errorf("validation failed for %q: %w", cfg.Name, err)
			}
		}
		proc, _ := step.(Processable)

		// 3. Return the wrapped Step
		return &autoWrapper{
			name: cfg.Name, // Use the name from the config file
			impl: proc,
		}, nil
	}
}
