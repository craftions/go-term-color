//go:build windows
// +build windows

package color

import (
	"os"

	"golang.org/x/sys/windows"
)

func init() {
	// Enable Virtual Terminal Processing (native support for ANSI codes)
	// for Windows stdout. Works on Windows 10 and above.
	stdout := windows.Handle(os.Stdout.Fd())
	var originalMode uint32
	if err := windows.GetConsoleMode(stdout, &originalMode); err == nil {
		windows.SetConsoleMode(stdout, originalMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
	}

	// We also enable it for stderr
	stderr := windows.Handle(os.Stderr.Fd())
	if err := windows.GetConsoleMode(stderr, &originalMode); err == nil {
		windows.SetConsoleMode(stderr, originalMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
	}
}
