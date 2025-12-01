// Package core converts a color image to 8-bit single-channel grayscale.
//
// This is typically used as the first step before edge detection, thresholding,
// or any algorithm that expects a single-channel image.
//
// No configuration options — simply add to the pipeline:
//
//	[[pipeline.steps]]
//	name = "Grayscale"
package core

import (
	"github.com/Elliot727/gocvkit/processor"
	"gocv.io/x/gocv"
)

// Grayscale converts color images to grayscale.
type Grayscale struct{}

// Process converts src from BGR to grayscale and writes the result to dst.
func (Grayscale) Process(src gocv.Mat, dst *gocv.Mat) {
	gocv.CvtColor(src, dst, gocv.ColorBGRToGray)
}

func init() {
	processor.Register("Grayscale", &Grayscale{})
}
