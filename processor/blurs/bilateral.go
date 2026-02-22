// Package blurs provides bilateral filter implementation.
//
// The bilateral filter smooths images while preserving edges, using configurable
// diameter and sigma values for color and space.
package blurs

import (
	"fmt"

	"github.com/Elliot727/gocvkit/processor"
	"gocv.io/x/gocv"
)

// Bilateral defines the configuration for bilateral filtering.
type Bilateral struct {
	Diameter   int     `toml:"diameter"`    // Diameter is the diameter of each pixel neighborhood
	SigmaColor float64 `toml:"sigma_color"` // SigmaColor is the filter sigma in the color space
	SigmaSpace float64 `toml:"sigma_space"` // SigmaSpace is the filter sigma in the coordinate space
}

// Validate checks constraints before the pipeline starts.
func (b *Bilateral) Validate() error {
	if b.Diameter < 0 {
		return fmt.Errorf("diameter must be >= 0, got %d", b.Diameter)
	}
	if b.SigmaColor <= 0 {
		return fmt.Errorf("sigma_color must be > 0, got %f", b.SigmaColor)
	}
	if b.SigmaSpace <= 0 {
		return fmt.Errorf("sigma_space must be > 0, got %f", b.SigmaSpace)
	}
	return nil
}

// Process applies bilateral filtering using the configured parameters.
func (b *Bilateral) Process(src gocv.Mat, dst *gocv.Mat) error {
	if src.Empty() {
		return nil
	}
	gocv.BilateralFilter(src, dst, b.Diameter, b.SigmaColor, b.SigmaSpace)
	return nil
}

func (b *Bilateral) Close() {}

func init() {
	processor.Register("Bilateral", &Bilateral{
		Diameter:   9,
		SigmaColor: 75,
		SigmaSpace: 75,
	})
}
