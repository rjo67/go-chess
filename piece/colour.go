package piece

// Colour represents the colours
type Colour uint32

// the colours
const (
	WHITE Colour = iota
	BLACK
)

// Other returns the other colour
func (c Colour) Other() Colour {
	if c == WHITE {
		return BLACK
	}
	return WHITE
}

// AllColours to iterate over the colours
var AllColours = []Colour{WHITE, BLACK}

var colourMapping = []string{"W", "B"}

// ToString returns a letter describing the colour
func (c Colour) ToString() string {
	return colourMapping[c]
}
