// +build !windows

package pty

import (
	"io"
	"os"
	"os/exec"

	"github.com/creack/pty"
)

type console struct {
	file *os.File
	cmd  *exec.Cmd
	opt  *Options
}

func newPty(opt *Options) (Pty, error) {
	return &console{
		opt: opt,
	}, nil
}

// Start starts a process and wraps in a console
func (c *console) Start() error {
	if c.file != nil {
		return nil
	}

	if len(c.opt.Command) < 1 {
		return ErrInvalidCmd
	}

	cmd := exec.Command(c.opt.Command[0], c.opt.Command[1:]...)
	cmd.Dir = c.opt.Dir
	cmd.Env = c.opt.Env

	f, err := pty.StartWithSize(cmd, &pty.Winsize{Cols: uint16(c.opt.Cols), Rows: uint16(c.opt.Rows)})
	if err != nil {
		return err
	}

	c.cmd = cmd
	c.file = f
	return nil
}

func (c *console) Stdio() (io.ReadWriteCloser, error) {
	if c.file == nil {
		return nil, ErrProcessNotStarted
	}

	return c.file, nil
}

func (c *console) SetSize(cols uint32, rows uint32) error {
	if c.cmd == nil {
		return ErrProcessNotStarted
	}

	return pty.Setsize(c.file, &pty.Winsize{Cols: uint16(cols), Rows: uint16(rows)})
}

func (c *console) GetSize() (uint32, uint32, error) {
	if c.cmd == nil {
		return 0, 0, ErrProcessNotStarted
	}

	cols, rows, err := pty.Getsize(c.file)
	if err != nil {
		return 0, 0, err
	}

	return uint32(cols), uint32(rows), nil
}

func (c *console) Process() (*os.Process, error) {
	if c.cmd == nil {
		return nil, ErrProcessNotStarted
	}

	return c.cmd.Process, nil
}
