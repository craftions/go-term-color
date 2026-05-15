//go:build windows && !appengine
// +build windows,!appengine

package colorable

import (
	"io"
	"os"

	"github.com/craftions/go-term-check/terminal"
	"golang.org/x/sys/windows"
)

var (
	isTerminal = func(fd uintptr) bool {
		return terminal.Check(fd).Terminal
	}
	getConsoleMode = windows.GetConsoleMode
	setConsoleMode = windows.SetConsoleMode
)

// NewColorable returns new instance of writer which handles escape sequence from File.
// For modern Windows terminals, it enables VIRTUAL_TERMINAL_PROCESSING.
func NewColorable(file *os.File) io.Writer {
	if file == nil {
		panic("nil passed instead of *os.File to NewColorable()")
	}

	fd := file.Fd()
	if !isTerminal(fd) {
		return file
	}

	// Verify and enable virtual terminal processing support
	var mode uint32
	err := getConsoleMode(windows.Handle(fd), &mode)
	if err == nil {
		mode |= windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
		_ = setConsoleMode(windows.Handle(fd), mode)
	}

	return file
}

// NewColorableStdout returns new instance of writer which handles escape sequence for stdout.
func NewColorableStdout() io.Writer {
	return NewColorable(os.Stdout)
}

// NewColorableStderr returns new instance of writer which handles escape sequence for stderr.
func NewColorableStderr() io.Writer {
	return NewColorable(os.Stderr)
}
