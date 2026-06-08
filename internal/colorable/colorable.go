//go:build windows && !appengine
// +build windows,!appengine

package colorable

import (
	"io"
	"os"

	"golang.org/x/sys/windows"
)

var (
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

	// Verify and enable virtual terminal processing support
	var mode uint32
	err := getConsoleMode(windows.Handle(fd), &mode)
	if err == nil {
		mode |= windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
		_ = setConsoleMode(windows.Handle(fd), mode)
	}

	return file
}
