package iostreams

import (
	"io"
	"os"
)

// IOStreams represents the standard input/output streams
type IOStreams struct {
	In     io.Reader
	Out    io.Writer
	ErrOut io.Writer
}

// New creates a new IOStreams instance with standard streams
func New() *IOStreams {
	return &IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}
}

// TestOptions represents options for test IOStreams
type TestOptions struct {
	In     io.Reader
	Out    io.Writer
	ErrOut io.Writer
}

// NewTest creates a new IOStreams for testing
func NewTest(options TestOptions) *IOStreams {
	io := &IOStreams{
		In:     options.In,
		Out:    options.Out,
		ErrOut: options.ErrOut,
	}
	
	// Use defaults if not provided
	if io.In == nil {
		io.In = &io.PipeReader{}
	}
	if io.Out == nil {
		io.Out = &io.PipeWriter{}
	}
	if io.ErrOut == nil {
		io.ErrOut = &io.PipeWriter{}
	}
	
	return io
}

// ColorEnabled returns true if color output is enabled
func (io *IOStreams) ColorEnabled() bool {
	// For demo purposes, always enable color
	return true
}

// IsTerminal returns true if output is a terminal
func (io *IOStreams) IsTerminal() bool {
	// For demo purposes, assume it's a terminal
	return true
}

// StartProgressIndicator starts a progress indicator
func (io *IOStreams) StartProgressIndicator() {
	// For demo, no-op
}

// StopProgressIndicator stops a progress indicator
func (io *IOStreams) StopProgressIndicator() {
	// For demo, no-op
}