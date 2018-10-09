package move

import (
	"fmt"

	"github.com/rjo67/chess/bitset"
	"github.com/rjo67/chess/piece"
	"github.com/rjo67/chess/piece/colour"
	"github.com/rjo67/chess/ray"
	"github.com/rjo67/chess/square"
)

// Move stores information about a move
type Move struct {
	col       colour.Colour
	from, to  square.Square
	pieceType piece.Piece
	capture   bool
	castle    bool
}

// New creates a new non-capture move
func New(col colour.Colour, from, to square.Square, pieceType piece.Piece) Move {
	return Move{col: col, from: from, to: to, pieceType: pieceType, capture: false}
}

// NewCapture creates a new capture move
func NewCapture(col colour.Colour, from, to square.Square, pieceType piece.Piece) Move {
	return Move{col: col, from: from, to: to, pieceType: pieceType, capture: true}
}

// CastleKingsSide creates O-O
func CastleKingsSide(col colour.Colour) Move {
	if col == colour.White {
		return Move{col: col, from: square.E1, to: square.G1, castle: true, pieceType: piece.KING}
	}
	return Move{col: col, from: square.E8, to: square.G8, castle: true, pieceType: piece.KING}
}

// CastleQueensSide creates O-O-O
func CastleQueensSide(col colour.Colour) Move {
	if col == colour.White {
		return Move{col: col, from: square.E1, to: square.C1, castle: true, pieceType: piece.KING}
	}
	return Move{col: col, from: square.E8, to: square.C8, castle: true, pieceType: piece.KING}
}

func (m Move) String() string {
	if m.castle {
		if m.to == square.G1 || m.to == square.G8 {
			return "O-O"
		}
		return "O-O-O"
	}
	if m.capture {
		return fmt.Sprintf("%sx%s", m.from.String(), m.to.String())
	}
	return fmt.Sprintf("%s%s", m.from.String(), m.to.String())
}

// Search2 implements the algorithm as described in secition 2 of http://www.craftychess.com/hyatt/bitmaps.html
// i.e. using 'normal' bitmaps.
// The returned bitset contains all possible squares which can be moved to in the given direction.
// NB: it is up to the caller to determine whether the 'blocker' square (if present, the last 'set-bit' in the bitset) is
//     occupied by an opponent's piece (==capture) or our own colour (==no move)
func Search2(startSquare int, direction ray.Direction, occupiedSquares bitset.BitSet) bitset.BitSet {
	attackRay := ray.AttackRays[startSquare][direction]           // get the squares attacked in the given direction
	blockers := attackRay.And(occupiedSquares)                    // returns blockers along the direction
	blockingSquare := direction.NextSetBit(blockers, startSquare) // find the first blocking square
	if blockingSquare != 99 {
		attackRay = attackRay.Xor(ray.AttackRays[blockingSquare][direction]) // remove squares 'after' the blocking square
	}
	return attackRay
}
