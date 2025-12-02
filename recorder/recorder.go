// Package recorder provides video recording capabilities for the GoCVKit framework.
//
// It handles video file creation, manages format changes during recording,
// and provides automatic file rotation when pipeline parameters change
// (e.g. when switching from grayscale to color or changing image dimensions).
package recorder

import (
	"fmt"
	"path/filepath"
	"strings"

	"gocv.io/x/gocv"
)

// Recorder manages video recording with automatic file rotation when format changes.
type Recorder struct {
	writer *gocv.VideoWriter
	fps    float64
	fourcc string

	// File naming
	baseName string
	ext      string
	counter  int

	// Current format state
	width    int
	height   int
	channels int
}

// NewRecorder creates a new Recorder that writes video files to the specified path.
// If no extension is provided, it defaults to .mp4. The recorder automatically
// handles file rotation when the input format changes during pipeline updates.
func NewRecorder(path string) *Recorder {
	// Split "output.mp4" into "output" and ".mp4"
	// so we can insert numbers later: "output-1.mp4"
	ext := filepath.Ext(path)
	base := strings.TrimSuffix(path, ext)
	if ext == "" {
		ext = ".mp4"
	}

	return &Recorder{
		baseName: base,
		ext:      ext,
		fps:      30.0,
		fourcc:   "mp4v",
	}
}

// SetFPS sets the frame rate for the output video file.
// Only positive values are accepted; negative or zero values are ignored.
func (r *Recorder) SetFPS(fps float64) {
	if fps > 0 {
		r.fps = fps
	}
}

// Write adds the given frame to the video file.
// The recorder automatically handles format changes by creating new files
// when the input dimensions or channel count changes.
func (r *Recorder) Write(frame gocv.Mat) error {
	if frame.Empty() {
		return nil
	}

	currentCols := frame.Cols()
	currentRows := frame.Rows()
	currentCh := frame.Channels()

	// CHECK: Did the format change since the last frame?
	// If dimensions or channels changed, we MUST start a new file.
	if r.writer != nil {
		if currentCols != r.width || currentRows != r.height || currentCh != r.channels {
			fmt.Printf("🔄 Pipeline changed (%dx%d %dc -> %dx%d %dc). Rotating video file...\n",
				r.width, r.height, r.channels, currentCols, currentRows, currentCh)
			r.Close() // Close the old file
		}
	}

	// INITIALIZE: Open a new writer if needed
	if r.writer == nil {
		r.width = currentCols
		r.height = currentRows
		r.channels = currentCh

		isColor := true
		if r.channels == 1 {
			isColor = false
		}

		// Create filename: "output-0.mp4", "output-1.mp4", etc.
		filename := fmt.Sprintf("%s-%d%s", r.baseName, r.counter, r.ext)
		r.counter++

		w, err := gocv.VideoWriterFile(filename, r.fourcc, r.fps, r.width, r.height, isColor)
		if err != nil {
			return fmt.Errorf("failed to open recorder: %w", err)
		}
		r.writer = w
	}

	return r.writer.Write(frame)
}

// Close releases all resources used by the recorder and finalizes the video file.
// Safe to call multiple times.
func (r *Recorder) Close() {
	if r.writer != nil {
		r.writer.Close()
		r.writer = nil
	}
}
