// Package core provides image flipping filter implementation.
//
// The Flip filter flips images horizontally, vertically, or both ways depending on the mode.
package core

import (
	"fmt"

	"github.com/Elliot727/gocvkit/processor"
	"gocv.io/x/gocv"
)

// Flip defines the configuration for image flipping.
type Flip struct {
	Mode     string `toml:"mode"` // Mode specifies the flip direction: "horizontal", "vertical", or "both"
	modeCode int
}

func (f *Flip) Validate() error {
	switch f.Mode {
	case "horizontal":
		f.modeCode = 1
	case "vertical":
		f.modeCode = 0
	case "both":
		f.modeCode = -1
	default:
		return fmt.Errorf("invalid flip mode %q (use 'horizontal', 'vertical', or 'both')", f.Mode)
	}
	return nil
}

// Process flips the image using the configured mode.
func (f *Flip) Process(src gocv.Mat, dst *gocv.Mat) error {
	if src.Empty() {
		return nil
	}
	gocv.Flip(src, dst, f.modeCode)
	return nil
}

func (f *Flip) Close() {}

func init() {
	processor.Register("Flip", &Flip{
		Mode: "horizontal",
	})
}
