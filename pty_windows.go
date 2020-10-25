package pty

import (
	"io"
	"os"
	"strings"

	winpty "github.com/iamacarpet/go-winpty"
	"golang.org/x/sys/windows"
)

type console struct {
	file *winpty.WinPTY
	opt  *Options
}

func newPty(opt *Options) (Pty, error) {
	return &console{
		opt:  opt,
		file: nil,
	}, nil
}

func (c *console) Start() error {
	if c.file != nil {
		return nil
	}

	opts := winpty.Options{
		InitialCols: c.opt.Cols,
		InitialRows: c.opt.Rows,
		Command:     strings.Join(c.opt.Command, " "),
		Dir:         c.opt.Dir,
		Env:         c.opt.Env,
	}

	cmd, err := winpty.OpenWithOptions(opts)
	if err != nil {
		return err
	}

	c.file = cmd
	return nil
}

func (c *console) Stdio() (io.ReadWriteCloser, error) {
	if c.file == nil {
		return nil, ErrProcessNotStarted
	}

	file := struct {
		io.Reader
		io.Writer
		io.Closer
	}{
		Reader: c.file.StdOut,
		Writer: c.file.StdIn,
		Closer: closer(c.file.Close),
	}
	return file, nil
}

func (c *console) SetSize(cols uint32, rows uint32) error {
	if c.file == nil {
		return ErrProcessNotStarted
	}

	c.opt.Cols = cols
	c.opt.Rows = rows
	c.file.SetSize(cols, rows)
	return nil
}

func (c *console) GetSize() (uint32, uint32, error) {
	if c.file == nil {
		return 0, 0, ErrProcessNotStarted
	}

	return c.opt.Cols, c.opt.Rows, nil
}

func (c *console) Process() (*os.Process, error) {
	if c.file == nil {
		return nil, ErrProcessNotStarted
	}

	handle := c.file.GetProcHandle()
	pid, err := windows.GetProcessId(windows.Handle(handle))
	if err != nil {
		return nil, err
	}

	return os.FindProcess(int(pid))
}

type closer func()

func (c closer) Close() error {
	c()
	return nil
}
