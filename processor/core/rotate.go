// Package core provides image rotation filter implementation.
//
// The Rotate filter rotates images by the specified angle in degrees.
package core

import (
	"image"

	"github.com/Elliot727/gocvkit/processor"
	"gocv.io/x/gocv"
)

// Rotate defines the configuration for image rotation.
type Rotate struct {
	Angle float64 `toml:"angle"` // Angle is the rotation angle in degrees (90, 180, 270 for optimized rotations)
}

// Process rotates the image by the configured angle.
func (r *Rotate) Process(src gocv.Mat, dst *gocv.Mat) {
	switch r.Angle {
	case 90:
		gocv.Rotate(src, dst, gocv.Rotate90Clockwise)
		return
	case 180:
		gocv.Rotate(src, dst, gocv.Rotate180Clockwise)
		return
	case 270:
		gocv.Rotate(src, dst, gocv.Rotate90CounterClockwise)
		return
	}

	center := image.Pt(src.Cols()/2, src.Rows()/2)
	mat := gocv.GetRotationMatrix2D(center, r.Angle, 1.0)
	defer mat.Close()

	gocv.WarpAffine(src, dst, mat, image.Pt(src.Cols(), src.Rows()))
}

func init() {
	processor.Register("Rotate", &Rotate{
		Angle: 90,
	})
}
