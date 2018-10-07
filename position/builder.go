package position

import (
	"fmt"

	"github.com/rjo67/chess/bitset"
	"github.com/rjo67/chess/piece"
	"github.com/rjo67/chess/square"
)

// Builder implements the builder pattern for a position. Create with NewBuilder
type Builder struct {
	pieces               []map[piece.Piece]bitset.BitSet
	piecesInitialised    []map[piece.Piece]bool // whether the appropriate entry in the map has already been set
	activeColour         piece.Colour
	castlingAvailability string
	enpassantSquare      *square.Square
	halfmoveClock        int
	fullmoveNbr          int
}

// NewBuilder returns a builder for a position
func NewBuilder() *Builder {
	b := Builder{}
	b.pieces = make([]map[piece.Piece]bitset.BitSet, 2)
	b.piecesInitialised = make([]map[piece.Piece]bool, 2)
	// init the maps
	for _, colour := range piece.AllColours {
		b.pieces[colour] = make(map[piece.Piece]bitset.BitSet)
		b.piecesInitialised[colour] = make(map[piece.Piece]bool)
	}
	return &b
}

// AddPiece adds the given piece bitset to the builder
func (b *Builder) AddPiece(colour piece.Colour, pieceID piece.Piece, bs *bitset.BitSet) *Builder {
	if b.piecesInitialised[colour][pieceID] {
		panic(fmt.Sprintf("AddPiece called multiple times for colour: %d, piece type: %d", colour, pieceID))
	}
	b.pieces[colour][pieceID] = *bs
	b.piecesInitialised[colour][pieceID] = true
	return b
}

// ActiveColour sets the active colour of the position
func (b *Builder) ActiveColour(colour piece.Colour) *Builder {
	b.activeColour = colour
	return b
}

// CastlingAvailability sets the castling rights of the position
func (b *Builder) CastlingAvailability(castlingAvailability string) *Builder {
	b.castlingAvailability = castlingAvailability
	return b
}

// EnpassantSquare sets the enpassant square of the position
func (b *Builder) EnpassantSquare(enpassantSquare *square.Square) *Builder {
	b.enpassantSquare = enpassantSquare
	return b
}

// HalfmoveClock sets the halfmove clock of the position
func (b *Builder) HalfmoveClock(clock int) *Builder {
	b.halfmoveClock = clock
	return b
}

// FullmoveNbr sets the fullmove nbr of the position
func (b *Builder) FullmoveNbr(moveNbr int) *Builder {
	b.fullmoveNbr = moveNbr
	return b
}

// Build builds a position object
func (b *Builder) Build() Position {
	posn := NewPosition(b.pieces[piece.WHITE], b.pieces[piece.BLACK])
	posn.setEnpassantSquare(b.enpassantSquare)
	posn.setHalfmoveClock(b.halfmoveClock)
	posn.setFullmoveNbr(b.fullmoveNbr)
	posn.setCastlingAvailability(b.castlingAvailability)

	return posn
}
