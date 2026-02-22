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
	Angle       float64 `toml:"angle"` // Angle is the rotation angle in degrees (90, 180, 270 for optimized rotations)
	isOptimized bool
	optCode     gocv.RotateFlag
	hasMatrix   bool
	mat         gocv.Mat
}

func (r *Rotate) Validate() error {
	angle := r.Angle
	for angle < 0 {
		angle += 360
	}
	for angle >= 360 {
		angle -= 360
	}

	switch angle {
	case 90:
		r.isOptimized = true
		r.optCode = gocv.Rotate90Clockwise
	case 180:
		r.isOptimized = true
		r.optCode = gocv.Rotate180Clockwise
	case 270:
		r.isOptimized = true
		r.optCode = gocv.Rotate90CounterClockwise
	default:
		r.isOptimized = false
		r.hasMatrix = false
	}

	return nil
}

// Process rotates the image by the configured angle.
func (r *Rotate) Process(src gocv.Mat, dst *gocv.Mat) error {
	if src.Empty() {
		return nil
	}

	if r.isOptimized {
		gocv.Rotate(src, dst, r.optCode)
		return nil
	}

	if !r.hasMatrix || r.mat.Rows() != src.Rows() || r.mat.Cols() != src.Cols() {
		if !r.mat.Empty() {
			r.mat.Close()
		}
		center := image.Pt(src.Cols()/2, src.Rows()/2)
		r.mat = gocv.GetRotationMatrix2D(center, r.Angle, 1.0)
		r.hasMatrix = true
	}

	gocv.WarpAffine(src, dst, r.mat, image.Pt(src.Cols(), src.Rows()))
	return nil
}

func (r *Rotate) Close() {
	if r.hasMatrix {
		if !r.mat.Empty() {
			r.mat.Close()
		}
		r.mat = gocv.NewMat()
		r.hasMatrix = false
	}
}

func init() {
	processor.Register("Rotate", &Rotate{
		Angle: 90,
	})
}
