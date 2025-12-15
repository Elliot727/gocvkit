// Package streamer provides an HTTP MJPEG streaming server.
//
// MJPEGStreamer enables real-time streaming of video frames over HTTP using the MJPEG format.
// It implements the http.Handler interface to serve streams to multiple clients simultaneously.
// The streamer handles concurrent client connections, frame broadcasting, and rate limiting
// to maintain optimal performance.
package streamer

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"gocv.io/x/gocv"
)

// MJPEGStreamer represents an HTTP-based MJPEG streaming server.
// It manages multiple client connections and broadcasts frames to all connected clients.
type MJPEGStreamer struct {
	mu          sync.Mutex            // mu provides thread-safe access to clients and latestFrame
	clients     map[chan []byte]struct{} // clients stores active client channels for frame delivery
	latestFrame []byte                // latestFrame keeps the most recently encoded frame
	lastSent    time.Time             // lastSent tracks the time of the last frame broadcast
	interval    time.Duration         // interval sets the minimum time between consecutive broadcasts
}

// NewMJPEGStreamer creates and initializes a new MJPEG streamer instance.
// The default interval is set to ~15 FPS (time.Second / 15) to balance quality and performance.
func NewMJPEGStreamer() *MJPEGStreamer {
	return &MJPEGStreamer{
		clients:  make(map[chan []byte]struct{}),
		interval: time.Second / 15,
	}
}

// ServeHTTP handles incoming HTTP requests and establishes a streaming connection.
// It implements the http.Handler interface, allowing the streamer to be registered
// as an HTTP endpoint. Each client receives a continuous stream of JPEG frames
// using the multipart/x-mixed-replace protocol.
func (s *MJPEGStreamer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Set headers for MJPEG streaming
	w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary=frame")
	w.Header().Set("Cache-Control", "no-cache")

	// Verify that the ResponseWriter supports flushing
	flusher, ok := w.(http.Flusher)
	if !ok {
		return
	}

	// Create a channel for this client to receive frames
	clientChan := make(chan []byte, 1)

	// Add client to the list of active clients
	s.mu.Lock()
	s.clients[clientChan] = struct{}{}
	latest := s.latestFrame
	s.mu.Unlock()

	// Cleanup: remove client when connection closes
	defer func() {
		s.mu.Lock()
		delete(s.clients, clientChan)
		s.mu.Unlock()
	}()

	// writeFrame is a helper function to send a JPEG frame to the client
	writeFrame := func(b []byte) bool {
		if _, err := fmt.Fprintf(w, "--frame\r\nContent-Type: image/jpeg\r\nContent-Length: %d\r\n\r\n", len(b)); err != nil {
			return false
		}
		if _, err := w.Write(b); err != nil {
			return false
		}
		if _, err := w.Write([]byte("\r\n")); err != nil {
			return false
		}
		flusher.Flush()
		return true
	}

	// Send the latest frame immediately if available
	if latest != nil {
		if !writeFrame(latest) {
			return
		}
	}

	// Stream frames as they arrive
	for {
		select {
		case <-r.Context().Done():
			// Client disconnected or request cancelled
			return
		case frame := <-clientChan:
			// Send new frame to client
			if !writeFrame(frame) {
				return
			}
		}
	}
}

// Broadcast encodes a frame to JPEG and sends it to all connected clients.
// It applies rate limiting to prevent overwhelming clients and network.
// The quality parameter controls JPEG compression (0-100, higher is better quality).
func (s *MJPEGStreamer) Broadcast(frame gocv.Mat, quality int) {
	// Apply rate limiting to prevent sending too many frames
	if time.Since(s.lastSent) < s.interval {
		return
	}

	// Encode the frame to JPEG with the specified quality
	buf, _ := gocv.IMEncodeWithParams(".jpg", frame, []int{gocv.IMWriteJpegQuality, quality})
	jpegBytes := buf.GetBytes()
	buf.Close()

	// Update shared state and broadcast to all clients
	s.mu.Lock()
	s.lastSent = time.Now()
	s.latestFrame = jpegBytes
	for client := range s.clients {
		select {
		case client <- jpegBytes:
		// Skip slow clients to prevent blocking others
		default:
		}
	}
	s.mu.Unlock()
}
