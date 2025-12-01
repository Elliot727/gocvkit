// Package blurs provides Gaussian blur filter implementation.
//
// The Gaussian blur filter applies a Gaussian kernel with configurable size and sigma.
package blurs

import (
	"image"

	"github.com/Elliot727/gocvkit/processor"
	"gocv.io/x/gocv"
)

// GaussianBlur performs Gaussian filtering with configurable kernel and sigma.
type GaussianBlur struct {
	Kernel int     `toml:"kernel"` // Kernel is the size of the Gaussian kernel (will be made odd if even)
	Sigma  float64 `toml:"sigma"`  // Sigma is the standard deviation for Gaussian kernel
}

// Process applies Gaussian blur using the configured parameters.
func (g *GaussianBlur) Process(src gocv.Mat, dst *gocv.Mat) {
	k := g.Kernel

	// Enforce positive odd kernel size
	if k < 1 {
		k = 1
	}
	if k%2 == 0 {
		k++
	}

	gocv.GaussianBlur(src, dst, image.Pt(k, k), g.Sigma, g.Sigma, gocv.BorderDefault)
}

func init() {
	// Register directly with default values.
	// The smart Register function handles the wrapping.
	processor.Register("GaussianBlur", &GaussianBlur{
		Kernel: 9,
		Sigma:  1.8,
	})
}
