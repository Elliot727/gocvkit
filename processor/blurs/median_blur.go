// Package blurs provides median blur filter implementation.
//
// The median blur filter applies median filtering with a configurable kernel size.
package blurs

import (
	"fmt"

	"github.com/Elliot727/gocvkit/processor"
	"gocv.io/x/gocv"
)

// MedianBlur performs median filtering with a configurable kernel size.
type MedianBlur struct {
	K int `toml:"k"` // K is the kernel size for median blur (will be made odd if even)
}

// Validate checks constraints before the pipeline starts.
func (m *MedianBlur) Validate() error {
	if m.K < 1 {
		return fmt.Errorf("kernel size must be >= 1, got %d", m.K)
	}
	if m.K%2 == 0 {
		return fmt.Errorf("kernel size must be odd, got %d", m.K)
	}
	return nil
}

// Process applies median blur using the configured kernel size.
func (m *MedianBlur) Process(src gocv.Mat, dst *gocv.Mat) error {
	if src.Empty() {
		return nil
	}
	gocv.MedianBlur(src, dst, m.K)
	return nil
}

func (m *MedianBlur) Close() {}

func init() {
	// Register directly with a sensible default
	processor.Register("MedianBlur", &MedianBlur{
		K: 5,
	})
}
