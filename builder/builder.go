// Package builder creates a processing pipeline from the loaded configuration
// by looking up each step in the processor registry.
//
// This is the glue that makes the entire plugin-style architecture work:
// config → builder → []processor.Step → pipeline.Pipeline
package builder

import (
	"fmt"

	"github.com/Elliot727/gocvkit/config"
	"github.com/Elliot727/gocvkit/processor"
)

// BuildPipeline constructs the ordered list of processing steps from the config.
// Returns an error if any step name is unknown or its factory fails.
func BuildPipeline(cfg *config.Config) ([]processor.Step, error) {
	var steps []processor.Step

	for i, sc := range cfg.Pipeline.Steps {
		factory, ok := processor.Get(sc.Name)
		if !ok {
			return nil, fmt.Errorf("pipeline step %d: unknown processor %q", i, sc.Name)
		}

		step, err := factory(sc)
		if err != nil {
			return nil, fmt.Errorf("pipeline step %d (%s): %w", i, sc.Name, err)
		}

		steps = append(steps, step)
	}

	return steps, nil
}
