// Package core provides Otsu thresholding filter implementation.
//
// The Otsu filter applies automatic thresholding using Otsu's method to determine
// the optimal threshold value.
package core

import (
	"fmt"

	"github.com/Elliot727/gocvkit/processor"
	"gocv.io/x/gocv"
)

// Otsu defines the configuration for Otsu thresholding.
type Otsu struct {
	MaxValue float32 `toml:"max_value"` // MaxValue is the maximum value to use with the threshold
	Invert   bool    `toml:"invert"`    // Invert indicates whether to invert the threshold result
	flags    gocv.ThresholdType
}

func (o *Otsu) Validate() error {
	if o.MaxValue <= 0 {
		return fmt.Errorf("max_value must be > 0, got %f", o.MaxValue)
	}

	// Pre-calculate the combined flag
	o.flags = gocv.ThresholdBinary | gocv.ThresholdOtsu
	if o.Invert {
		o.flags |= gocv.ThresholdBinaryInv
	}

	return nil
}

// Process applies Otsu thresholding using the configured parameters.
func (o *Otsu) Process(src gocv.Mat, dst *gocv.Mat) error {
	if src.Empty() {
		return nil
	}

	gocv.Threshold(src, dst, 0, o.MaxValue, o.flags)
	return nil
}

func (o *Otsu) Close() {}

func init() {
	processor.Register("Otsu", &Otsu{
		MaxValue: 255,
		Invert:   false,
	})
}
