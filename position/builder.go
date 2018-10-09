package position

import (
	"fmt"

	"github.com/rjo67/chess/bitset"
	"github.com/rjo67/chess/piece"
	"github.com/rjo67/chess/piece/colour"
	"github.com/rjo67/chess/square"
)

// Builder implements the builder pattern for a position. Create with NewBuilder
type Builder struct {
	pieces                         []map[piece.Piece]bitset.BitSet
	piecesInitialised              []map[piece.Piece]bool // whether the appropriate entry in the map has already been set
	activeColour                   colour.Colour
	castlingAvailabilityKingsSide  []bool
	castlingAvailabilityQueensSide []bool
	enpassantSquare                *square.Square
	halfmoveClock                  int
	fullmoveNbr                    int
}

// NewBuilder returns a builder for a position
func NewBuilder() *Builder {
	b := Builder{}
	b.pieces = make([]map[piece.Piece]bitset.BitSet, 2)
	b.piecesInitialised = make([]map[piece.Piece]bool, 2)
	b.castlingAvailabilityKingsSide = make([]bool, 2)
	b.castlingAvailabilityQueensSide = make([]bool, 2)
	// init the maps
	for _, col := range colour.AllColours {
		b.pieces[col] = make(map[piece.Piece]bitset.BitSet)
		b.piecesInitialised[col] = make(map[piece.Piece]bool)
	}
	return &b
}

// AddPiece adds the given piece bitset to the builder
func (b *Builder) AddPiece(col colour.Colour, pieceID piece.Piece, bs *bitset.BitSet) *Builder {
	if b.piecesInitialised[col][pieceID] {
		panic(fmt.Sprintf("AddPiece called multiple times for colour: %d, piece type: %d", col, pieceID))
	}
	b.pieces[col][pieceID] = *bs
	b.piecesInitialised[col][pieceID] = true
	return b
}

// ActiveColour sets the active colour of the position
func (b *Builder) ActiveColour(col colour.Colour) *Builder {
	b.activeColour = col
	return b
}

// CastlingAvailability sets the castling rights of the position
func (b *Builder) CastlingAvailability(col colour.Colour, kingsside bool, canCastle bool) *Builder {
	if kingsside {
		b.castlingAvailabilityKingsSide[col] = canCastle
	}
	b.castlingAvailabilityQueensSide[col] = canCastle
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
	posn := NewPosition(b.pieces[colour.White], b.pieces[colour.Black])
	posn.setEnpassantSquare(b.enpassantSquare)
	posn.setHalfmoveClock(b.halfmoveClock)
	posn.setFullmoveNbr(b.fullmoveNbr)
	for _, col := range colour.AllColours {
		posn.setCastlingAvailability(col, true, b.castlingAvailabilityKingsSide[col])
		posn.setCastlingAvailability(col, false, b.castlingAvailabilityQueensSide[col])
	}

	return posn
}
