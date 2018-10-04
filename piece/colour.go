package piece

// Colour represents the colours
type Colour uint32

// the colours
const (
	WHITE Colour = iota
	BLACK
)

// AllColours to iterate over the colours
var AllColours = []Colour{WHITE, BLACK}

var colourMapping = []string{"W", "B"}

// ToString returns a letter describing the colour
func (col Colour) ToString() string {
	return colourMapping[col]
}
