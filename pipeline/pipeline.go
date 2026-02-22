// Package pipeline executes a sequence of image processing steps with minimal memory allocations.
//
// It uses a double-buffer (ping-pong) strategy: two pre-allocated Mats (bufA and bufB)
// are alternately used as input/output for each step. This avoids creating new Mats
// on every frame and keeps per-frame overhead extremely low.
//
// Pipelines are safe to replace at runtime (via config hot-reload) because the
// buffers belong to the Pipeline instance and are closed only after the old
// pipeline is no longer in use.
package pipeline

import (
	"fmt"
	"os"
	"time"

	"github.com/Elliot727/gocvkit/processor"

	"gocv.io/x/gocv"
)

// StepStats holds performance metrics for a single processing step.
type StepStats struct {
	Name      string
	Calls     int64
	TotalTime time.Duration
	MaxTime   time.Duration
}

// Pipeline holds an ordered list of processing steps and two reusable buffers.
type Pipeline struct {
	Steps []processor.Step // Steps contains the ordered list of processing steps to execute
	bufA  gocv.Mat         // bufA is the first internal scratch buffer for double-buffering
	bufB  gocv.Mat         // bufB is the second internal scratch buffer for double-buffering
	stats []StepStats
}

// New creates a new pipeline from a slice of processing steps.
// The two internal buffers are pre-allocated and reused for the lifetime of the pipeline.
func New(steps []processor.Step) *Pipeline {
	stats := make([]StepStats, len(steps))
	for i, step := range steps {
		stats[i] = StepStats{Name: step.Name()}
	}
	return &Pipeline{
		Steps: steps,
		bufA:  gocv.NewMat(),
		bufB:  gocv.NewMat(),
	}
}

// Close releases the internal scratch buffers.
// Safe to call multiple times.
func (p *Pipeline) Close() {
	p.printReport()
	p.bufA.Close()
	p.bufB.Close()
	for _, step := range p.Steps {
		step.Close()
	}
}

// Run executes the full pipeline on src and writes the final result to dst.
//
// The function uses zero-allocation double-buffering: each step alternately
// reads from one buffer and writes to the other. The original src Mat is never
// modified and can be safely reused or closed by the caller.
// Run executes the full pipeline with profiling.
func (p *Pipeline) Run(src gocv.Mat, dst *gocv.Mat) error {
	if src.Ptr() == nil {
		return nil
	}
	if src.Empty() {
		return nil
	}
	if dst == nil {
		return fmt.Errorf("dst is nil")
	}

	if p.bufA.Empty() || p.bufB.Empty() {
		p.bufA = gocv.NewMatWithSize(src.Rows(), src.Cols(), src.Type())
		p.bufB = gocv.NewMatWithSize(src.Rows(), src.Cols(), src.Type())
	} else if p.bufA.Rows() != src.Rows() || p.bufA.Cols() != src.Cols() {
		p.bufA.Close()
		p.bufB.Close()
		p.bufA = gocv.NewMatWithSize(src.Rows(), src.Cols(), src.Type())
		p.bufB = gocv.NewMatWithSize(src.Rows(), src.Cols(), src.Type())
	}

	if len(p.Steps) == 0 {
		src.CopyTo(dst)
		return nil
	}

	if len(p.stats) != len(p.Steps) {
		p.stats = make([]StepStats, len(p.Steps))
		for i, step := range p.Steps {
			p.stats[i] = StepStats{Name: step.Name()}
		}
	}

	src.CopyTo(&p.bufA)
	in := &p.bufA
	out := &p.bufB

	for i, step := range p.Steps {
		start := time.Now()
		if in.Empty() {
			return fmt.Errorf("step %s: input mat is empty before processing", step.Name())
		}
		if err := step.Process(*in, out); err != nil {
			return fmt.Errorf("step %s failed: %w", step.Name(), err)
		}

		if out.Empty() {
			return fmt.Errorf("step %s produced an empty output matrix; pipeline halted to prevent crash", step.Name())
		}

		elapsed := time.Since(start)

		p.stats[i].Calls++
		p.stats[i].TotalTime += elapsed
		if elapsed > p.stats[i].MaxTime {
			p.stats[i].MaxTime = elapsed
		}

		in, out = out, in
	}

	in.CopyTo(dst)
	return nil
}

// printReport outputs a formatted table to stderr.
func (p *Pipeline) printReport() {
	if len(p.stats) == 0 {
		fmt.Fprintln(os.Stderr, "\n--- Pipeline Performance Report ---")
		fmt.Fprintln(os.Stderr, "No steps executed. Pipeline was empty.")
		return
	}

	fmt.Fprintln(os.Stderr, "\n--- Pipeline Performance Report ---")
	fmt.Fprintf(os.Stderr, "%-20s | %-10s | %-10s | %-10s | %-10s\n", "Step", "Calls", "Avg (ms)", "Max (ms)", "% Total")
	fmt.Fprintln(os.Stderr, "---------------------------------------------------------------")

	var grandTotal time.Duration
	for _, s := range p.stats {
		grandTotal += s.TotalTime
	}

	for _, s := range p.stats {
		avg := float64(0)
		if s.Calls > 0 {
			avg = float64(s.TotalTime.Nanoseconds()) / float64(s.Calls) / 1e6 // to ms
		}
		maxMs := float64(s.MaxTime.Nanoseconds()) / 1e6
		percent := float64(0)
		if grandTotal > 0 {
			percent = float64(s.TotalTime) / float64(grandTotal) * 100
		}

		fmt.Fprintf(os.Stderr, "%-20s | %-10d | %-10.3f | %-10.3f | %-9.2f%%\n",
			s.Name, s.Calls, avg, maxMs, percent)
	}
	fmt.Fprintln(os.Stderr, "-----------------------------------------------")

	// Safe access: we know len > 0 here because of the check at the top
	frameCount := p.stats[0].Calls
	fmt.Fprintf(os.Stderr, "Total Pipeline Time: %.2f ms over %d frames\n",
		float64(grandTotal.Nanoseconds())/1e6, frameCount)
	fmt.Fprintln(os.Stderr, "")
}
