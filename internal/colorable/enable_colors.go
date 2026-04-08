//go:build windows && !appengine
// +build windows,!appengine

// Este archivo activa el modo de colores moderno en la consola de Windows cuando es posible.
// Las versiones recientes de Windows (10 y 11) tienen una capacidad llamada "VT mode"
// que permite que la consola entienda directamente los códigos de color ANSI (como \033[32m)
// Sin embargo, esta capacidad viene desactivada por defecto por compatibilidad con programas antiguos
package colorable

import (
	"os"
	"unsafe"
)

// EnableColorsStdout enable colors if possible.
func EnableColorsStdout(enabled *bool) func() {
	var mode uint32
	h := os.Stdout.Fd()
	if r, _, _ := procGetConsoleMode.Call(h, uintptr(unsafe.Pointer(&mode))); r != 0 {
		if r, _, _ = procSetConsoleMode.Call(h, uintptr(mode|cENABLE_VIRTUAL_TERMINAL_PROCESSING)); r != 0 {
			if enabled != nil {
				*enabled = true
			}
			return func() {
				procSetConsoleMode.Call(h, uintptr(mode))
			}
		}
	}
	if enabled != nil {
		*enabled = true
	}
	return func() {}
}
