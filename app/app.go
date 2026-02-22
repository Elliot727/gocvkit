// Package app is the main orchestrator of the GoCVKit framework.
//
// It combines camera input, configurable processing pipeline, window display,
// and live config hot-reloading into a tiny, safe, and idiomatic public API:
//
//	app, _ := gocvkit.NewApp("config.toml")
//	defer app.Close()
//	app.Run(nil) // or pass func(*gocv.Mat) for per-frame processing
//
// All concurrency, resource management, graceful shutdown (Ctrl+C, Esc/q),
// and zero-leak pipeline swapping are handled automatically.
package app

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Elliot727/gocvkit/builder"
	"github.com/Elliot727/gocvkit/camera"
	"github.com/Elliot727/gocvkit/config"
	"github.com/Elliot727/gocvkit/display"
	"github.com/Elliot727/gocvkit/pipeline"
	"github.com/Elliot727/gocvkit/recorder"
	"github.com/Elliot727/gocvkit/streamer"
	"github.com/fsnotify/fsnotify"

	"gocv.io/x/gocv"
)

// App represents a fully configured and running computer vision application.
type App struct {
	mu         sync.RWMutex       // mu provides thread-safe access to mutable fields
	Camera     *camera.Camera     // Camera handles video input from webcam or file
	Recorder   *recorder.Recorder // Recorder manages video file output
	Streamer   *streamer.MJPEGStreamer
	Display    *display.Display   // Display shows processed frames in a window
	Pipeline   *pipeline.Pipeline // Pipeline processes frames through configured steps
	Config     *config.Config     // Config holds the current application configuration
	configPath string             // configPath is the path to the config file for hot-reloading
}

// New creates and returns a new App instance from the given TOML config file.
// Config changes are automatically detected and applied at runtime.
func New(cfgPath string) (*App, error) {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return nil, err
	}

	cam, err := camera.NewCamera(cfg.Camera.DeviceID, cfg.Camera.File)
	if err != nil {
		return nil, err
	}

	output := cfg.App.Output
	if output == "" {
		output = "gocvkit_capture.mp4"
	}

	rec := recorder.NewRecorder(output)

	win := display.New(cfg.App.WindowName)

	str := streamer.NewMJPEGStreamer()

	if cfg.Stream.Quality == 0 {
		cfg.Stream.Quality = 75
	}

	steps, err := builder.BuildPipeline(cfg)
	if err != nil {
		cam.Close()

		win.Close()
		return nil, err
	}

	a := &App{
		Camera:     cam,
		Recorder:   rec,
		Streamer:   str,
		Display:    win,
		Pipeline:   pipeline.New(steps),
		Config:     cfg,
		configPath: cfgPath,
	}

	if cfg.Stream.Enabled {
		mux := http.NewServeMux()
		mux.Handle(cfg.Stream.Path, str)

		addr := fmt.Sprintf(":%d", cfg.Stream.Port)
		go http.ListenAndServe(addr, mux)
	}
	go a.watchConfig() // fire-and-forget hot reload
	return a, nil
}

// Close releases all resources (camera, window, pipeline).
func (a *App) Close() {
	a.Camera.Close()
	a.Display.Close()

	a.Recorder.Close()
	a.mu.Lock()
	if a.Pipeline != nil {
		a.Pipeline.Close()
	}
	a.mu.Unlock()
}

