package color

import (
	"fmt"
	"io"
	"os"
)

// Sprint formats using the default formats for its operands and returns the resulting string.
func (color *Color) Sprint(a ...any) string {
	return color.render(fmt.Sprint(a...))
}

// Sprintf formats according to a format specifier and returns the resulting string.
func (color *Color) Sprintf(format string, a ...any) string {
	return color.render(fmt.Sprintf(format, a...))
}

// Fprint formats using the default formats for its operands and writes to w.
func (color *Color) Fprint(w io.Writer, a ...any) (n int, err error) {
	return fmt.Fprint(w, color.Sprint(a...))
}

// Fprintf formats according to a format specifier and writes to w.
func (color *Color) Fprintf(w io.Writer, format string, a ...any) (n int, err error) {
	return fmt.Fprint(w, color.Sprintf(format, a...))
}

// Fprintln formats using the default formats for its operands and writes to w.
func (color *Color) Fprintln(w io.Writer, a ...any) (n int, err error) {
	return fmt.Fprintln(w, color.Sprint(a...))
}

// Print prints the formatted text to stdout
func (color *Color) Print(a ...any) (n int, err error) {
	return color.Fprint(os.Stdout, a...)
}

// Println prints the formatted text to stdout, followed by a newline
func (color *Color) Println(a ...any) (n int, err error) {
	return color.Fprintln(os.Stdout, a...)
}

// Printf prints the formatted text using format to stdout
func (color *Color) Printf(format string, a ...any) (n int, err error) {
	return color.Fprintf(os.Stdout, format, a...)
}

// --- Simplified helper functions (inspired by fatih/color) ---

// RedString returns text wrapped in red ANSI.
func RedString(format string, a ...any) string { return New(FgRed).Sprintf(format, a...) }

// GreenString returns text wrapped in green ANSI.
func GreenString(format string, a ...any) string { return New(FgGreen).Sprintf(format, a...) }

// YellowString returns text wrapped in yellow ANSI.
func YellowString(format string, a ...any) string { return New(FgYellow).Sprintf(format, a...) }

// BlueString returns text wrapped in blue ANSI.
func BlueString(format string, a ...any) string { return New(FgBlue).Sprintf(format, a...) }

// MagentaString returns text wrapped in magenta ANSI.
func MagentaString(format string, a ...any) string { return New(FgMagenta).Sprintf(format, a...) }

// CyanString returns text wrapped in cyan ANSI.
func CyanString(format string, a ...any) string { return New(FgCyan).Sprintf(format, a...) }

// WhiteString returns text wrapped in white ANSI.
func WhiteString(format string, a ...any) string { return New(FgWhite).Sprintf(format, a...) }

// Formatters that print adding a newline at the end if it's not string-wrapping (behavior of `color.Red("...")`)

// Red prints text in red color. Adds a newline at the end.
func Red(format string, a ...any) { New(FgRed).Fprintf(os.Stdout, format+"\n", a...) }

// Green prints text in green color. Adds a newline at the end.
func Green(format string, a ...any) { New(FgGreen).Fprintf(os.Stdout, format+"\n", a...) }

// Yellow prints text in yellow color. Adds a newline at the end.
func Yellow(format string, a ...any) { New(FgYellow).Fprintf(os.Stdout, format+"\n", a...) }

// Blue prints text in blue color. Adds a newline at the end.
func Blue(format string, a ...any) { New(FgBlue).Fprintf(os.Stdout, format+"\n", a...) }

// Magenta prints text in magenta color. Adds a newline at the end.
func Magenta(format string, a ...any) { New(FgMagenta).Fprintf(os.Stdout, format+"\n", a...) }

// Cyan prints text in cyan color. Adds a newline at the end.
func Cyan(format string, a ...any) { New(FgCyan).Fprintf(os.Stdout, format+"\n", a...) }

// White prints text in white color. Adds a newline at the end.
func White(format string, a ...any) { New(FgWhite).Fprintf(os.Stdout, format+"\n", a...) }
