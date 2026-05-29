//go:build windows

package colorable

import (
	"os"
	"testing"

	"golang.org/x/sys/windows"
)

func TestColorable_PanicNil(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	NewColorable(nil)
}

func TestColorable_NotTerminal(t *testing.T) {
	r, w, _ := os.Pipe()
	defer r.Close()
	defer w.Close()

	_ = NewColorable(w)
}

func TestColorable_IsTerminal(t *testing.T) {
	oldGetConsoleMode := getConsoleMode
	oldSetConsoleMode := setConsoleMode
	
	getConsoleMode = func(handle windows.Handle, mode *uint32) error {
		*mode = 0
		return nil
	}
	setConsoleMode = func(handle windows.Handle, mode uint32) error {
		return nil
	}
	
	defer func() {
		getConsoleMode = oldGetConsoleMode
		setConsoleMode = oldSetConsoleMode
	}()

	_ = NewColorable(os.Stdout)
}
