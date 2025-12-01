// Package display provides a minimal, clean wrapper around gocv.Window.
//
// It exposes only the three operations needed by the main app:
//   - showing frames
//   - reading keyboard input (for quit detection)
//   - proper cleanup
//
// This keeps the rest of the codebase decoupled from gocv.Window details.
package display

import "gocv.io/x/gocv"

// Display represents an on-screen window for showing processed frames.
type Display struct {
	window *gocv.Window
}

// New creates a new named window.
func New(title string) *Display {
	return &Display{window: gocv.NewWindow(title)}
}

// Show displays the given Mat in the window.
func (d *Display) Show(img gocv.Mat) {
	d.window.IMShow(img)
}

// Key waits up to delay milliseconds for a key press and returns the key code.
// Commonly used with delay = 1 to keep the window responsive.
func (d *Display) Key(delay int) int {
	return d.window.WaitKey(delay)
}

// Close destroys the window and releases native resources.
func (d *Display) Close() {
	d.window.Close()
}
