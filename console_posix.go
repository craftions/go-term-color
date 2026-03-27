//go:build !windows
// +build !windows

package color

func init() {
	// Los sistemas POSIX (Linux, macOS, Solaris) operan de forma nativa con secuencias ANSI de escape
	// por lo tanto, no se requiere configuración adicional a nivel de kernel/consola de SO aquí.
}
