//go:build windows
// +build windows

package color

import (
	"os"

	"golang.org/x/sys/windows"
)

func init() {
	// Habilitar Virtual Terminal Processing (soporte nativo de códigos ANSI) 
	// para el stdout de Windows. Funciona a partir de Windows 10.
	stdout := windows.Handle(os.Stdout.Fd())
	var originalMode uint32
	if err := windows.GetConsoleMode(stdout, &originalMode); err == nil {
		windows.SetConsoleMode(stdout, originalMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
	}

	// Lo habilitamos también para el stderr
	stderr := windows.Handle(os.Stderr.Fd())
	if err := windows.GetConsoleMode(stderr, &originalMode); err == nil {
		windows.SetConsoleMode(stderr, originalMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
	}
}
