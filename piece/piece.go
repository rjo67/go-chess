package piece

// Piece represents a chess piece
type Piece uint32

// the pieces
const (
	PAWN Piece = iota
	ROOK
	KNIGHT
	BISHOP
	QUEEN
	KING
)

var pieceMapping = []string{"P", "R", "K", "B", "Q", "K"}

// ToString returns a letter describing the piece
func (p Piece) ToString() string {
	return pieceMapping[p]
}
