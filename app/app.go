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
	"log"
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

	"github.com/fsnotify/fsnotify"
	"gocv.io/x/gocv"
)

// App represents a fully configured and running computer vision application.
type App struct {
	mu         sync.RWMutex       // mu provides thread-safe access to mutable fields
	Camera     *camera.Camera     // Camera handles video input from webcam or file
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

	win := display.New(cfg.App.WindowName)

	steps, err := builder.BuildPipeline(cfg)
	if err != nil {
		cam.Close()
		win.Close()
		return nil, err
	}

	a := &App{
		Camera:     cam,
		Display:    win,
		Pipeline:   pipeline.New(steps),
		Config:     cfg,
		configPath: cfgPath,
	}

	go a.watchConfig() // fire-and-forget hot reload
	return a, nil
}

// Close releases all resources (camera, window, pipeline).
func (a *App) Close() {
	a.Camera.Close()
	a.Display.Close()

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
	if a.Config.Camera.File != "" {
		fps := a.Camera.FPS()

		// Fallback for files with missing/bad metadata
		if fps <= 0 || fps > 200 {
			fps = 30.0
		}

		// 1000ms / FPS = delay in ms (e.g., 30fps -> 33ms)
		delay = int(1000.0 / fps)
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

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case m, ok := <-results:
			if !ok {
				return nil // Pipeline finished (e.g. end of video)
			}

			frameCallback(&m)
			a.Display.Show(m)

			// Wait for the correct duration (1ms for webcam, ~33ms for video)
			if key := a.Display.Key(delay); key == 27 || key == 'q' || key == 'Q' {
				return nil
			}

			m.Close()
		}
	}
}

// watchConfig monitors the config file and safely replaces the pipeline on change.
func (a *App) watchConfig() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}
	defer watcher.Close()

	if err := watcher.Add(a.configPath); err != nil {
		return
	}

	var last time.Time
	for {
		select {
		case ev, ok := <-watcher.Events:
			if !ok || ev.Op&fsnotify.Write == 0 {
				continue
			}
			if time.Since(last) < 200*time.Millisecond {
				continue
			}
			last = time.Now()

			if cfg, err := config.Load(a.configPath); err == nil {
				if steps, err := builder.BuildPipeline(cfg); err == nil {
					newP := pipeline.New(steps)

					a.mu.Lock()
					old := a.Pipeline
					a.Pipeline = newP
					a.Config = cfg
					a.mu.Unlock()

					if old != nil {
						time.AfterFunc(150*time.Millisecond, old.Close)
					}

					log.Println("Pipeline hot-reloaded!")
				}
			}
		case <-watcher.Errors:
			// ignore
		}
	}
}
