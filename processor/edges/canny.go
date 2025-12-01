// Package edges provides Canny edge detection filter implementation.
//
// The Canny filter detects edges in an image using configurable low and high thresholds.
package edges

import (
	"github.com/Elliot727/gocvkit/processor"
	"gocv.io/x/gocv"
)

// Canny defines the configuration for Canny edge detection filter.
type Canny struct {
	Low  float64 `toml:"low"`  // Low is the lower threshold for edge detection
	High float64 `toml:"high"` // High is the upper threshold for edge detection
}

// Process applies Canny edge detection using the configured low and high thresholds.
func (c *Canny) Process(src gocv.Mat, dst *gocv.Mat) {
	gocv.Canny(src, dst, float32(c.Low), float32(c.High))
}

func init() {
	// Clean and obvious
	processor.Register("Canny", &Canny{
		Low:  50,
		High: 150,
	})
}
