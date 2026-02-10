package edges

import (
	"github.com/Elliot727/gocvkit/processor"
	"gocv.io/x/gocv"
)

type BackgroundSubtractor struct {
	Algorithm    string  `toml:"algorithm"`
	LearningRate float64 `toml:"learning_rate"`

	mog2 *gocv.BackgroundSubtractorMOG2
	knn  *gocv.BackgroundSubtractorKNN
}

func (b *BackgroundSubtractor) Process(src gocv.Mat, dst *gocv.Mat) {
	if b.mog2 == nil && b.knn == nil {
		switch b.Algorithm {
		case "MOG2":
			sub := gocv.NewBackgroundSubtractorMOG2()
			b.mog2 = &sub
		case "KNN":
			sub := gocv.NewBackgroundSubtractorKNN()
			b.knn = &sub
		default:
			panic("unsupported background subtractor: " + b.Algorithm)
		}
	}

	if b.mog2 != nil {
		b.mog2.ApplyWithLearningRate(src, dst, b.LearningRate)
		return
	}

	b.knn.Apply(src, dst)
}

func init() {
	processor.Register("BackgroundSubtractor", &BackgroundSubtractor{
		Algorithm:    "MOG2",
		LearningRate: 0.01,
	})
}
