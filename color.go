package color

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/craftions/go-term-check/terminal"
)

// NoColor determines if the use of colors should be omitted.
// It will be true if the output is redirected to a file, or if environment
// variables like NO_COLOR were specified.
var NoColor bool

// CheckIfTerminal abstracts and wraps the complete inter-OS terminal check,
// replacing conventional external terminal dependencies.
func CheckIfTerminal(fd uintptr) bool {
	// 1. Standard cross-platform base verification
	if terminal.IsTerminal(fd) {
		return true
	}
	// 2. Vital extension to tolerate emulations on Windows (PTY/MSYS - Git Bash / Mintty)
	if runtime.GOOS == "windows" {
		if terminal.IsCygwinTerminal(fd) {
			return true
		}
	}
	return false
}

func init() {
	setupNoColor()
}

// setupNoColor extracts the initialization logic to facilitate testing.
func setupNoColor() {
	if os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb" {
		NoColor = true
		return
	}

	// Using our own integrated check to automatically enable/disable ANSI colors
	NoColor = !CheckIfTerminal(os.Stdout.Fd())
}

// Attribute represents an ANSI escape visual attribute.
type Attribute int

// Base attributes
const (
	Reset Attribute = iota
	Bold
	Faint
	Italic
	Underline
	BlinkSlow
	BlinkRapid
	ReverseVideo
	Concealed
	CrossedOut
)

// Foreground Colors
const (
	FgBlack Attribute = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
)

// Background Colors
const (
	BgBlack Attribute = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
)

// Color represents a set of ANSI attributes that will be applied to the text.
type Color struct {
	params []Attribute
}

// New creates a new Color object with the provided attributes.
func New(value ...Attribute) *Color {
	color := &Color{
		params: make([]Attribute, 0),
	}
	color.Add(value...)
	return color
}

// Add appends new attributes to the color.
func (color *Color) Add(value ...Attribute) *Color {
	color.params = append(color.params, value...)
	return color
}

// format concatenates the numerical ANSI parameters
func (color *Color) format() string {
	format := make([]string, len(color.params))
	for i, v := range color.params {
		format[i] = fmt.Sprintf("%d", int(v))
	}
	return strings.Join(format, ";")
}

// wrap applies the ANSI sequences to the text, if colors are enabled.
func (color *Color) wrap(format string, a ...interface{}) string {
	text := fmt.Sprintf(format, a...)
	if NoColor {
		return text
	}
	return fmt.Sprintf("\x1b[%sm%s\x1b[0m", color.format(), text)
}

// Print prints the formatted text to stdout
func (color *Color) Print(a ...interface{}) (n int, err error) {
	text := fmt.Sprint(a...)
	if NoColor {
		return fmt.Print(text)
	}
	return fmt.Printf("\x1b[%sm%s\x1b[0m", color.format(), text)
}

// Println prints the formatted text to stdout, followed by a newline
func (color *Color) Println(a ...interface{}) (n int, err error) {
	text := fmt.Sprint(a...)
	if NoColor {
		return fmt.Println(text)
	}
	return fmt.Printf("\x1b[%sm%s\x1b[0m\n", color.format(), text)
}

// Printf prints the formatted text using format to stdout
func (color *Color) Printf(format string, a ...interface{}) (n int, err error) {
	text := fmt.Sprintf(format, a...)
	if NoColor {
		return fmt.Print(text)
	}
	return fmt.Printf("\x1b[%sm%s\x1b[0m", color.format(), text)
}

// --- Simplified helper functions (inspired by fatih/color) ---

// RedString returns text wrapped in red ANSI.
func RedString(format string, a ...interface{}) string { return New(FgRed).wrap(format, a...) }

// GreenString returns text wrapped in green ANSI.
func GreenString(format string, a ...interface{}) string { return New(FgGreen).wrap(format, a...) }

// YellowString returns text wrapped in yellow ANSI.
func YellowString(format string, a ...interface{}) string { return New(FgYellow).wrap(format, a...) }

// BlueString returns text wrapped in blue ANSI.
func BlueString(format string, a ...interface{}) string { return New(FgBlue).wrap(format, a...) }

// MagentaString returns text wrapped in magenta ANSI.
func MagentaString(format string, a ...interface{}) string { return New(FgMagenta).wrap(format, a...) }

// CyanString returns text wrapped in cyan ANSI.
func CyanString(format string, a ...interface{}) string { return New(FgCyan).wrap(format, a...) }

// WhiteString returns text wrapped in white ANSI.
func WhiteString(format string, a ...interface{}) string { return New(FgWhite).wrap(format, a...) }

// Formatters that print adding a newline at the end if it's not string-wrapping (behavior of `color.Red("...")`)

// Red prints text in red color. Adds a newline at the end.
func Red(format string, a ...interface{}) { New(FgRed).Printf(format+"\n", a...) }

// Green prints text in green color. Adds a newline at the end.
func Green(format string, a ...interface{}) { New(FgGreen).Printf(format+"\n", a...) }

// Yellow prints text in yellow color. Adds a newline at the end.
func Yellow(format string, a ...interface{}) { New(FgYellow).Printf(format+"\n", a...) }

// Blue prints text in blue color. Adds a newline at the end.
func Blue(format string, a ...interface{}) { New(FgBlue).Printf(format+"\n", a...) }

// Magenta prints text in magenta color. Adds a newline at the end.
func Magenta(format string, a ...interface{}) { New(FgMagenta).Printf(format+"\n", a...) }

// Cyan prints text in cyan color. Adds a newline at the end.
func Cyan(format string, a ...interface{}) { New(FgCyan).Printf(format+"\n", a...) }

// White prints text in white color. Adds a newline at the end.
func White(format string, a ...interface{}) { New(FgWhite).Printf(format+"\n", a...) }
