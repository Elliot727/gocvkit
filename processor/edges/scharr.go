// Package edges provides Scharr edge detection filter implementation.
//
// The Scharr filter performs edge detection using the Scharr operator, which is more
// accurate than the Sobel operator for certain applications.
package edges

import (
	"github.com/Elliot727/gocvkit/processor"
	"gocv.io/x/gocv"
)

// Scharr defines the configuration for Scharr edge detection filter.
type Scharr struct{}

func (s *Scharr) Validate() error {
	return nil
}

// Process applies Scharr edge detection.
func (s *Scharr) Process(src gocv.Mat, dst *gocv.Mat) error {
	if src.Empty() {
		return nil
	}

	gocv.Scharr(src, dst, gocv.MatTypeCV16S, 1, 0, 1, 0, gocv.BorderDefault)
	return nil
}

func (s *Scharr) Close() {}

func init() {
	processor.Register("Scharr", &Scharr{})
}
