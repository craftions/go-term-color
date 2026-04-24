//go:build windows && !appengine
// +build windows,!appengine

// This file enables modern color mode in the Windows console when possible.
// Recent versions of Windows (10 and 11) have a capability called "VT mode"
// which allows the console to directly understand ANSI color codes (like \033[32m).
// However, this capability comes disabled by default for compatibility with older programs.
package colorable

import (
	"os"
	"unsafe"
)

// EnableColorsStdout enable colors if possible.
func EnableColorsStdout(enabled *bool) func() {
	handle := os.Stdout.Fd()
	var mode uint32

	setEnabled := func() {
		if enabled != nil {
			*enabled = true
		}
	}

	if r, _, _ := procGetConsoleMode.Call(handle, uintptr(unsafe.Pointer(&mode))); r == 0 {
		setEnabled()
		return func() {}
	}

	newMode := mode | cENABLE_VIRTUAL_TERMINAL_PROCESSING
	if r, _, _ := procSetConsoleMode.Call(handle, uintptr(newMode)); r != 0 {
		setEnabled()
		return func() {
			procSetConsoleMode.Call(handle, uintptr(mode))
		}
	}

	setEnabled()
	return func() {}
}
