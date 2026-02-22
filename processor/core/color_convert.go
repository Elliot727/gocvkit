// Package core provides color space conversion filter implementation.
//
// The ColorConvert filter converts images between different color spaces using configurable codes.
package core

import (
	"fmt"

	"github.com/Elliot727/gocvkit/processor"
	"gocv.io/x/gocv"
)

// ColorConvert defines the configuration for color space conversion.
type ColorConvert struct {
	Code     string                   `toml:"code"` // Code specifies the color conversion code (e.g., "BGR2GRAY", "BGR2HSV", "HSV2BGR", etc.)
	codeEnum gocv.ColorConversionCode // Pre-calculated enum

}

// Validate checks the conversion code and pre-calculates the enum.
// Prevents runtime panics from typos in the TOML file.
func (c *ColorConvert) Validate() error {
	switch c.Code {
	case "BGR2GRAY":
		c.codeEnum = gocv.ColorBGRToGray
	case "BGR2HSV":
		c.codeEnum = gocv.ColorBGRToHSV
	case "HSV2BGR":
		c.codeEnum = gocv.ColorHSVToBGR
	case "BGR2LAB":
		c.codeEnum = gocv.ColorBGRToLab
	case "LAB2BGR":
		c.codeEnum = gocv.ColorLabToBGR
	case "BGR2YUV":
		c.codeEnum = gocv.ColorBGRToYUV
	case "YUV2BGR":
		c.codeEnum = gocv.ColorYUVToBGR
	default:
		return fmt.Errorf("unsupported color conversion code %q", c.Code)
	}
	return nil
}

// Process converts the image using the pre-calculated enum.
func (c *ColorConvert) Process(src gocv.Mat, dst *gocv.Mat) error {
	if src.Empty() {
		return nil
	}
	gocv.CvtColor(src, dst, c.codeEnum)
	return nil
}

func (c *ColorConvert) Close() {}

func init() {
	processor.Register("ColorConvert", &ColorConvert{
		Code: "BGR2GRAY",
	})
}
