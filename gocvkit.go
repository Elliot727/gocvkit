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

	// Import sub-packages to alias them.
	// This automatically runs their init() functions, so "Grayscale" is registered.
	"github.com/Elliot727/gocvkit/processor/blurs"
	"github.com/Elliot727/gocvkit/processor/core"
	"github.com/Elliot727/gocvkit/processor/edges"
)

// NewApp creates a fully configured App instance from a TOML config path.
func NewApp(cfgPath string) (*app.App, error) {
	return app.New(cfgPath)
}

// RegisterProcessor allows external registration of custom processes.
func RegisterProcessor(name string, item any) {
	processor.Register(name, item)
}

// ---------------------------------------------------------
// EXPORTED FILTERS (Type Aliases)
// ---------------------------------------------------------
// This allows users to use 'gocvkit.Canny' instead of 'edges.Canny'.

type Grayscale = core.Grayscale
type Erode = core.Erode
type Dilate = core.Dilate
type MorphClose = core.MorphClose
type Otsu = core.Otsu

type Canny = edges.Canny
type Sobel = edges.Sobel
type Laplacian = edges.Laplacian
type Scharr = edges.Scharr

type GaussianBlur = blurs.GaussianBlur
type MedianBlur = blurs.MedianBlur
type Bilateral = blurs.Bilateral
