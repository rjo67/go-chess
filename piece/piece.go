package piece

import (
	"github.com/rjo67/chess/bitset"
	"github.com/rjo67/chess/square"
)

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

// AllPieces to iterate over the piece types
var AllPieces = []Piece{PAWN, ROOK, KNIGHT, BISHOP, QUEEN, KING}

// RooksStartPosn contains the start position for the rooks
var RooksStartPosn = []bitset.BitSet{bitset.NewFromSquares(square.A1, square.H1), bitset.NewFromSquares(square.A8, square.H8)}

// KnightsStartPosn contains the start position for the knights
var KnightsStartPosn = []bitset.BitSet{bitset.NewFromSquares(square.B1, square.G1), bitset.NewFromSquares(square.B8, square.G8)}

// BishopsStartPosn contains the start position for the bishops
var BishopsStartPosn = []bitset.BitSet{bitset.NewFromSquares(square.C1, square.F1), bitset.NewFromSquares(square.C8, square.F8)}

// QueensStartPosn contains the start position for the queens
var QueensStartPosn = []bitset.BitSet{bitset.NewFromSquares(square.D1), bitset.NewFromSquares(square.D8)}

// KingsStartPosn contains the start position for the kings
var KingsStartPosn = []bitset.BitSet{bitset.NewFromSquares(square.E1), bitset.NewFromSquares(square.E8)}

// PawnsStartPosn contains the start position for the pawns
var PawnsStartPosn = []bitset.BitSet{bitset.NewFromSquares(square.A2, square.B2, square.C2, square.D2, square.E2, square.F2, square.G2, square.H2),
	bitset.NewFromSquares(square.A7, square.B7, square.C7, square.D7, square.E7, square.F7, square.G7, square.H7)}

var pieceMapping = [][]string{{"P", "R", "N", "B", "Q", "K"}, {"p", "r", "n", "b", "q", "k"}}

// ToString returns a letter describing the piece
func (p Piece) ToString(colour Colour) string {
	return pieceMapping[colour][p]
}
