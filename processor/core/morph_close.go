// Package core provides morphological close operation filter implementation.
//
// The morphological close operation is useful for closing small holes and gaps in objects.
package core

import (
	"image"

	"github.com/Elliot727/gocvkit/processor"
	"gocv.io/x/gocv"
)

// MorphClose defines the configuration for morphological close operation.
type MorphClose struct {
	KernelSize int `toml:"kernel"`      // KernelSize is the size of the structuring element for morphological close
	Iterations int `toml:"iterations"`  // Iterations is the number of times morphological close is applied
}

// Process applies morphological close operation using the configured parameters.
func (m *MorphClose) Process(src gocv.Mat, dst *gocv.Mat) {
	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(m.KernelSize, m.KernelSize))
	defer kernel.Close()
	gocv.MorphologyEx(src, dst, gocv.MorphClose, kernel)
	for i := 1; i < m.Iterations; i++ {
		temp := gocv.NewMat()
		defer temp.Close()
		gocv.MorphologyEx(*dst, &temp, gocv.MorphClose, kernel)
		temp.CopyTo(dst)
	}
}

func init() {
	processor.Register("MorphClose", &MorphClose{
		KernelSize: 5,
		Iterations: 1,
	})
}
