// Package core provides morphological erosion filter implementation.
//
// The erosion filter applies a morphological erosion operation with configurable
// kernel size and number of iterations.
package core

import (
	"fmt"
	"image"

	"github.com/Elliot727/gocvkit/processor"
	"gocv.io/x/gocv"
)

// Erode defines the configuration for morphological erosion.
type Erode struct {
	KernelSize int `toml:"kernel"`     // KernelSize is the size of the structuring element for erosion
	Iterations int `toml:"iterations"` // Iterations is the number of times erosion is applied
	// Pre-allocated resources
	kernel gocv.Mat
}

// Validate checks constraints and pre-allocates the kernel.
func (e *Erode) Validate() error {
	if e.KernelSize < 1 {
		return fmt.Errorf("kernel size must be >= 1, got %d", e.KernelSize)
	}
	if e.KernelSize%2 == 0 {
		return fmt.Errorf("kernel size must be odd, got %d", e.KernelSize)
	}
	if e.Iterations < 1 {
		return fmt.Errorf("iterations must be >= 1, got %d", e.Iterations)
	}

	e.kernel = gocv.GetStructuringElement(gocv.MorphRect, image.Pt(e.KernelSize, e.KernelSize))
	if e.kernel.Empty() {
		return fmt.Errorf("failed to create structuring element")
	}
	return nil
}

// Process applies morphological erosion using the configured parameters.
func (e *Erode) Process(src gocv.Mat, dst *gocv.Mat) error {
	if src.Empty() {
		return nil
	}

	gocv.Erode(src, dst, e.kernel)

	if e.Iterations > 1 {
		temp := gocv.NewMat()
		defer temp.Close()

		for i := 1; i < e.Iterations; i++ {
			gocv.Erode(*dst, &temp, e.kernel)
			temp.CopyTo(dst)
		}
	}
	return nil
}
func (e *Erode) Close() {
	if !e.kernel.Empty() {
		e.kernel.Close()
	}

}

func init() {
	processor.Register("Erode", &Erode{
		KernelSize: 3,
		Iterations: 1,
	})
}
