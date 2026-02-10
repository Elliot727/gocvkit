// Package core provides color space conversion filter implementation.
//
// The ColorConvert filter converts images between different color spaces using configurable codes.
package core

import (
	"github.com/Elliot727/gocvkit/processor"
	"gocv.io/x/gocv"
)

// ColorConvert defines the configuration for color space conversion.
type ColorConvert struct {
	Code string `toml:"code"` // Code specifies the color conversion code (e.g., "BGR2GRAY", "BGR2HSV", "HSV2BGR", etc.)
}

// Process converts the image from one color space to another using the configured code.
func (c *ColorConvert) Process(src gocv.Mat, dst *gocv.Mat) {
	gocv.CvtColor(src, dst, colorCode(c.Code))
}

func colorCode(code string) gocv.ColorConversionCode {
	switch code {
	case "BGR2GRAY":
		return gocv.ColorBGRToGray
	case "BGR2HSV":
		return gocv.ColorBGRToHSV
	case "HSV2BGR":
		return gocv.ColorHSVToBGR
	case "BGR2LAB":
		return gocv.ColorBGRToLab
	case "LAB2BGR":
		return gocv.ColorLabToBGR
	case "BGR2YUV":
		return gocv.ColorBGRToYUV
	case "YUV2BGR":
		return gocv.ColorYUVToBGR
	default:
		panic("unsupported color conversion: " + code)
	}
}

func init() {
	processor.Register("ColorConvert", &ColorConvert{
		Code: "BGR2GRAY",
	})
}
