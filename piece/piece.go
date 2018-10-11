package piece

import (
	"github.com/rjo67/chess/bitset"
	"github.com/rjo67/chess/piece/colour"
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

// PromotedPawnPieceCandidates specifies the pieces which a pawn can promote to
var PromotedPawnPieceCandidates = []Piece{QUEEN, ROOK, BISHOP, KNIGHT}

// StartPosn stores the start positions of all pieces, keyed by colour and piece type (ordering as given by AllColours and AllPieces)
var StartPosn = [][]bitset.BitSet{
	// WHITE
	{
		bitset.NewFromSquares(square.A2, square.B2, square.C2, square.D2, square.E2, square.F2, square.G2, square.H2),
		bitset.NewFromSquares(square.A1, square.H1),
		bitset.NewFromSquares(square.B1, square.G1),
		bitset.NewFromSquares(square.C1, square.F1),
		bitset.NewFromSquares(square.D1),
		bitset.NewFromSquares(square.E1),
	},
	// BLACK
	{
		bitset.NewFromSquares(square.A7, square.B7, square.C7, square.D7, square.E7, square.F7, square.G7, square.H7),
		bitset.NewFromSquares(square.A8, square.H8),
		bitset.NewFromSquares(square.B8, square.G8),
		bitset.NewFromSquares(square.C8, square.F8),
		bitset.NewFromSquares(square.D8),
		bitset.NewFromSquares(square.E8),
	},
}

var (
	// WhitePawnString is the string identifying a white pawn
	WhitePawnString = "P"
	// WhiteRookString is the string identifying a white rook
	WhiteRookString = "R"
	// WhiteKnightString is the string identifying a white knight
	WhiteKnightString = "N"
	// WhiteBishopString is the string identifying a white bishop
	WhiteBishopString = "B"
	// WhiteQueenString is the string identifying a white queen
	WhiteQueenString = "Q"
	// WhiteKingString is the string identifying a white king
	WhiteKingString = "K"
	// BlackPawnString is the string identifying a black pawn
	BlackPawnString = "p"
	// BlackRookString is the string identifying a black rook
	BlackRookString = "r"
	// BlackKnightString is the string identifying a black knight
	BlackKnightString = "n"
	// BlackBishopString is the string identifying a black bishop
	BlackBishopString = "b"
	// BlackQueenString is the string identifying a black queen
	BlackQueenString = "q"
	// BlackKingString is the string identifying a black king
	BlackKingString = "k"
)

// PieceMapping contains the strings identifying the pieces, keyed by colour
var PieceMapping = [][]string{{WhitePawnString, WhiteRookString, WhiteKnightString, WhiteBishopString, WhiteQueenString, WhiteKingString},
	{BlackPawnString, BlackRookString, BlackKnightString, BlackBishopString, BlackQueenString, BlackKingString}}

// StringToPiece is a mapping between the piece strings and the piece types
var StringToPiece = map[string]Piece{
	WhitePawnString:   PAWN,
	WhiteRookString:   ROOK,
	WhiteKnightString: KNIGHT,
	WhiteBishopString: BISHOP,
	WhiteQueenString:  QUEEN,
	WhiteKingString:   KING,
	BlackPawnString:   PAWN,
	BlackRookString:   ROOK,
	BlackKnightString: KNIGHT,
	BlackBishopString: BISHOP,
	BlackQueenString:  QUEEN,
	BlackKingString:   KING,
}

// String returns a letter describing the piece
func (p Piece) String(colour colour.Colour) string {
	return PieceMapping[colour][p]
}

// FromString returns a piece from the given letter.FromString
// No distinction between black or white pieces (i.e. "P" or "p" will both return PAWN)
func FromString(colour colour.Colour, str string) Piece {
	return StringToPiece[str]
}
