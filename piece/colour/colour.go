package colour

// Colour represents the colours
type Colour uint32

// the colours
const (
	White Colour = iota
	Black
	AnyColour // any colour
)

// Other returns the other colour
func (c Colour) Other() Colour {
	if c == White {
		return Black
	}
	return White
}

// AllColours to iterate over the colours
var AllColours = []Colour{White, Black}

var colourMapping = []string{"W", "B"}

// String returns a letter describing the colour
func (c Colour) String() string {
	return colourMapping[c]
}
