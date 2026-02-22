// Package core provides adaptive thresholding filter implementation.
//
// The Adaptive filter applies adaptive thresholding using the Gaussian method with configurable parameters.
package core

import (
	"fmt"

	"github.com/Elliot727/gocvkit/processor"
	"gocv.io/x/gocv"
)

// Adaptive defines the configuration for adaptive thresholding.
type Adaptive struct {
	MaxValue  float32 `toml:"max_value"`  // MaxValue is the maximum value to use with the threshold
	BlockSize int     `toml:"block_size"` // BlockSize is the size of the pixel neighborhood for adaptive thresholding
	C         float32 `toml:"c"`          // C is the constant subtracted from the mean or weighted mean
}

// Validate checks constraints before the pipeline starts.
func (a *Adaptive) Validate() error {
	if a.MaxValue <= 0 {
		return fmt.Errorf("max_value must be > 0, got %f", a.MaxValue)
	}
	if a.BlockSize <= 1 {
		return fmt.Errorf("block_size must be > 1, got %d", a.BlockSize)
	}
	if a.BlockSize%2 == 0 {
		return fmt.Errorf("block_size must be odd, got %d", a.BlockSize)
	}
	// C can be negative (subtracting from mean), so no check needed there.
	return nil
}

// Process applies adaptive thresholding using the configured parameters.
func (a *Adaptive) Process(src gocv.Mat, dst *gocv.Mat) error {
	if src.Empty() {
		return nil
	}
	gocv.AdaptiveThreshold(
		src,
		dst,
		a.MaxValue,
		gocv.AdaptiveThresholdGaussian,
		gocv.ThresholdBinary,
		a.BlockSize,
		a.C,
	)
	return nil
}

func (a *Adaptive) Clos() {}

func init() {
	processor.Register("Adaptive", &Adaptive{
		MaxValue:  255,
		BlockSize: 11,
		C:         2,
	})
}
