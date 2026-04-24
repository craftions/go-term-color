//go:build windows && !appengine
// +build windows,!appengine

package colorable

import (
	"bytes"
	"os"
	"testing"
)

func TestAtoiWithDefault(t *testing.T) {
	val, err := atoiWithDefault("", 10)
	if err != nil || val != 10 {
		t.Errorf("Expected 10, got %d (err: %v)", val, err)
	}

	val, err = atoiWithDefault("5", 10)
	if err != nil || val != 5 {
		t.Errorf("Expected 5, got %d (err: %v)", val, err)
	}

	_, err = atoiWithDefault("invalid", 10)
	if err == nil {
		t.Errorf("Expected error for invalid int")
	}
}

func TestDoTitleSequence(t *testing.T) {
	// "0;TITLE\x07" -> valid
	seq := []byte{'0', ';', 'T', 'I', 'T', 'L', 'E', 0x07}
	err := doTitleSequence(bytes.NewReader(seq))
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Invalid sequence format
	seq = []byte{'1'} // not 0 or 2
	err = doTitleSequence(bytes.NewReader(seq))
	if err != nil {
		t.Errorf("Expected no error for unhandled case, got %v", err)
	}

	seq = []byte{'0', 'X'} // not ';'
	err = doTitleSequence(bytes.NewReader(seq))
	if err != nil {
		t.Errorf("Expected no error for unhandled case, got %v", err)
	}
}

func TestWriterBasic(t *testing.T) {
	var out bytes.Buffer
	w := &writer{
		out: &out,
	}

	data := []byte("hello world")
	n, err := w.Write(data)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if n != len(data) {
		t.Errorf("Expected %d, got %d", len(data), n)
	}
}

func TestWriterANSI(t *testing.T) {
	var out bytes.Buffer
	w := &writer{
		out: &out,
	}

	sequences := []string{
		"\x1b[31m", // red
		"\x1b[38;5;12m", // 256 color
		"\x1b[38;2;255;0;0m", // rgb color
		"\x1b[41m", // bg red
		"\x1b[48;5;12m", // bg 256
		"\x1b[48;2;255;0;0m", // bg rgb
		"\x1b[0m", // reset
		"\x1b[1m", // bold
		"\x1b[4m", // underline
		"\x1b[7m", // reverse
		"\x1b[2A", // cursor up
		"\x1b[2B", // cursor down
		"\x1b[2C", // cursor forward
		"\x1b[2D", // cursor back
		"\x1b[2E", // cursor next line
		"\x1b[2F", // cursor prev line
		"\x1b[2G", // cursor h abs
		"\x1b[2;2H", // cursor pos
		"\x1b[2J", // erase in display
		"\x1b[2K", // erase in line
		"\x1b[2X", // erase chars
		"\x1b[s", // save cursor pos
		"\x1b[u", // restore cursor pos
		"\x1b[?25h", // show cursor
		"\x1b[?25l", // hide cursor
		"\x1b]0;Title\x07", // title
		"\x1b[90m", // bright fg
		"\x1b[100m", // bright bg
		"\x1b[38;5m", // invalid
		"\x1b[38;2m", // invalid
		"\x1b[48;5m", // invalid
		"\x1b[48;2m", // invalid
		"\x1b[30m\x1b[31m\x1b[32m\x1b[33m\x1b[34m\x1b[35m\x1b[36m\x1b[37m", // colors
		"\x1b[40m\x1b[41m\x1b[42m\x1b[43m\x1b[44m\x1b[45m\x1b[46m\x1b[47m", // bg colors
		"\x1b[91m\x1b[92m\x1b[93m\x1b[94m\x1b[95m\x1b[96m\x1b[97m", // bright colors
		"\x1b[101m\x1b[102m\x1b[103m\x1b[104m\x1b[105m\x1b[106m\x1b[107m", // bright bg colors
		"\x1b[22m\x1b[24m\x1b[27m\x1b[39m\x1b[49m", // resets
	}
	for _, seq := range sequences {
		w.Write([]byte(seq))
	}
}

func TestEnableColorsStdout(t *testing.T) {
	var enabled bool
	restore := EnableColorsStdout(&enabled)
	if !enabled {
		t.Errorf("Expected enabled to be true")
	}
	restore()

	restore = EnableColorsStdout(nil)
	restore()
}

func TestNewColorableWindowsConsole(t *testing.T) {
	// Try opening the real console to force full evaluation
	f, err := os.OpenFile("CONOUT$", os.O_RDWR, 0644)
	if err == nil {
		defer f.Close()
		w := NewColorable(f)
		if w == nil {
			t.Errorf("Expected valid writer, got nil")
		}
	}

	// Test passing a regular file (non-terminal)
	tmp, err := os.CreateTemp("", "test")
	if err == nil {
		defer os.Remove(tmp.Name())
		defer tmp.Close()

		w := NewColorable(tmp)
		// Should return the same file because it is not a terminal
		if w != tmp {
			t.Errorf("Expected same file for non-terminal, got %T", w)
		}
	}
}

