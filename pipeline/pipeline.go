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
	"github.com/Elliot727/gocvkit/processor"

	"gocv.io/x/gocv"
)

// Pipeline holds an ordered list of processing steps and two reusable buffers.
type Pipeline struct {
	Steps []processor.Step // Steps contains the ordered list of processing steps to execute
	bufA  gocv.Mat         // bufA is the first internal scratch buffer for double-buffering
	bufB  gocv.Mat         // bufB is the second internal scratch buffer for double-buffering
}

// New creates a new pipeline from a slice of processing steps.
// The two internal buffers are pre-allocated and reused for the lifetime of the pipeline.
func New(steps []processor.Step) *Pipeline {
	return &Pipeline{
		Steps: steps,
		bufA:  gocv.NewMat(),
		bufB:  gocv.NewMat(),
	}
}

// Close releases the internal scratch buffers.
// Safe to call multiple times.
func (p *Pipeline) Close() {
	p.bufA.Close()
	p.bufB.Close()
}

// Run executes the full pipeline on src and writes the final result to dst.
//
// The function uses zero-allocation double-buffering: each step alternately
// reads from one buffer and writes to the other. The original src Mat is never
// modified and can be safely reused or closed by the caller.
func (p *Pipeline) Run(src gocv.Mat, dst *gocv.Mat) error {
	if src.Empty() {
		return nil
	}

	// Copy source into first buffer
	src.CopyTo(&p.bufA)

	in := &p.bufA
	out := &p.bufB

	// Execute each step, swapping buffers
	for _, step := range p.Steps {
		step.Process(*in, out)
		in, out = out, in
	}

	// Copy final result to destination
	in.CopyTo(dst)
	return nil
}
