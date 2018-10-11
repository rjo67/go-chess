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
	col           colour.Colour
	from, to      square.Square
	pieceType     piece.Piece
	castle        bool
	promotedPiece *piece.Piece // set if promotion
	capturedPiece *piece.Piece // set if capture
}

// New creates a new non-capture move
func New(col colour.Colour, from, to square.Square, pieceType piece.Piece) Move {
	return Move{col: col, from: from, to: to, pieceType: pieceType}
}

// NewCapture creates a new capture move
func NewCapture(col colour.Colour, from, to square.Square, pieceType piece.Piece, capturedPieceType piece.Piece) Move {
	return Move{col: col, from: from, to: to, pieceType: pieceType, capturedPiece: &capturedPieceType}
}

// NewPromotion creates a new promotion move
func NewPromotion(col colour.Colour, from, to square.Square, toPiece piece.Piece) Move {
	return Move{col: col, from: from, to: to, pieceType: piece.PAWN, promotedPiece: &toPiece}
}

// NewPromotionCapture creates a new promotion move with capture
func NewPromotionCapture(col colour.Colour, from, to square.Square, toPiece piece.Piece, capturedPieceType piece.Piece) Move {
	return Move{col: col, from: from, to: to, pieceType: piece.PAWN, promotedPiece: &toPiece, capturedPiece: &capturedPieceType}
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

// IsKingsMove returns true when this move involves the king (castling excluded)
func (m Move) IsKingsMove() bool { return m.pieceType == piece.KING }

// IsCastles returns true when this move was "castles"
func (m Move) IsCastles() bool { return m.castle }

// IsCapture returns true when this move was a capture
func (m Move) IsCapture() bool { return m.capturedPiece != nil }

// CapturedPiece returns the captured piece (only call if IsCapture()==true)
func (m Move) CapturedPiece() piece.Piece { return *m.capturedPiece }

// From returns the move's source square
func (m Move) From() square.Square { return m.from }

// To returns the move's target square
func (m Move) To() square.Square { return m.to }

// Colour returns the colour of the move
func (m Move) Colour() colour.Colour { return m.col }

// PieceType returns the move's piece
func (m Move) PieceType() piece.Piece { return m.pieceType }

func (m Move) String() string {
	if m.castle {
		if m.to == square.G1 || m.to == square.G8 {
			return "O-O"
		}
		return "O-O-O"
	}
	var promotion string
	if m.promotedPiece != nil {
		promotion = fmt.Sprintf("=%s", m.promotedPiece.String(colour.White))
	}
	if m.capturedPiece != nil {
		return fmt.Sprintf("%sx%s%s", m.from.String(), m.to.String(), promotion)
	}
	return fmt.Sprintf("%s%s%s", m.from.String(), m.to.String(), promotion)
}

// Search2 implements the algorithm as described in secition 2 of http://www.craftychess.com/hyatt/bitmaps.html
// i.e. using 'normal' bitmaps.
// The returned bitset contains all possible squares which can be moved to in the given direction.
// The blockingSquare will be set to the 'blocker' or 99 if no blocker
// NB: it is up to the caller to determine whether the 'blocker' square is occupied by an opponent's piece (==capture) or our own colour (==no move)
func Search2(startSquare int, direction ray.Direction, occupiedSquares bitset.BitSet) (attackRay bitset.BitSet, blockingSquare int) {
	attackRay = ray.AttackRays[startSquare][direction]           // get the squares attacked in the given direction
	blockers := attackRay.And(occupiedSquares)                   // returns blockers along the direction
	blockingSquare = direction.NextSetBit(blockers, startSquare) // find the first blocking square
	if blockingSquare != 99 {
		attackRay = attackRay.Xor(ray.AttackRays[blockingSquare][direction]) // remove squares 'after' the blocking square
	}
	return attackRay, blockingSquare
}
