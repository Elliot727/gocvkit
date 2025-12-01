// Package camera provides a clean, unified wrapper around gocv.VideoCapture
// that works identically for webcam devices and video files.
//
// Usage from config:
//
//	device_id = 0          → opens default webcam
//	file = "video.mp4"     → opens video file (device_id is ignored)
//
// The wrapper hides the difference between the two sources and adds
// convenient helpers (Width, Height, FPS).
package camera

import "gocv.io/x/gocv"

// Camera represents an open video source (webcam or file).
type Camera struct {
	device int                 // device is the camera device ID when using webcam
	file   string              // file is the path to video file when using file input
	cap    *gocv.VideoCapture // cap is the underlying video capture instance
}

// NewCamera opens either a webcam (by device ID) or a video file.
// If file is non-empty, it takes precedence over device.
func NewCamera(device int, file string) (*Camera, error) {
	var cap *gocv.VideoCapture
	var err error

	if file != "" {
		cap, err = gocv.VideoCaptureFile(file)
	} else {
		cap, err = gocv.OpenVideoCapture(device)
	}

	if err != nil || cap == nil {
		return nil, err
	}

	return &Camera{
		device: device,
		file:   file,
		cap:    cap,
	}, nil
}

// Read reads the next frame into the provided Mat.
// Returns false if no more frames are available (e.g. end of file or camera disconnected).
func (c *Camera) Read(frame *gocv.Mat) bool {
	return c.cap.Read(frame)
}

// Close releases the underlying VideoCapture.
func (c *Camera) Close() {
	c.cap.Close()
}

// Width returns the frame width of the video source.
func (c *Camera) Width() int {
	return int(c.cap.Get(gocv.VideoCaptureFrameWidth))
}

// Height returns the frame height of the video source.
func (c *Camera) Height() int {
	return int(c.cap.Get(gocv.VideoCaptureFrameHeight))
}

// FPS returns the frames per second of the video source (may be 0.0 for some webcams).
func (c *Camera) FPS() float64 {
	return c.cap.Get(gocv.VideoCaptureFPS)
}
