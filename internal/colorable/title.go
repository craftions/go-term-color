// Change the console window title using ANSI codes.
package colorable

import (
	"bytes"
	"unsafe"

	syscall "golang.org/x/sys/windows"
)

// `\033]0;TITLESTR\007`
func doTitleSequence(er *bytes.Reader) error {
	var c byte
	var err error

	c, err = er.ReadByte()
	if err != nil {
		return err
	}
	if c != '0' && c != '2' {
		return nil
	}
	c, err = er.ReadByte()
	if err != nil {
		return err
	}
	if c != ';' {
		return nil
	}
	title := make([]byte, 0, 80)
	for {
		c, err = er.ReadByte()
		if err != nil {
			return err
		}
		if c == 0x07 || c == '\n' {
			break
		}
		title = append(title, c)
	}
	if len(title) > 0 {
		title8, err := syscall.UTF16PtrFromString(string(title))
		if err == nil {
			procSetConsoleTitle.Call(uintptr(unsafe.Pointer(title8)))
		}
	}
	return nil
}

func (w *writer) handleTitle(er *bytes.Reader) error {

	w.rest.WriteByte(0x1b)
	w.rest.WriteByte(']')

	er.WriteTo(&w.rest)

	if bytes.IndexByte(w.rest.Bytes(), 0x07) == -1 {
		return nil
	}

	er = bytes.NewReader(w.rest.Bytes()[2:])

	err := doTitleSequence(er)

	if err != nil {
		return err
	}

	w.rest.Reset()

	return nil
}
