// Package core provides morphological close operation filter implementation.
//
// The morphological close operation is useful for closing small holes and gaps in objects.
package core

import (
	"fmt"
	"image"

	"github.com/Elliot727/gocvkit/processor"
	"gocv.io/x/gocv"
)

// MorphClose defines the configuration for morphological close operation.
type MorphClose struct {
	KernelSize int `toml:"kernel"`     // KernelSize is the size of the structuring element for morphological close
	Iterations int `toml:"iterations"` // Iterations is the number of times morphological close is applied
	kernel     gocv.Mat
	temp       gocv.Mat
}

func (m *MorphClose) Validate() error {
	if m.KernelSize < 1 {
		return fmt.Errorf("kernel size must be >= 1, got %d", m.KernelSize)
	}
	if m.KernelSize%2 == 0 {
		return fmt.Errorf("kernel size must be odd, got %d", m.KernelSize)
	}
	if m.Iterations < 1 {
		return fmt.Errorf("iterations must be >= 1, got %d", m.Iterations)
	}

	// Pre-create kernel once
	m.kernel = gocv.GetStructuringElement(gocv.MorphRect, image.Pt(m.KernelSize, m.KernelSize))
	if m.kernel.Empty() {
		return fmt.Errorf("failed to create structuring element")
	}
	return nil
}

// Process applies morphological close operation using the configured parameters.
func (m *MorphClose) Process(src gocv.Mat, dst *gocv.Mat) error {
	if src.Empty() {
		return nil
	}

	// First pass: src -> dst
	gocv.MorphologyEx(src, dst, gocv.MorphClose, m.kernel)

	// Subsequent passes: dst -> dst (via temp)
	if m.Iterations > 1 {
		// Lazy init / Resize check for temp buffer
		if m.temp.Empty() || m.temp.Rows() != src.Rows() || m.temp.Cols() != src.Cols() {
			if !m.temp.Empty() {
				m.temp.Close()
			}
			m.temp = gocv.NewMatWithSize(src.Rows(), src.Cols(), src.Type())
		}

		for i := 1; i < m.Iterations; i++ {
			gocv.MorphologyEx(*dst, &m.temp, gocv.MorphClose, m.kernel)
			m.temp.CopyTo(dst)
		}
	}
	return nil
}

func (m *MorphClose) Close() {
	if !m.kernel.Empty() {
		m.kernel.Close()
	}
	if !m.temp.Empty() {
		m.temp.Close()
	}
}

func init() {
	processor.Register("MorphClose", &MorphClose{
		KernelSize: 5,
		Iterations: 1,
	})
}
