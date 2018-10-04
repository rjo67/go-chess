package position

import (
	"strings"

	"github.com/rjo67/chess/bitset"
	"github.com/rjo67/chess/piece"
)

// Position represents a chess position
type Position struct {
	// the pieces, keyed by Colour and Piece
	pieces [][]bitset.BitSet
}

// NewPosition creates a new position
// The bitset arrays are in the order as given by the piece constants
func NewPosition(whitePieces []bitset.BitSet, blackPieces []bitset.BitSet) Position {
	var posn Position
	posn.pieces = make([][]bitset.BitSet, 2)
	posn.pieces[piece.WHITE] = whitePieces
	posn.pieces[piece.BLACK] = blackPieces

	return posn
}

// StartPosition creates a new start position
func StartPosition() Position {
	return NewPosition([]bitset.BitSet{piece.PawnsStartPosn[piece.WHITE], piece.RooksStartPosn[piece.WHITE], piece.KnightsStartPosn[piece.WHITE],
		piece.BishopsStartPosn[piece.WHITE], piece.QueensStartPosn[piece.WHITE], piece.KingsStartPosn[piece.WHITE]},
		[]bitset.BitSet{piece.PawnsStartPosn[piece.BLACK], piece.RooksStartPosn[piece.BLACK], piece.KnightsStartPosn[piece.BLACK],
			piece.BishopsStartPosn[piece.BLACK], piece.QueensStartPosn[piece.BLACK], piece.KingsStartPosn[piece.BLACK]})
}

// ToString delivers a string representation of the current position
func (p Position) ToString() string {
	var squares [64]string

	// populate squares with the bitset contents
	for _, colour := range piece.AllColours {
		for _, pieceType := range piece.AllPieces {
			positions := p.pieces[colour][pieceType].SetBits()
			for _, sq := range positions {
				squares[sq-1] = pieceType.ToString(colour)
			}
		}
	}

	var sb strings.Builder
	sb.WriteString("+--------+\n")
	for i := 63; i >= 0; i-- {
		if i%8 == 7 {
			sb.WriteString("|")
		}
		if squares[i] == "" {
			sb.WriteString(".")
		} else {
			sb.WriteString(squares[i])
		}
		if i%8 == 0 {
			sb.WriteString("|\n")
		}
	}
	sb.WriteString("+--------+\n")

	return sb.String()
}
