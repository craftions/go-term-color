// Show or hide the console cursor in Windows when using ANSI commands.
package colorable

import (
	"unsafe"
)

func (w *writer) handleMode(cmd byte, params string) error {

	var ci consoleCursorInfo

	switch cmd {

	// show cursor
	case 'h':

		if params == "?25" {

			procGetConsoleCursorInfo.Call(
				uintptr(w.handle),
				uintptr(unsafe.Pointer(&ci)),
			)

			ci.visible = 1

			procSetConsoleCursorInfo.Call(
				uintptr(w.handle),
				uintptr(unsafe.Pointer(&ci)),
			)
		}

	// hide cursor
	case 'l':

		if params == "?25" {

			procGetConsoleCursorInfo.Call(
				uintptr(w.handle),
				uintptr(unsafe.Pointer(&ci)),
			)

			ci.visible = 0

			procSetConsoleCursorInfo.Call(
				uintptr(w.handle),
				uintptr(unsafe.Pointer(&ci)),
			)
		}

	}

	return nil
}
