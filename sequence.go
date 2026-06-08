package color

import (
	"strconv"
	"strings"
)

// sequence concatenates the numerical ANSI attributes
func (color *Color) sequence() string {
	seq := make([]string, len(color.attributes))
	for i, v := range color.attributes {
		seq[i] = strconv.Itoa(int(v))
	}
	return strings.Join(seq, ";")
}

// apply applies the ANSI sequences to the text, if colors are enabled.
func (color *Color) apply(text string) string {
	if NoColor {
		return text
	}
	return "\x1b[" + color.sequence() + "m" + text + "\x1b[0m"
}