// Run starts the capture -> process -> display loop.
// The function orchestrates the entire pipeline: reading frames from the camera,
// processing them through the configured pipeline steps, and displaying the results.
//
// The optional frameCallback is called with the final processed frame
// before it is displayed (useful for overlays, saving, logging, etc.).
// Pass nil if not needed.
//
// Run blocks until the user presses Esc/q or sends Ctrl+C.
// Returns any error that occurs during execution or context cancellation.
func (a *App) Run(frameCallback func(*gocv.Mat)) error {
	if frameCallback == nil {
		frameCallback = func(*gocv.Mat) {}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() { <-sig; cancel() }()

	frames := make(chan gocv.Mat, 10)
	results := make(chan gocv.Mat, 10)

	delay := 1
	recFPS := 30.0
	if a.Config.Camera.File != "" {
		fps := a.Camera.FPS()

		// Fallback for files with missing/bad metadata
		if fps <= 0 || fps > 200 {
			fps = 30.0
		}
		recFPS = fps
		// 1000ms / FPS = delay in ms (e.g., 30fps -> 33ms)
		delay = int(1000.0 / fps)
	}

	if a.Config.App.Record {
		if fps := a.Camera.FPS(); fps > 0 {
			recFPS = fps
		}
		a.Recorder.SetFPS(recFPS)
	}

	go func() {
		defer close(frames)
		for ctx.Err() == nil {
			img := gocv.NewMat()

			// Read frame
			if ok := a.Camera.Read(&img); !ok || img.Empty() {
				img.Close()
				return // End of file or error
			}

			select {
			case frames <- img:
				// Frame sent successfully
			case <-ctx.Done():
				img.Close()
				return
			}
		}
	}()

	go func() {
		defer close(results)
		for img := range frames {
			out := gocv.NewMat()

			a.mu.RLock()
			err := a.Pipeline.Run(img, &out)
			a.mu.RUnlock()

			img.Close() // We are done with the input frame

			if err != nil {
				out.Close()
				log.Printf("Pipeline error: %v", err)
				continue
			}

			select {
			case results <- out:
			case <-ctx.Done():
				out.Close()
				return
			}
		}
	}()

	showFPS := false

	// Setup "Bucket" variables for stable FPS calculation
	fpsTicker := time.Now() // The starting gun
	fpsCounter := 0         // The bucket of frames
	fpsText := "FPS: --"    // The text we actually draw (updated rarely)

	// Pre-allocate colors
	green := color.RGBA{0, 255, 0, 0}
	blackShadow := color.RGBA{0, 0, 0, 0}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case m, ok := <-results:
			if !ok {
				return nil
			}

			// 1. Run User Callback
			frameCallback(&m)

			// 2. Update FPS Logic
			fpsCounter++

			// Only recalculate the text string every 500ms (0.5 seconds)
			if time.Since(fpsTicker) >= 500*time.Millisecond {
				// Calculate average: Frames / Seconds Elapsed
				currentFPS := float64(fpsCounter) / time.Since(fpsTicker).Seconds()
				fpsText = fmt.Sprintf("FPS: %.1f", currentFPS)

				// Reset the bucket
				fpsCounter = 0
				fpsTicker = time.Now()
			}

			// 3. Draw the Overlay (if enabled)
			if showFPS {
				// FIX: If image is Grayscale (1-channel), convert to BGR (3-channel).
				// Otherwise, Green text (0, 255, 0) is drawn as Black (0) on a Black background.
				if m.Channels() == 1 {
					gocv.CvtColor(m, &m, gocv.ColorGrayToBGR)
				}

				// Draw drop shadow first (black), then text (green) for readability
				gocv.PutText(&m, fpsText, image.Pt(11, 31), gocv.FontHersheyPlain, 1.5, blackShadow, 3)
				gocv.PutText(&m, fpsText, image.Pt(10, 30), gocv.FontHersheyPlain, 1.5, green, 2)
			}

			// 4. Record (Smart Recorder handles format changes)
			if a.Config.App.Record {
				a.Recorder.Write(m)
			}

			if a.Config.Stream.Enabled {
				a.Streamer.Broadcast(m, a.Config.Stream.Quality)
			}

			// 5. Display
			a.Display.Show(m)

			// 6. Handle Input
			// Use the calculated 'delay' from earlier in the function
			key := a.Display.Key(delay)

			// Quit on 'q' or Esc (27)
			if key == 27 || key == 'q' || key == 'Q' {
				return nil
			}

			// Toggle FPS on 'f'
			if key == 'f' || key == 'F' {
				showFPS = !showFPS
			}

			m.Close()
		}
	}
}

// watchConfig monitors the config file and safely replaces the pipeline on change.
func (a *App) watchConfig() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("Failed to create config watcher: %v", err)
		return
	}
	defer watcher.Close()

	if err := watcher.Add(a.configPath); err != nil {
		log.Printf("Failed to add config path to watcher: %v", err)
		return
	}

	var last time.Time
	for {
		select {
		case ev, ok := <-watcher.Events:
			if !ok || ev.Op&fsnotify.Write == 0 {
				continue
			}
			// Debounce rapid saves
			if time.Since(last) < 200*time.Millisecond {
				continue
			}
			last = time.Now()

			// 1. Load Config
			cfg, err := config.Load(a.configPath)
			if err != nil {
				log.Printf("❌ Config reload failed: %v", err)
				continue // Skip to next event
			}

			// 2. Build Pipeline
			steps, err := builder.BuildPipeline(cfg)
			if err != nil {
				// CRITICAL: Log the error here so the user sees it!
				log.Printf("Pipeline build failed (config ignored): %v", err)
				continue // Keep running with the OLD pipeline
			}

			// 3. Swap Pipeline
			newP := pipeline.New(steps)

			a.mu.Lock()
			old := a.Pipeline
			a.Pipeline = newP
			a.Config = cfg
			a.mu.Unlock()

			// 4. Cleanup Old Pipeline safely
			if old != nil {
				// Give the running loop a moment to finish using the old pipeline
				time.AfterFunc(150*time.Millisecond, func() {
					old.Close()
				})
			}

			log.Println("Pipeline hot-reloaded successfully!")

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Watcher error: %v", err)
		}
	}
}
