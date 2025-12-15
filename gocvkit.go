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

// Grayscale is an alias for core.Grayscale, providing color to grayscale conversion.
type Grayscale = core.Grayscale

// Erode is an alias for core.Erode, providing morphological erosion filtering.
type Erode = core.Erode

// Dilate is an alias for core.Dilate, providing morphological dilation filtering.
type Dilate = core.Dilate

// MorphClose is an alias for core.MorphClose, providing morphological closing operation.
type MorphClose = core.MorphClose

// Otsu is an alias for core.Otsu, providing Otsu threshold filtering.
type Otsu = core.Otsu

// Canny is an alias for edges.Canny, providing Canny edge detection.
type Canny = edges.Canny

// Sobel is an alias for edges.Sobel, providing Sobel edge detection.
type Sobel = edges.Sobel

// Laplacian is an alias for edges.Laplacian, providing Laplacian edge detection.
type Laplacian = edges.Laplacian

// Scharr is an alias for edges.Scharr, providing Scharr edge detection.
type Scharr = edges.Scharr

// GaussianBlur is an alias for blurs.GaussianBlur, providing Gaussian blur filtering.
type GaussianBlur = blurs.GaussianBlur

// MedianBlur is an alias for blurs.MedianBlur, providing median blur filtering.
type MedianBlur = blurs.MedianBlur

// Bilateral is an alias for blurs.Bilateral, providing bilateral filtering.
type Bilateral = blurs.Bilateral
