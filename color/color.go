package color

import (
	"go-term-color/internal/colorable"
	"io"
)

// Attribute defines a single SGR Code
type Attribute int

// Color defines a custom color object which is defined by SGR parameters.
type Color struct {
	params  []Attribute
	noColor *bool
}

// Output defines the standard output of the print functions.
var Output io.Writer = colorable.NewColorableStdout()

// Error defines a color supporting writer for os.Stderr.
var Error io.Writer = colorable.NewColorableStderr()

// NoColor defines if the output is colorized or not.
var NoColor bool
