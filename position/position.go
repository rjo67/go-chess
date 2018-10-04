package position

import (
	"github.com/rjo67/chess/bitset"
	"github.com/rjo67/chess/piece"
)

// Position represents a chess position
type Position struct {
	// the pieces, keyed by Colour and Piece
	pieces [][]bitset.BitSet
}

// NewPosition creates a new position
func NewPosition(whitePieces []bitset.BitSet, blackPieces []bitset.BitSet) Position {
	var posn Position
	posn.pieces[piece.WHITE] = whitePieces
	posn.pieces[piece.BLACK] = blackPieces

	return posn
}
