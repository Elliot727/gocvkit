// Package core provides Otsu thresholding filter implementation.
//
// The Otsu filter applies automatic thresholding using Otsu's method to determine
// the optimal threshold value.
package core

import (
	"github.com/Elliot727/gocvkit/processor"
	"gocv.io/x/gocv"
)

// Otsu defines the configuration for Otsu thresholding.
type Otsu struct {
	MaxValue float32 `toml:"max_value"` // MaxValue is the maximum value to use with the threshold
	Invert   bool    `toml:"invert"`    // Invert indicates whether to invert the threshold result
}

// Process applies Otsu thresholding using the configured parameters.
func (o *Otsu) Process(src gocv.Mat, dst *gocv.Mat) {
	gocv.Threshold(src, dst, 0, o.MaxValue, gocv.ThresholdBinary|gocv.ThresholdOtsu)
}

func init() {
	processor.Register("Otsu", &Otsu{
		MaxValue: 255,
		Invert:   false,
	})
}
