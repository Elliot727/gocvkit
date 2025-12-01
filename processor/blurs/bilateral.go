// Package blurs provides bilateral filter implementation.
//
// The bilateral filter smooths images while preserving edges, using configurable
// diameter and sigma values for color and space.
package blurs

import (
	"github.com/Elliot727/gocvkit/processor"
	"gocv.io/x/gocv"
)

// Bilateral defines the configuration for bilateral filtering.
type Bilateral struct {
	Diameter   int     `toml:"diameter"`   // Diameter is the diameter of each pixel neighborhood
	SigmaColor float64 `toml:"sigma_color"` // SigmaColor is the filter sigma in the color space
	SigmaSpace float64 `toml:"sigma_space"` // SigmaSpace is the filter sigma in the coordinate space
}

// Process applies bilateral filtering using the configured parameters.
func (b *Bilateral) Process(src gocv.Mat, dst *gocv.Mat) {
	gocv.BilateralFilter(src, dst, b.Diameter, b.SigmaColor, b.SigmaSpace)
}

func init() {
	processor.Register("Bilateral", &Bilateral{
		Diameter:   9,
		SigmaColor: 75,
		SigmaSpace: 75,
	})
}
