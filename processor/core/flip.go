// Package core provides image flipping filter implementation.
//
// The Flip filter flips images horizontally, vertically, or both ways depending on the mode.
package core

import (
	"github.com/Elliot727/gocvkit/processor"
	"gocv.io/x/gocv"
)

// Flip defines the configuration for image flipping.
type Flip struct {
	Mode string `toml:"mode"` // Mode specifies the flip direction: "horizontal", "vertical", or "both"
}

// Process flips the image using the configured mode.
func (f *Flip) Process(src gocv.Mat, dst *gocv.Mat) {
	var code int
	switch f.Mode {
	case "horizontal":
		code = 1
	case "vertical":
		code = 0
	case "both":
		code = -1
	default:
		panic("unsupported flip mode: " + f.Mode)
	}

	gocv.Flip(src, dst, code)
}

func init() {
	processor.Register("Flip", &Flip{
		Mode: "horizontal",
	})
}
