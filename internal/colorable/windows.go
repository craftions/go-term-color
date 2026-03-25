//This file defines the writer that allows the interpretation of ANSI sequences
// (colors, cursor, clear screen, etc.) in the Windows console.
//go:build windows && !appengine
// +build windows,!appengine

package colorable

import (
	"bytes"
	"io"
	"math"
	"os"
	"sync"
	"unsafe"

	syscall "golang.org/x/sys/windows"

	gotermcheck "github.com/mattn/go-isatty"
)

const (
	foregroundBlue      = 0x1
	foregroundGreen     = 0x2
	foregroundRed       = 0x4
	foregroundIntensity = 0x8
	foregroundMask      = (foregroundRed | foregroundBlue | foregroundGreen | foregroundIntensity)
	backgroundBlue      = 0x10
	backgroundGreen     = 0x20
	backgroundRed       = 0x40
	backgroundIntensity = 0x80
	backgroundMask      = (backgroundRed | backgroundBlue | backgroundGreen | backgroundIntensity)
	commonLvbUnderscore = 0x8000

	cENABLE_VIRTUAL_TERMINAL_PROCESSING = 0x4
)

const (
	genericRead  = 0x80000000
	genericWrite = 0x40000000
)

const (
	consoleTextmodeBuffer = 0x1
)

type wchar uint16
type short int16
type dword uint32
type word uint16

type coord struct {
	x short
	y short
}

type smallRect struct {
	left   short
	top    short
	right  short
	bottom short
}

type consoleScreenBufferInfo struct {
	size              coord
	cursorPosition    coord
	attributes        word
	window            smallRect
	maximumWindowSize coord
}

type consoleCursorInfo struct {
	size    dword
	visible int32
}

var (
	kernel32                       = syscall.NewLazySystemDLL("kernel32.dll")
	procGetConsoleScreenBufferInfo = kernel32.NewProc("GetConsoleScreenBufferInfo")
	procSetConsoleTextAttribute    = kernel32.NewProc("SetConsoleTextAttribute")
	procSetConsoleCursorPosition   = kernel32.NewProc("SetConsoleCursorPosition")
	procFillConsoleOutputCharacter = kernel32.NewProc("FillConsoleOutputCharacterW")
	procFillConsoleOutputAttribute = kernel32.NewProc("FillConsoleOutputAttribute")
	procGetConsoleCursorInfo       = kernel32.NewProc("GetConsoleCursorInfo")
	procSetConsoleCursorInfo       = kernel32.NewProc("SetConsoleCursorInfo")
	procSetConsoleTitle            = kernel32.NewProc("SetConsoleTitleW")
	procGetConsoleMode             = kernel32.NewProc("GetConsoleMode")
	procSetConsoleMode             = kernel32.NewProc("SetConsoleMode")
	procCreateConsoleScreenBuffer  = kernel32.NewProc("CreateConsoleScreenBuffer")
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

	if gotermcheck.IsTerminal(file.Fd()) {
		var mode uint32
		if r, _, _ := procGetConsoleMode.Call(file.Fd(), uintptr(unsafe.Pointer(&mode))); r != 0 && mode&cENABLE_VIRTUAL_TERMINAL_PROCESSING != 0 {
			return file
		}
		var csbi consoleScreenBufferInfo
		handle := syscall.Handle(file.Fd())
		procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))
		return &writer{out: file, handle: handle, oldattr: csbi.attributes, oldpos: coord{0, 0}}
	}
	return file
}

// NewColorableStdout returns new instance of writer which handles escape sequence for stdout.
func NewColorableStdout() io.Writer {
	return NewColorable(os.Stdout)
}

// NewColorableStderr returns new instance of writer which handles escape sequence for stderr.
func NewColorableStderr() io.Writer {
	return NewColorable(os.Stderr)
}

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

func (c consoleColor) foregroundAttr() (attr word) {
	if c.red {
		attr |= foregroundRed
	}
	if c.green {
		attr |= foregroundGreen
	}
	if c.blue {
		attr |= foregroundBlue
	}
	if c.intensity {
		attr |= foregroundIntensity
	}
	return
}

func (c consoleColor) backgroundAttr() (attr word) {
	if c.red {
		attr |= backgroundRed
	}
	if c.green {
		attr |= backgroundGreen
	}
	if c.blue {
		attr |= backgroundBlue
	}
	if c.intensity {
		attr |= backgroundIntensity
	}
	return
}

func (a hsv) dist(b hsv) float32 {
	dh := a.h - b.h
	switch {
	case dh > 0.5:
		dh = 1 - dh
	case dh < -0.5:
		dh = -1 - dh
	}
	ds := a.s - b.s
	dv := a.v - b.v
	return float32(math.Sqrt(float64(dh*dh + ds*ds + dv*dv)))
}
