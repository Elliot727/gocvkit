package main

import (
	"log"

	"github.com/Elliot727/gocvkit"
	"gocv.io/x/gocv"
)

type CannyEdge struct {
	Low     float64 `toml:"low"`
	High    float64 `toml:"high"`
	Enabled bool    `toml:"enabled"`
}

func (c *CannyEdge) Process(src gocv.Mat, dst *gocv.Mat) {
	if !c.Enabled {
		src.CopyTo(dst)
		return
	}

	// Apply Canny edge detection using the gocvkit.Canny processor
	internalCanny := &gocvkit.Canny{
		Low:  c.Low,
		High: c.High,
	}
	internalCanny.Process(src, dst)
}

func init() {
	// Register the CannyEdge processor
	gocvkit.RegisterProcessor("CannyEdge", &CannyEdge{})
}

func main() {
	app, err := gocvkit.NewApp("config.toml")
	if err != nil {
		log.Fatal("Failed to create app:", err)
	}
	defer app.Close()

	if err := app.Run(nil); err != nil {
		log.Fatal("Run failed:", err)
	}
}
