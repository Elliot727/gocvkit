// Package blurs provides Gaussian blur filter implementation.
//
// The Gaussian blur filter applies a Gaussian kernel with configurable size and sigma.
package blurs

import (
	"fmt"
	"image"

	"github.com/Elliot727/gocvkit/processor"
	"gocv.io/x/gocv"
)

// GaussianBlur performs Gaussian filtering with configurable kernel and sigma.
type GaussianBlur struct {
	Kernel int     `toml:"kernel"` // Kernel is the size of the Gaussian kernel (will be made odd if even)
	Sigma  float64 `toml:"sigma"`  // Sigma is the standard deviation for Gaussian kernel
}

// Validate checks constraints before the pipeline starts.
func (g *GaussianBlur) Validate() error {
	if g.Kernel < 1 {
		return fmt.Errorf("kernel must be >= 1, got %d", g.Kernel)
	}
	if g.Kernel%2 == 0 {
		return fmt.Errorf("kernel must be odd, got %d (auto-fixing is magic, don't do it)", g.Kernel)
	}
	if g.Sigma < 0 {
		return fmt.Errorf("sigma must be >= 0, got %f", g.Sigma)
	}
	return nil
}

// Process applies Gaussian blur using the configured parameters.
func (g *GaussianBlur) Process(src gocv.Mat, dst *gocv.Mat) error {
	if src.Empty() {
		return nil
	}
	gocv.GaussianBlur(src, dst, image.Pt(g.Kernel, g.Kernel), g.Sigma, g.Sigma, gocv.BorderDefault)
	return nil
}

func (g *GaussianBlur) Close() {}

func init() {
	// Register directly with default values.
	// The smart Register function handles the wrapping.
	processor.Register("GaussianBlur", &GaussianBlur{
		Kernel: 9,
		Sigma:  1.8,
	})
}
