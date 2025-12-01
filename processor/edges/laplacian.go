// Package edges provides Laplacian edge detection filter implementation.
//
// The Laplacian filter performs edge detection using the Laplacian operator, which
// calculates the second derivative of the image.
package edges

import (
	"github.com/Elliot727/gocvkit/processor"
	"gocv.io/x/gocv"
)

// Laplacian defines the configuration for Laplacian edge detection filter.
type Laplacian struct {
	K int `toml:"k"` // K is the aperture size for the Laplacian operator
}

// Process applies Laplacian edge detection using the configured parameter.
func (l *Laplacian) Process(src gocv.Mat, dst *gocv.Mat) {
	k := l.K
	if k < 1 {
		k = 1
	}
	gocv.Laplacian(src, dst, gocv.MatTypeCV16S, k, 1, 0, gocv.BorderDefault)
}

func init() {
	processor.Register("Laplacian", &Laplacian{
		K: 3,
	})
}
