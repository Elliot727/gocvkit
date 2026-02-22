# GoCVKit

**The OpenCV framework Go actually needed.**
No boilerplate. Hot reload. Pre-allocated pipelines. No memory leaks.

---

## What is it?

GoCVKit is a modular computer vision framework built on `gocv`. It replaces the typical spaghetti code of camera loops, mutexes, and manual resource management with a declarative TOML configuration and a robust Go API.

It handles the boring stuff (camera init, window management, config watching, graceful shutdown) so you can focus on the image processing logic.

## Quick Start

Get a live camera feed with edge detection in 10 lines of code.

```go
package main

import (
    "log"
    "github.com/Elliot727/gocvkit"
)

func main() {
    app, err := gocvkit.NewApp("config.toml")
    if err != nil {
        log.Fatal(err)
    }
    defer app.Close()

    // Blocks until 'q' or Ctrl+C
    if err := app.Run(nil); err != nil {
        log.Fatal(err)
    }
}
```

That's it. No manual `Mat` closing, no race conditions, no leaky goroutines.

## Installation

```bash
go get github.com/elliot727/gocvkit
```

### Prerequisites
You need OpenCV 4.x installed on your system. GoCV does not bundle it.
- **macOS**: `brew install opencv`
- **Ubuntu**: `sudo apt install libopencv-dev`
- **Windows**: Follow the [gocv Windows guide](https://gocv.io/getting-started/windows/)

## Configuration

Define your pipeline in `config.toml`. Change this file while the app is running to see instant updates.

```toml
[app]
window_name = "GoCVKit – Edge Detection"
record = true          # Optional: Record output
output = "capture.mp4"

[camera]
device_id = 0
# file = "input.mp4"   # Or process a video file

[stream]
enabled = true         # Optional: MJPEG Stream
port = 8080
path = "/stream"
quality = 75

# Define the processing chain
[[pipeline.steps]]
name = "Grayscale"

[[pipeline.steps]]
name = "GaussianBlur"
kernel = 9
sigma = 1.8

[[pipeline.steps]]
name = "Canny"
low = 50
high = 150
```

## Controls

- **`q`** or **`Esc`**: Quit cleanly.
- **`f`**: Toggle FPS overlay.

## Key Features

- **Declarative Pipelines**: Define complex CV chains in TOML.
- **Hot Reloading**: Edit config → Save → Pipeline updates instantly. Invalid configs are rejected with clear errors; the old pipeline keeps running.
- **Pre-Allocated Hot Path**: All heavy resources (kernels, buffers) are allocated at startup. The frame loop performs **near-zero heap allocations**, minimizing GC pressure.
- **Robust Error Handling**: Typos in config names or parameters fail fast with descriptive messages, not segfaults.
- **Built-in Streaming & Recording**: MJPEG HTTP server and MP4 recording out of the box.
- **Extensible Plugin System**: Register custom filters from anywhere in your codebase.
- **Graceful Shutdown**: Handles SIGINT/SIGTERM properly. Resources are always freed.

## Built-in Processors

| Category | Processor | Config Params | Description |
| :--- | :--- | :--- | :--- |
| **Core** | `Grayscale` | – | BGR → Gray conversion |
| | `Flip` | `mode` | Horizontal, vertical, or both |
| | `Resize` | `width`, `height` | Fixed resolution resize |
| | `Rotate` | `angle` | Arbitrary rotation (optimized for 90° steps) |
| | `ColorConvert` | `code` | HSV, LAB, YUV, etc. |
| **Blurs** | `GaussianBlur` | `kernel`, `sigma` | Standard Gaussian smoothing |
| | `MedianBlur` | `k` | Salt-and-pepper noise removal |
| | `Bilateral` | `diameter`, `sigmas` | Edge-preserving smoothing |
| **Morphology**| `Erode` | `kernel`, `iterations` | Morphological erosion |
| | `Dilate` | `kernel`, `iterations` | Morphological dilation |
| | `MorphClose` | `kernel`, `iterations` | Close gaps in foreground |
| **Edges** | `Canny` | `low`, `high` | Standard edge detection |
| | `Sobel` | `sobel_size` | Gradient calculation |
| | `Laplacian` | `k` | Second derivative edges |
| | `Scharr` | – | High-precision gradients |
| **Threshold** | `Otsu` | `max_value`, `invert` | Automatic thresholding |
| | `Adaptive` | `block_size`, `c` | Local adaptive thresholding |
| **Advanced** | `BackgroundSubtractor`| `algorithm`, `lr` | MOG2 or KNN motion detection |

## Advanced Usage

### Frame Callbacks
Inject custom logic (overlays, logging, external APIs) into the render loop without modifying the pipeline.

```go
app.Run(func(frame *gocv.Mat) {
    // Draw custom HUD, send to API, etc.
    // Frame is valid and processed.
})
```

### Custom Filters
Implement the `Processable` interface. GoCVKit handles the reflection, config parsing, and lifecycle management.

```go
type RedTint struct {
    Intensity float64 `toml:"intensity"`
}

// Validate runs once at startup. Return error to stop invalid configs.
func (r *RedTint) Validate() error {
    if r.Intensity < 0 || r.Intensity > 1 {
        return fmt.Errorf("intensity must be 0.0–1.0")
    }
    return nil
}

// Process runs every frame. Near-zero allocs expected.
func (r *RedTint) Process(src gocv.Mat, dst *gocv.Mat) error {
    if src.Empty() { return nil }
    // ... implementation ...
    return nil
}

// Close runs on shutdown/reload. Free C++ resources here.
func (r *RedTint) Close() { /* cleanup */ }

func init() {
    gocvkit.RegisterProcessor("RedTint", &RedTint{Intensity: 0.5})
}
```

## Architecture

Designed for stability and performance:
- **`app`**: Orchestrator. Manages concurrency, signals, and lifecycle.
- **`pipeline`**: Double-buffered execution engine. Swaps pre-allocated mats to avoid per-frame mallocs.
- **`processor`**: Plugin registry. Uses reflection to map TOML params to structs safely.
- **`builder`**: Constructs pipelines from config, validating every step before execution.

## Performance Reality Check

We claim **"Near-Zero Allocation"**, not "Zero".
- **True**: We pre-allocate image buffers and kernels. The hot path avoids explicit `make` or `new` for image data.
- **False**: Go's runtime, `time.Now()`, channel sends, and cgo boundary crossings still incur minor overhead.
- **Result**: Deterministic latency suitable for real-time video (30+ FPS on modern hardware), with minimal GC stutter.

## Use Cases

- Rapid prototyping of CV algorithms.
- Home automation / Security cameras.
- Educational tools for teaching OpenCV.
- Production services where stability matters more than cleverness.

## License

MIT © 2025 elliot727

---

**Star it if it saves you time.**
**Contribute if you find a bug.**
**Stop writing camera loops by hand.**

Made by [@elliot727](https://github.com/elliot727)
