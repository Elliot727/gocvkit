# GoCVKit v0.1.0
**The OpenCV framework Go was waiting for**
Zero boilerplate • Hot reload • Zero leaks • Actually fun to use

---

## What is GoCVKit?

GoCVKit is a modular GoCV framework for working with live camera or video streams with configurable pipelines and filters. It simplifies computer vision development by providing a clean, idiomatic API that handles all the complexity behind the scenes.

## Quick Start

To get started with a live camera application:

```go
package main

import (
    "log"
    "github.com/Elliot727/gocvkit"
)

func main() {
    app, err := gocvkit.NewApp("config.toml")
    if err != nil {
        log.Fatal("Failed to create app:", err)
    }
    defer app.Close()

    // Run the application with optional per-frame callback
    if err := app.Run(nil); err != nil {
        log.Fatal("Application error:", err)
    }
}
```

That's it! No boilerplate, no crashes, no resource leaks.

## Installation

```bash
go get github.com/elliot727/gocvkit
```

### Prerequisites

GoCVKit requires OpenCV 4.x:
- **macOS**: `brew install opencv`
- **Ubuntu**: `sudo apt install libopencv-dev`
- **Windows**: [Installation guide for Windows](https://gocv.io/getting-started/windows/)

## Configuration

Create a `config.toml` file to define your pipeline:

```toml
[app]
window_name = "GoCVKit – Live Edge Detection"

[camera]
device_id = 0
# file = "demo.mp4"  # Use video file instead of webcam

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

## Key Features

- **One-line startup**: Full application in ≤10 lines of code
- **Live config hot-reload**: Edit `config.toml` → instant pipeline update
- **Zero per-frame allocations**: Efficient double-buffer pipeline
- **Frame callbacks**: Overlay, logging, and recording support
- **Graceful shutdown**: Ctrl+C and Esc/q handling with zero resource leaks
- **Extensible plugin system**: Register custom processors from anywhere
- **Zero-boilerplate AutoConfig**: Dynamic parameters with reflection
- **Robust error handling**: Clear error messages for typos and unknown processors
- **Webcam & video file support**: Seamless input switching

## Built-in Processors

| Name           | Config Keys             | Description                          |
|----------------|-------------------------|--------------------------------------|
| `Grayscale`    | –                       | Convert BGR → grayscale              |
| `GaussianBlur` | `kernel`, `sigma`       | Noise reduction with Gaussian kernel |
| `MedianBlur`   | `k`                     | Remove salt-and-pepper noise         |
| `Canny`        | `low`, `high`           | Edge detection                       |
| `Sobel`        | `sobel_size`            | Gradient-based edge detection        |

## Advanced Usage

### Frame Callbacks

Add real-time overlays, logging, or recording:

```go
app.Run(func(frame *gocv.Mat) {
    // Process the final frame
    // Useful for overlays, saving, logging, etc.
})
```

### Adding Custom Filters

Create your own processors with zero boilerplate:

```go
// custom_filter.go
package main

import (
    "gocv.io/x/gocv"
    "github.com/Elliot727/gocvkit"
)

// RedTint adds a red tint to the image
type RedTint struct {
    Intensity float64 `toml:"intensity"` // 0.0 – 1.0
    Enabled   bool    `toml:"enabled"`
}

func (r *RedTint) Process(src gocv.Mat, dst *gocv.Mat) {
    if !r.Enabled {
        src.CopyTo(dst)
        return
    }
    src.CopyTo(dst)
    // Apply red tint logic here
    // Implementation details...
}

func init() {
    gocvkit.RegisterProcessor("RedTint", &RedTint{
        Intensity: 0.5,
        Enabled:   true,
    })
}
```

Then add to your config.toml:

```toml
[[pipeline.steps]]
name = "RedTint"
intensity = 0.85
enabled = true
```

## Architecture

GoCVKit follows a clean, modular architecture:

- **app**: Main orchestrator handling camera, display, pipeline, and config
- **builder**: Creates processing pipelines from configuration
- **camera**: Unified wrapper for webcam and video file input
- **config**: TOML configuration loading with custom unmarshaling
- **display**: Window display wrapper
- **pipeline**: Efficient double-buffered processing pipeline
- **processor**: Extensible filter system with auto-configuration

## Use Cases

Perfect for:

- **Rapid prototyping** of computer vision applications
- **Teaching computer vision** concepts
- **Live demonstrations** and presentations
- **Real-time vision applications**
- **Anyone who values their sanity**

## License

MIT © 2025 elliot727

---

**Star if you like it**
**Contribute if you love it**
**Tell everyone** — Go deserves this.

Made with passion by [@elliot727](https://github.com/elliot727)

**GoCVKit v0.1.0 — The future of real-time computer vision in Go starts here.**
