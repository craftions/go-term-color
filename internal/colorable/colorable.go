//go:build windows && !appengine
// +build windows,!appengine

// Makes ANSI color codes work on ANY version of Windows.
// The program writes text with color codes
// NewColorable() analyzes the terminal:
// Is it a modern terminal? (Windows 10+)
// Yes -> Uses the terminal directly (fast)
// No -> Wraps everything in a "translator"
// The translator converts ANSI codes to Windows commands
// Colors appear correctly
package colorable

import (
	"bytes"
	"io"
	"os"
	"sync"
	"unsafe"

	syscall "golang.org/x/sys/windows"

	"github.com/craftions/go-term-check/terminal"
)

// writer provides colorable Writer to the console
type writer struct {
	out       io.Writer
	handle    syscall.Handle
	althandle syscall.Handle
	oldattr   word
	oldpos    coord
	rest      bytes.Buffer
	mutex     sync.Mutex
}

// NewColorable returns new instance of writer which handles escape sequence from File.
func NewColorable(file *os.File) io.Writer {
	if file == nil {
		panic("nil passed instead of *os.File to NewColorable()")
	}

	fd := file.Fd()
	if !terminal.IsTerminal(fd) {
		return file
	}

	// Verify virtual terminal processing support
	var mode uint32
	ret, _, _ := procGetConsoleMode.Call(fd, uintptr(unsafe.Pointer(&mode)))
	if ret != 0 && mode&cENABLE_VIRTUAL_TERMINAL_PROCESSING != 0 {
		return file
	}

	// If there's no support, use the custom writer
	var csbi consoleScreenBufferInfo
	handle := syscall.Handle(fd)
	procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))

	return &writer{
		out:     file,
		handle:  handle,
		oldattr: csbi.attributes,
		oldpos:  coord{0, 0},
	}
}

// NewColorableStdout returns new instance of writer which handles escape sequence for stdout.
func NewColorableStdout() io.Writer {
	return NewColorable(os.Stdout)
}

// NewColorableStderr returns new instance of writer which handles escape sequence for stderr.
func NewColorableStderr() io.Writer {
	return NewColorable(os.Stderr)
}
