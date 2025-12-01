// Package gocvkit is a modular GoCV framework for working with live camera
// or video streams with configurable pipelines and filters.
//
// Features:
//   - Modular processing pipeline (Grayscale, Blur, Canny, Sobel)
//   - TOML configuration for camera, display, and processor settings
//   - Supports webcam or video file input
//   - Display window with optional frame callbacks
//   - Easy extension with custom filters
//
// Example usage:
//
//	import "gocvkit"
//
//	func main() {
//	    app, err := gocvkit.NewApp("config.toml")
//	    if err != nil {
//	        log.Fatalf("failed to start app: %v", err)
//	    }
//	    defer app.Close()
//
//	    // Run with optional per-frame callback
//	    app.Run(nil)
//	}
package gocvkit

import (
	"github.com/Elliot727/gocvkit/app"
	"github.com/Elliot727/gocvkit/processor"
)

// NewApp is a top-level convenience function that creates and returns
// a fully configured App instance from a TOML config path.
func NewApp(cfgPath string) (*app.App, error) {
	return app.New(cfgPath)
}

// RegisterProcessor is a top-level convenience function that allows you
// to regiser processes externally.
func RegisterProcessor(name string, item any) {
	processor.Register(name, item)
}
