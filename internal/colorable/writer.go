// Read what the program wants to print to the console, detect if there are ANSI codes,
// and call the correct function to make it work in Windows.
package colorable

import (
	"bytes"
	"strconv"
)

func (w *writer) Write(data []byte) (int, error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	er := w.prepareReader(data)

	err := w.processStream(er)
	if err != nil {
		return 0, err
	}

	return len(data), nil
}

func (w *writer) prepareReader(data []byte) *bytes.Reader {

	var er *bytes.Reader

	if w.rest.Len() > 0 {

		var rest bytes.Buffer

		w.rest.WriteTo(&rest)
		w.rest.Reset()

		rest.Write(data)

		er = bytes.NewReader(rest.Bytes())

	} else {

		er = bytes.NewReader(data)

	}

	return er
}

func (w *writer) processStream(er *bytes.Reader) error {

	var plaintext bytes.Buffer

	for {

		c1, err := er.ReadByte()

		if err != nil {
			plaintext.WriteTo(w.out)
			break
		}

		if c1 != 0x1b {

			plaintext.WriteByte(c1)
			continue
		}

		_, err = plaintext.WriteTo(w.out)

		if err != nil {
			return err
		}

		err = w.handleEscape(er)

		if err != nil {
			return err
		}

	}

	return nil
}

func (w *writer) handleEscape(er *bytes.Reader) error {

	c2, err := er.ReadByte()

	if err != nil {
		return err
	}

	switch c2 {

	case '>':
		return nil

	case ']':
		return w.handleTitle(er)

	case '7':
		return w.saveCursor()

	case '8':
		return w.restoreCursor()

	case 0x5b:
		cmd, params, err := w.parseCSI(er)

		if err != nil {
			return err
		}

		return w.executeCSI(cmd, params)

	default:
		return nil
	}
}

func (w *writer) parseCSI(er *bytes.Reader) (byte, string, error) {

	var buf bytes.Buffer
	var m byte

	w.rest.Reset()

	er.WriteTo(&w.rest)

	for i, c := range w.rest.Bytes() {

		if ('a' <= c && c <= 'z') ||
			('A' <= c && c <= 'Z') ||
			c == '@' {

			m = c

			params := buf.String()

			w.rest = *bytes.NewBuffer(
				w.rest.Bytes()[i+1:],
			)

			return m, params, nil
		}

		buf.WriteByte(c)
	}

	return 0, "", nil
}

func (w *writer) executeCSI(cmd byte, params string) error {

	switch cmd {

	case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'f':
		return w.handleCursor(cmd, params)

	case 'J', 'K', 'X':
		return w.handleErase(cmd, params)

	case 'm':
		return w.handleColor(params)

	case 'h', 'l':
		return w.handleMode(cmd, params)

	case 's', 'u':
		return w.handleSaveRestore(cmd)

	}

	return nil
}

// returns Atoi(s) unless s == "" in which case it returns def
func atoiWithDefault(s string, def int) (int, error) {
	if s == "" {
		return def, nil
	}
	return strconv.Atoi(s)
}
