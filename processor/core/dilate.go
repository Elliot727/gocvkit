// Package core provides morphological dilation filter implementation.
//
// The dilation filter applies a morphological dilation operation with configurable
// kernel size and number of iterations.
package core

import (
	"image"

	"github.com/Elliot727/gocvkit/processor"
	"gocv.io/x/gocv"
)

// Dilate defines the configuration for morphological dilation.
type Dilate struct {
	KernelSize int `toml:"kernel"`      // KernelSize is the size of the structuring element for dilation
	Iterations int `toml:"iterations"`  // Iterations is the number of times dilation is applied
}

// Process applies morphological dilation using the configured parameters.
func (d *Dilate) Process(src gocv.Mat, dst *gocv.Mat) {
	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(d.KernelSize, d.KernelSize))
	defer kernel.Close()
	gocv.Dilate(src, dst, kernel)
	for i := 1; i < d.Iterations; i++ {
		temp := gocv.NewMat()
		defer temp.Close()
		gocv.Dilate(*dst, &temp, kernel)
		temp.CopyTo(dst)
	}
}

func init() {
	processor.Register("Dilate", &Dilate{
		KernelSize: 3,
		Iterations: 1,
	})
}
