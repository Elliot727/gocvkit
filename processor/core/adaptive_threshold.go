// Package core provides adaptive thresholding filter implementation.
//
// The Adaptive filter applies adaptive thresholding using the Gaussian method with configurable parameters.
package core

import (
	"github.com/Elliot727/gocvkit/processor"
	"gocv.io/x/gocv"
)

// Adaptive defines the configuration for adaptive thresholding.
type Adaptive struct {
	MaxValue  float32 `toml:"max_value"`  // MaxValue is the maximum value to use with the threshold
	BlockSize int     `toml:"block_size"` // BlockSize is the size of the pixel neighborhood for adaptive thresholding
	C         float32 `toml:"c"`          // C is the constant subtracted from the mean or weighted mean
}

// Process applies adaptive thresholding using the configured parameters.
func (a *Adaptive) Process(src gocv.Mat, dst *gocv.Mat) {
	gocv.AdaptiveThreshold(
		src,
		dst,
		a.MaxValue,
		gocv.AdaptiveThresholdGaussian,
		gocv.ThresholdBinary,
		a.BlockSize,
		a.C,
	)
}

func init() {
	processor.Register("Adaptive", &Adaptive{
		MaxValue:  255,
		BlockSize: 11,
		C:         2,
	})
}
