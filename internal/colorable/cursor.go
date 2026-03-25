// cursor movement handler
// It is used to make the ANSI codes that move the cursor work in Windows,
// because Windows does not understand them directly.

package colorable

import (
	"strconv"
	"unsafe"
)

// Moves the cursor according to the ANSI command.
func (w *writer) handleCursor(cmd byte, params string) error {

	var csbi consoleScreenBufferInfo

	procGetConsoleScreenBufferInfo.Call(
		uintptr(w.handle),
		uintptr(unsafe.Pointer(&csbi)),
	)

	n, _ := strconv.Atoi(params)

	if n == 0 {
		n = 1
	}

	switch cmd {

	case 'A':
		csbi.cursorPosition.y -= short(n)

	case 'B':
		csbi.cursorPosition.y += short(n)

	case 'C':
		csbi.cursorPosition.x += short(n)

	case 'D':
		csbi.cursorPosition.x -= short(n)

	case 'G':
		csbi.cursorPosition.x = short(n - 1)

	case 'H', 'f':
		csbi.cursorPosition.x = 0
		csbi.cursorPosition.y = 0

	case 'E':
		csbi.cursorPosition.x = 0
		csbi.cursorPosition.y += short(n)

	case 'F':
		csbi.cursorPosition.x = 0
		csbi.cursorPosition.y -= short(n)
	}

	procSetConsoleCursorPosition.Call(
		uintptr(w.handle),
		*(*uintptr)(unsafe.Pointer(&csbi.cursorPosition)),
	)

	return nil
}

// Saves the current cursor position.
func (w *writer) saveCursor() error {

	var csbi consoleScreenBufferInfo

	procGetConsoleScreenBufferInfo.Call(
		uintptr(w.handle),
		uintptr(unsafe.Pointer(&csbi)),
	)

	w.oldpos = csbi.cursorPosition

	return nil
}

// Returns the cursor to the saved position.
func (w *writer) restoreCursor() error {

	procSetConsoleCursorPosition.Call(
		uintptr(w.handle),
		*(*uintptr)(unsafe.Pointer(&w.oldpos)),
	)

	return nil
}

// Handles save and restore commands
func (w *writer) handleSaveRestore(cmd byte) error {

	var csbi consoleScreenBufferInfo

	switch cmd {

	case 's':

		procGetConsoleScreenBufferInfo.Call(
			uintptr(w.handle),
			uintptr(unsafe.Pointer(&csbi)),
		)

		w.oldpos = csbi.cursorPosition

	case 'u':

		procSetConsoleCursorPosition.Call(
			uintptr(w.handle),
			*(*uintptr)(unsafe.Pointer(&w.oldpos)),
		)

	}

	return nil
}
