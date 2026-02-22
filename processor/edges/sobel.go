// Package edges provides Sobel edge detection filter implementation.
//
// The Sobel filter performs edge detection using the Sobel operator with a configurable kernel size.
package edges

import (
	"fmt"

	"github.com/Elliot727/gocvkit/processor"
	"gocv.io/x/gocv"
)

// Sobel performs Sobel edge detection with configurable kernel size.
type Sobel struct {
	K int `toml:"sobel_size"` // K is the kernel size for Sobel edge detection (will be made odd if even)
}

func (s *Sobel) Validate() error {
	if s.K <= 0 {
		return fmt.Errorf("kernel size must be > 0, got %d", s.K)
	}
	if s.K%2 == 0 {
		return fmt.Errorf("kernel size must be odd, got %d", s.K)
	}
	return nil
}

// Process applies the Sobel operator and normalizes the output to 8-bit.
func (s *Sobel) Process(src gocv.Mat, dst *gocv.Mat) error {
	if src.Empty() {
		return nil
	}

	gocv.Sobel(src, dst, gocv.MatTypeCV16S, 1, 1, s.K, 1.0, 0, gocv.BorderDefault)

	gocv.ConvertScaleAbs(*dst, dst, 1, 0)

	return nil
}

func (s *Sobel) Close() {}

func init() {
	processor.Register("Sobel", &Sobel{
		K: 3,
	})
}
