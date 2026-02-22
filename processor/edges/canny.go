// Package edges provides Canny edge detection filter implementation.
//
// The Canny filter detects edges in an image using configurable low and high thresholds.
package edges

import (
	"fmt"

	"github.com/Elliot727/gocvkit/processor"
	"gocv.io/x/gocv"
)

// Canny defines the configuration for Canny edge detection filter.
type Canny struct {
	Low  float64 `toml:"low"`  // Low is the lower threshold for edge detection
	High float64 `toml:"high"` // High is the upper threshold for edge detection
}

func (c *Canny) Validate() error {
	if c.Low < 0 {
		return fmt.Errorf("low threshold must be >= 0, got %f", c.Low)
	}
	if c.High < 0 {
		return fmt.Errorf("high threshold must be >= 0, got %f", c.High)
	}
	if c.Low > c.High {
		return fmt.Errorf("low threshold (%f) cannot be greater than high threshold (%f)", c.Low, c.High)
	}
	return nil
}

// Process applies Canny edge detection using the configured low and high thresholds.
func (c *Canny) Process(src gocv.Mat, dst *gocv.Mat) error {
	if src.Empty() {
		return nil
	}
	gocv.Canny(src, dst, float32(c.Low), float32(c.High))
	return nil
}

func (c *Canny) Close() {}

func init() {
	processor.Register("Canny", &Canny{
		Low:  50,
		High: 150,
	})
}
