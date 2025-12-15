// Package config handles loading and parsing of TOML configuration files.
//
// It provides the main Config struct that represents the complete application
// configuration loaded from a TOML file. The package includes custom
// unmarshaling logic to handle dynamic pipeline step parameters efficiently.
package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// Config represents the complete application configuration.
type Config struct {
	App struct {
		WindowName string `toml:"window_name"` // WindowName is the title for the display window
		Record     bool   `toml:"record"`      // Record enables video recording when set to true
		Output     string `toml:"output"`      // Output is the path for the recorded video file
	} `toml:"app"`

	Camera struct {
		DeviceID int    `toml:"device_id"` // DeviceID is the camera device index (ignored if File is set)
		File     string `toml:"file"`      // File is the path to a video file (takes precedence over DeviceID)
	} `toml:"camera"`

	Stream struct {
		Enabled bool   `toml:"enabled"`
		Port    int    `toml:"port"`
		Path    string `toml:"path"`
		Quality int    `toml:"quality"`
	}

	Pipeline struct {
		Steps []StepConfig `toml:"steps"` // Steps contains the ordered list of processing steps
	} `toml:"pipeline"`
}

// StepConfig holds the name and a map of ALL other parameters.
// We removed the struct tags because we are using UnmarshalTOML below.
type StepConfig struct {
	Name   string                 // Name of the processor step
	Params map[string]interface{} // Params contains all additional configuration parameters
}

// UnmarshalTOML is a hook called automatically by the TOML parser.
// It gives us the raw map, allowing us to manually extract 'name'
// and keep everything else as params.
func (s *StepConfig) UnmarshalTOML(data interface{}) error {
	// 1. Cast the raw data to a map
	raw, ok := data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("expected TOML table for step, got %T", data)
	}

	// 2. Extract the Name field manually
	if name, ok := raw["name"].(string); ok {
		s.Name = name
		// Remove 'name' from the map so it doesn't appear in Params
		delete(raw, "name")
	} else {
		return fmt.Errorf("pipeline step missing 'name' field")
	}

	// 3. Assign the remaining fields to Params
	s.Params = raw
	return nil
}

// Load reads and parses the TOML configuration file at the given path.
// Returns a Config struct with default values applied if not present in the file.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.App.WindowName == "" {
		cfg.App.WindowName = "GoCV Live"
	}

	return &cfg, nil
}
