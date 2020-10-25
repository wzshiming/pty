package pty

import (
	"errors"
	"io"
	"os"

	"github.com/wzshiming/loginshell"
)

// Pty communication interface
type Pty interface {
	// Start
	Start() error

	// Stdio returns stdio of pty
	Stdio() (io.ReadWriteCloser, error)

	// SetSize sets the console size
	SetSize(cols uint32, rows uint32) error

	// GetSize gets the console size
	GetSize() (uint32, uint32, error)

	// Process returns the process
	Process() (*os.Process, error)
}

type Options struct {
	// Command holds command line arguments, including the command as Command[0].
	// If the Command field is empty or nil, Run uses {Path}.
	Command []string

	// Dir sets the current working directory for the command
	Dir string

	// Env sets the environment variables. Use the format VAR=VAL.
	Env []string

	// Initial size for Columns and Rows
	Cols uint32
	Rows uint32
}

// NewPty creates a new pty
func NewPty(opt Options) (Pty, error) {
	if opt.Env == nil {
		opt.Env = os.Environ()
	}
	if opt.Dir == "" {
		opt.Dir = "."
	}
	if opt.Cols == 0 {
		opt.Cols = 80
	}
	if opt.Rows == 0 {
		opt.Rows = 24
	}
	if opt.Command == nil {
		sh, err := loginshell.Shell()
		if err != nil {
			return nil, err
		}
		opt.Command = []string{sh}
	}
	return newPty(&opt)
}

var (
	ErrProcessNotStarted = errors.New("process has not been started")
	ErrInvalidCmd        = errors.New("invalid command")
)
