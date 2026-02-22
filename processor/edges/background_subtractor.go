package edges

import (
	"fmt"

	"github.com/Elliot727/gocvkit/processor"
	"gocv.io/x/gocv"
)

type BackgroundSubtractor struct {
	Algorithm    string  `toml:"algorithm"`
	LearningRate float64 `toml:"learning_rate"`

	mog2 *gocv.BackgroundSubtractorMOG2
	knn  *gocv.BackgroundSubtractorKNN
}

func (b *BackgroundSubtractor) Validate() error {
	if b.mog2 != nil || b.knn != nil {
		return nil
	}

	switch b.Algorithm {
	case "MOG2":
		sub := gocv.NewBackgroundSubtractorMOG2()
		b.mog2 = &sub
	case "KNN":
		sub := gocv.NewBackgroundSubtractorKNN()
		b.knn = &sub
	default:
		return fmt.Errorf("unsupported background subtractor algorithm %q (use 'MOG2' or 'KNN')", b.Algorithm)
	}
	return nil
}

// Process applies background subtraction using the pre-initialized backend.
func (b *BackgroundSubtractor) Process(src gocv.Mat, dst *gocv.Mat) error {
	if src.Empty() {
		return nil
	}

	if b.mog2 == nil && b.knn == nil {
		return fmt.Errorf("background subtractor not initialized")
	}

	if b.mog2 != nil {

		b.mog2.ApplyWithLearningRate(src, dst, b.LearningRate)
		return nil
	}

	b.knn.Apply(src, dst)
	return nil
}

func (b *BackgroundSubtractor) Close() {
	if b.mog2 != nil {
		b.mog2.Close()
		b.mog2 = nil
	}
	if b.knn != nil {
		b.knn.Close()
		b.knn = nil
	}
}

func init() {
	processor.Register("BackgroundSubtractor", &BackgroundSubtractor{
		Algorithm:    "MOG2",
		LearningRate: 0.01,
	})
}
