package position

import (
	"strings"

	"github.com/rjo67/chess/square"

	"github.com/rjo67/chess/bitset"
	"github.com/rjo67/chess/piece"
)

// Position represents a chess position
type Position struct {
	pieces               []map[piece.Piece]bitset.BitSet // array of map of piece bitsets, array-dim = colour
	activeColour         piece.Colour                    // whose move
	castlingAvailability string
	enpassantSquare      *square.Square
	halfmoveClock        int
	fullmoveNbr          int
}

// NewPosition creates a new position
// The bitset arrays are in the order as given by the piece constants
func NewPosition(whitePieces, blackPieces map[piece.Piece]bitset.BitSet) Position {
	var p Position
	p.pieces = make([]map[piece.Piece]bitset.BitSet, 2)

	p.pieces[piece.WHITE] = whitePieces
	p.pieces[piece.BLACK] = blackPieces

	return p
}

// StartPosition creates a new start position
func StartPosition() Position {
	pieces := make([]map[piece.Piece]bitset.BitSet, 2)

	for _, colour := range piece.AllColours {
		for _, pieceType := range piece.AllPieces {
			pieces[colour][pieceType] = piece.StartPosn[colour][pieceType]
		}
	}

	return NewPosition(pieces[piece.WHITE], pieces[piece.BLACK])
}

// EnpassantSquare returns the current enpassant square (or nil)
func (p Position) EnpassantSquare() *square.Square {
	return p.enpassantSquare
}

// HalfmoveClock returns the current halfmove clock
func (p Position) HalfmoveClock() int {
	return p.halfmoveClock
}

// FullmoveNbr returns the current fullmove nbr
func (p Position) FullmoveNbr() int {
	return p.fullmoveNbr
}

// CastlingAvailability returns the castling availabilty
func (p Position) CastlingAvailability() string {
	return p.castlingAvailability
}

// BitSetFor returns the bitset for the given piece and colour
func (p Position) BitSetFor(colour piece.Colour, piece piece.Piece) bitset.BitSet {
	return p.pieces[colour][piece]
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

// 'package private', only called by Builder
func (p *Position) setEnpassantSquare(sq *square.Square) {
	p.enpassantSquare = sq
}

// 'package private', only called by Builder
func (p *Position) setHalfmoveClock(clock int) {
	p.halfmoveClock = clock
}

// 'package private', only called by Builder
func (p *Position) setFullmoveNbr(move int) {
	p.fullmoveNbr = move
}

// 'package private', only called by Builder
func (p *Position) setCastlingAvailability(rights string) {
	p.castlingAvailability = rights
}
