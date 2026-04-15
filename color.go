package color

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/craftions/go-term-check/terminal"
)

// NoColor determina si se debe omitir el uso de colores.
// Será verdadero si la salida se redirige a un archivo, o si se especificaron
// variables de entorno como NO_COLOR.
var NoColor bool

// CheckIfTerminal abstrae y envuelve la comprobación completa inter-OS,
// reemplazando las dependencias externas convencionales de terminales.
func CheckIfTerminal(fd uintptr) bool {
	// 1. Verificación base estándar cruzada
	if terminal.IsTerminal(fd) {
		return true
	}
	// 2. Extensión vital para tolerar emulaciones sobre Windows (PTY/MSYS - Git Bash / Mintty)
	if runtime.GOOS == "windows" {
		if terminal.IsCygwinTerminal(fd) {
			return true
		}
	}
	return false
}

func init() {
	if os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb" {
		NoColor = true
		return
	}

	// Usando nuestro propio check integrado para activar/desactivar colores ANSI automáticamente
	NoColor = !CheckIfTerminal(os.Stdout.Fd())
}

// Attribute representa un atributo visual de escape ANSI.
type Attribute int

// Atributos base
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

// Colores Frontales (Foreground)
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

// Colores de Fondo (Background)
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

// Color representa un cojunto de atributos ANSI que se aplicarán al texto.
type Color struct {
	params []Attribute
}

// New crea un nuevo objeto Color con los atributos proporcionados.
func New(value ...Attribute) *Color {
	c := &Color{
		params: make([]Attribute, 0),
	}
	c.Add(value...)
	return c
}

// Add añade nuevos atributos al color.
func (c *Color) Add(value ...Attribute) *Color {
	c.params = append(c.params, value...)
	return c
}

// format concatena los parámetros numéricos de ANSI
func (c *Color) format() string {
	format := make([]string, len(c.params))
	for i, v := range c.params {
		format[i] = fmt.Sprintf("%d", int(v))
	}
	return strings.Join(format, ";")
}

// wrap aplica las secuencias ANSI al texto, si los colores están activados.
func (c *Color) wrap(format string, a ...interface{}) string {
	text := fmt.Sprintf(format, a...)
	if NoColor {
		return text
	}
	return fmt.Sprintf("\x1b[%sm%s\x1b[0m", c.format(), text)
}

// Print imprime el texto formateado en stdout
func (c *Color) Print(a ...interface{}) (n int, err error) {
	text := fmt.Sprint(a...)
	if NoColor {
		return fmt.Print(text)
	}
	return fmt.Printf("\x1b[%sm%s\x1b[0m", c.format(), text)
}

// Println imprime el texto formateado en stdout, seguido de un salto de línea
func (c *Color) Println(a ...interface{}) (n int, err error) {
	text := fmt.Sprint(a...)
	if NoColor {
		return fmt.Println(text)
	}
	return fmt.Printf("\x1b[%sm%s\x1b[0m\n", c.format(), text)
}

// Printf imprime el texto formateado usando formato en stdout
func (c *Color) Printf(format string, a ...interface{}) (n int, err error) {
	text := fmt.Sprintf(format, a...)
	if NoColor {
		return fmt.Print(text)
	}
	return fmt.Printf("\x1b[%sm%s\x1b[0m", c.format(), text)
}

// --- Funciones helper simplificadas (inspiradas en fatih/color) ---

// RedString devuelve texto envuelto en ANSI rojo.
func RedString(format string, a ...interface{}) string { return New(FgRed).wrap(format, a...) }

// GreenString devuelve texto envuelto en ANSI verde.
func GreenString(format string, a ...interface{}) string { return New(FgGreen).wrap(format, a...) }

// YellowString devuelve texto envuelto en ANSI amarillo.
func YellowString(format string, a ...interface{}) string { return New(FgYellow).wrap(format, a...) }

// BlueString devuelve texto envuelto en ANSI azul.
func BlueString(format string, a ...interface{}) string { return New(FgBlue).wrap(format, a...) }

// MagentaString devuelve texto envuelto en ANSI magenta.
func MagentaString(format string, a ...interface{}) string { return New(FgMagenta).wrap(format, a...) }

// CyanString devuelve texto envuelto en ANSI cyan.
func CyanString(format string, a ...interface{}) string { return New(FgCyan).wrap(format, a...) }

// WhiteString devuelve texto envuelto en ANSI blanco.
func WhiteString(format string, a ...interface{}) string { return New(FgWhite).wrap(format, a...) }

// Formateadores que imprimen añadiendo un salto de linea al final si no es string-wrapping (comportamiento de `color.Red("...")`)

// Red imprime texto en color rojo. Añade salto de línea al final.
func Red(format string, a ...interface{}) { New(FgRed).Printf(format+"\n", a...) }

// Green imprime texto en color verde. Añade salto de línea al final.
func Green(format string, a ...interface{}) { New(FgGreen).Printf(format+"\n", a...) }

// Yellow imprime texto en color amarillo. Añade salto de línea al final.
func Yellow(format string, a ...interface{}) { New(FgYellow).Printf(format+"\n", a...) }

// Blue imprime texto en color azul. Añade salto de línea al final.
func Blue(format string, a ...interface{}) { New(FgBlue).Printf(format+"\n", a...) }

// Magenta imprime texto en color magenta. Añade salto de línea al final.
func Magenta(format string, a ...interface{}) { New(FgMagenta).Printf(format+"\n", a...) }

// Cyan imprime texto en color cyan. Añade salto de línea al final.
func Cyan(format string, a ...interface{}) { New(FgCyan).Printf(format+"\n", a...) }

// White imprime texto en color blanco. Añade salto de línea al final.
func White(format string, a ...interface{}) { New(FgWhite).Printf(format+"\n", a...) }
