// handle ANSI commands that delete text in the Windows console.
package colorable

import (
	"unsafe"
)

// handle ANSI erase commands.
func (w *writer) handleErase(cmd byte, params string) error {

	var csbi consoleScreenBufferInfo
	var written dword

	procGetConsoleScreenBufferInfo.Call(
		uintptr(w.handle),
		uintptr(unsafe.Pointer(&csbi)),
	)

	cursor := csbi.cursorPosition

	switch cmd {

	case 'J':

		procFillConsoleOutputCharacter.Call(
			uintptr(w.handle),
			uintptr(' '),
			uintptr(1000),
			*(*uintptr)(unsafe.Pointer(&cursor)),
			uintptr(unsafe.Pointer(&written)),
		)

	case 'K':

		procFillConsoleOutputCharacter.Call(
			uintptr(w.handle),
			uintptr(' '),
			uintptr(200),
			*(*uintptr)(unsafe.Pointer(&cursor)),
			uintptr(unsafe.Pointer(&written)),
		)

	case 'X':

		procFillConsoleOutputCharacter.Call(
			uintptr(w.handle),
			uintptr(' '),
			uintptr(10),
			*(*uintptr)(unsafe.Pointer(&cursor)),
			uintptr(unsafe.Pointer(&written)),
		)

	}

	return nil
}
