// Package processor provides Sobel edge detection filter implementation.
//
// The Sobel filter performs edge detection using the Sobel operator with a configurable kernel size.
package processor

import (
	"gocv.io/x/gocv"
)

// Sobel performs Sobel edge detection with configurable kernel size.
type Sobel struct {
	K int `toml:"sobel_size"` // K is the kernel size for Sobel edge detection (will be made odd if even)
}

// Process applies the Sobel operator and normalizes the output to 8-bit.
func (s *Sobel) Process(src gocv.Mat, dst *gocv.Mat) {
	k := s.K

	// Enforce valid Sobel kernel: positive odd integer
	if k < 1 {
		k = 1
	}
	if k%2 == 0 {
		k++
	}

	// Apply Sobel in both directions (dx=1, dy=1)
	// We use CV16S to avoid overflow, then convert back to 8-bit
	gocv.Sobel(src, dst, gocv.MatTypeCV16S, 1, 1, k, 1.0, 0, gocv.BorderDefault)

	// Convert to absolute values and scale to 8-bit for display
	gocv.ConvertScaleAbs(*dst, dst, 1, 0)
}

func init() {
	// Register directly with default kernel size of 3
	Register("Sobel", &Sobel{
		K: 3,
	})
}
