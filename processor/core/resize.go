// Package core provides image resizing filter implementation.
//
// The Resize filter resizes images to the specified width and height dimensions.
package core

import (
	"image"

	"github.com/Elliot727/gocvkit/processor"
	"gocv.io/x/gocv"
)

// Resize defines the configuration for image resizing.
type Resize struct {
	Width  int `toml:"width"`  // Width is the target width for the resized image
	Height int `toml:"height"` // Height is the target height for the resized image
}

// Process resizes the image to the configured dimensions using area interpolation.
func (r *Resize) Process(src gocv.Mat, dst *gocv.Mat) {
	gocv.Resize(src, dst, image.Pt(r.Width, r.Height), 0, 0, gocv.InterpolationArea)
}

func init() {
	processor.Register("Resize", &Resize{
		Width:  224,
		Height: 224,
	})
}
