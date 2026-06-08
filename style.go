package color

// Color represents a set of ANSI attributes that will be applied to the text.
type Color struct {
	attributes []Attribute
}

// New creates a new Color object with the provided attributes.
func New(value ...Attribute) *Color {
	color := &Color{
		attributes: make([]Attribute, 0),
	}
	color.Add(value...)
	return color
}

// Add appends new attributes to the color.
func (color *Color) Add(value ...Attribute) *Color {
	color.attributes = append(color.attributes, value...)
	return color
}
