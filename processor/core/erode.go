// Package core provides morphological erosion filter implementation.
//
// The erosion filter applies a morphological erosion operation with configurable
// kernel size and number of iterations.
package core

import (
	"image"

	"github.com/Elliot727/gocvkit/processor"
	"gocv.io/x/gocv"
)

// Erode defines the configuration for morphological erosion.
type Erode struct {
	KernelSize int `toml:"kernel"`      // KernelSize is the size of the structuring element for erosion
	Iterations int `toml:"iterations"`  // Iterations is the number of times erosion is applied
}

// Process applies morphological erosion using the configured parameters.
func (e *Erode) Process(src gocv.Mat, dst *gocv.Mat) {
	k := e.KernelSize
	if k < 1 {
		k = 1
	}
	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(k, k))
	defer kernel.Close()

	gocv.Erode(src, dst, kernel)
	for i := 1; i < e.Iterations; i++ {
		tmp := gocv.NewMat()
		defer tmp.Close()
		gocv.Erode(*dst, &tmp, kernel)
		tmp.CopyTo(dst)
	}
}

func init() {
	processor.Register("Erode", &Erode{
		KernelSize: 3,
		Iterations: 1,
	})
}
