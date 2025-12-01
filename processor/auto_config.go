package processor

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/Elliot727/gocvkit/config"

	"github.com/BurntSushi/toml"
	"gocv.io/x/gocv"
)

// autoWrapper wraps a user's Processable to add the Name() method required by Step.
// This allows users to implement just Process() without worrying about names.
type autoWrapper struct {
	name string
	impl Processable
}

func (a *autoWrapper) Name() string { return a.name }

func (a *autoWrapper) Process(src gocv.Mat, dst *gocv.Mat) {
	a.impl.Process(src, dst)
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

		// 3. Return the wrapped Step
		return &autoWrapper{
			name: cfg.Name, // Use the name from the config file
			impl: newStepPtr.Interface().(Processable),
		}, nil
	}
}
