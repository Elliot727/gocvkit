// Package core provides morphological dilation filter implementation.
//
// The dilation filter applies a morphological dilation operation with configurable
// kernel size and number of iterations.
package core

import (
	"fmt"
	"image"

	"github.com/Elliot727/gocvkit/processor"
	"gocv.io/x/gocv"
)

// Dilate defines the configuration for morphological dilation.
type Dilate struct {
	KernelSize int `toml:"kernel"`     // KernelSize is the size of the structuring element for dilation
	Iterations int `toml:"iterations"` // Iterations is the number of times dilation is applied
	// Pre-allocated kernel to avoid recreation every frame
	kernel gocv.Mat
}

// Validate checks constraints and pre-allocates the kernel.
func (d *Dilate) Validate() error {
	if d.KernelSize < 1 {
		return fmt.Errorf("kernel size must be >= 1, got %d", d.KernelSize)
	}
	if d.KernelSize%2 == 0 {
		return fmt.Errorf("kernel size must be odd, got %d", d.KernelSize)
	}
	if d.Iterations < 1 {
		return fmt.Errorf("iterations must be >= 1, got %d", d.Iterations)
	}

	// Pre-create the kernel once.
	// If this fails, we fail at startup, not during video processing.
	d.kernel = gocv.GetStructuringElement(gocv.MorphRect, image.Pt(d.KernelSize, d.KernelSize))
	if d.kernel.Empty() {
		return fmt.Errorf("failed to create structuring element")
	}
	return nil
}

// Process applies morphological dilation using the configured parameters.
func (d *Dilate) Process(src gocv.Mat, dst *gocv.Mat) error {
	if src.Empty() {
		return nil
	}

	gocv.Dilate(src, dst, d.kernel)

	if d.Iterations > 1 {
		temp := gocv.NewMat()
		defer temp.Close()

		for i := 1; i < d.Iterations; i++ {
			gocv.Dilate(*dst, &temp, d.kernel)
			temp.CopyTo(dst)
		}
	}

	return nil
}

func (d *Dilate) Close() {
	if !d.kernel.Empty() {
		d.kernel.Close()
	}
}

func init() {
	processor.Register("Dilate", &Dilate{
		KernelSize: 3,
		Iterations: 1,
	})
}
