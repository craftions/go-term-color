//go:build windows && !appengine
// +build windows,!appengine

// Hace que los códigos de color ANSI funcionen en CUALQUIER versión de Windows
// El programa escribe texto con códigos de color
// NewColorable() analiza la terminal
// ¿Es terminal moderna? (Windows 10+)
// Sí → Usa la terminal directamente (rápido)
// No → Envuelve todo en un "traductor"
// El traductor convierte códigos ANSI a comandos de Windows
// Los colores aparecen correctamente
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

	// Verificar soporte de virtual terminal processing
	var mode uint32
	ret, _, _ := procGetConsoleMode.Call(fd, uintptr(unsafe.Pointer(&mode)))
	if ret != 0 && mode&cENABLE_VIRTUAL_TERMINAL_PROCESSING != 0 {
		return file
	}

	// Si no hay soporte, usar el writer personalizado
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
