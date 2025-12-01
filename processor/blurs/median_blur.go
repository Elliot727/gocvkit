// Package blurs provides median blur filter implementation.
//
// The median blur filter applies median filtering with a configurable kernel size.
package blurs

import (
	"github.com/Elliot727/gocvkit/processor"
	"gocv.io/x/gocv"
)

// MedianBlur performs median filtering with a configurable kernel size.
type MedianBlur struct {
	K int `toml:"k"` // K is the kernel size for median blur (will be made odd if even)
}

// Process applies median blur using the configured kernel size.
func (m *MedianBlur) Process(src gocv.Mat, dst *gocv.Mat) {
	k := m.K

	// Enforce valid median blur kernel: positive odd integer
	if k < 1 {
		k = 1
	}
	if k%2 == 0 {
		k++
	}

	gocv.MedianBlur(src, dst, k)
}

func init() {
	// Register directly with a sensible default
	processor.Register("MedianBlur", &MedianBlur{
		K: 5,
	})
}
