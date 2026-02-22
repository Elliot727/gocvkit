// Package edges provides Laplacian edge detection filter implementation.
//
// The Laplacian filter performs edge detection using the Laplacian operator, which
// calculates the second derivative of the image.
package edges

import (
	"fmt"

	"github.com/Elliot727/gocvkit/processor"
	"gocv.io/x/gocv"
)

// Laplacian defines the configuration for Laplacian edge detection filter.
type Laplacian struct {
	K int `toml:"k"` // K is the aperture size for the Laplacian operator
}

func (l *Laplacian) Validate() error {
	if l.K <= 0 {
		return fmt.Errorf("aperture size must be > 0, got %d", l.K)
	}
	if l.K%2 == 0 {
		return fmt.Errorf("aperture size must be odd, got %d", l.K)
	}
	return nil
}

// Process applies Laplacian edge detection using the configured parameter.
func (l *Laplacian) Process(src gocv.Mat, dst *gocv.Mat) error {
	if src.Empty() {
		return nil
	}
	gocv.Laplacian(src, dst, gocv.MatTypeCV16S, l.K, 1, 0, gocv.BorderDefault)
	return nil
}

func (l *Laplacian) Close() {}

func init() {
	processor.Register("Laplacian", &Laplacian{
		K: 3,
	})
}
