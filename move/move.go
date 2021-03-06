package move

import (
	"fmt"

	"github.com/rjo67/chess/bitset"
	"github.com/rjo67/chess/piece"
	"github.com/rjo67/chess/piece/colour"
	"github.com/rjo67/chess/ray"
	"github.com/rjo67/chess/square"
)

// various masks to access bits in 'castlingInfo'
var myColourKingssideMask uint32 = 0x1
var myColourQueenssideMask = myColourKingssideMask << 1
var opponentsColourKingssideMask uint32 = 0x4
var opponentsColourQueenssideMask = opponentsColourKingssideMask << 1

var fromSquareMask uint32 = 0x3F                                  // bits 1-6
var toSquareMask uint32 = 0xFC0                                   // bits 7-12
var castlingKingssideMask uint32 = 0x1000                         // bit 13
var castlingQueenssideMask = castlingKingssideMask << 1           // bit 14
var castlingMask = castlingKingssideMask | castlingQueenssideMask // 13 or 14
var movingPieceMask uint32 = 0x1C000                              // bit 15-17
var enpassantMask = castlingQueenssideMask << 4                   // bit 18
var promotionMask = enpassantMask << 1                            // bit 19
var promotionPieceMask uint32 = 0x180000                          // bit 20..21

// Move stores information about a move.
// info: bits 1..6  'from' square   (0..63)
//       bits 7..12 'to' square     (0..63)
//       bits 13-14 if the move was castles: King-side (1), Queen-side (2)
//       bits 15-17 type of moving piece
//       bit 18 if the move was enpassant
//       bit 19 if the move was promotion
//       bits 20-21 promotion piece type
//
// Castling info is stored in the int 'castlingInfo':
//   Bit 1, 2: whether could castle kingsside/queensside before making this move (mask: myColourKingssideMask, myColourQueenssideMask)
//   Bit 3, 4: whether OPPONENT could castle kingsside/queenside before making this move (mask: opponentsColourKingssideMask, opponentsColourQueensssideMask)
type Move struct {
	info                      uint32         // info about the move (see above)
	castlingInfo              uint32         // stores if we or opponent could castle before making this move. See description above (set during posn.MakeMove)
	capturedPiece             *piece.Piece   // set if capture
	enpassantPawnRealLocation *bitset.BitSet // set if e.p., contains the 'real' square where the pawn was, e.g. move.To()==E6, enpassantPawnRealLocation==E5
	enpassantSquare           *square.Square // set to the enpassant square if this move is a pawn move from rank2 to rank4
}

// New creates a new non-capture move
func New(col colour.Colour, from, to square.Square, pieceType piece.Piece) Move {
	m := Move{}
	// squares are stored in 6 bits 0..63
	m.info |= uint32((from - 1))
	m.info |= uint32((to - 1) << 6)
	m.info |= (uint32(pieceType)) << 14
	if pieceType == piece.PAWN {
		if col == colour.White && from.Rank() == 2 && to.Rank() == 4 {
			sq := square.Square(from + 8)
			m.enpassantSquare = &sq
		} else if col == colour.Black && from.Rank() == 7 && to.Rank() == 5 {
			sq := square.Square(from - 8)
			m.enpassantSquare = &sq
		}
	}
	return m
}

// NewCapture creates a new capture move
func NewCapture(col colour.Colour, from, to square.Square, pieceType piece.Piece, capturedPieceType piece.Piece) Move {
	m := New(col, from, to, pieceType)
	m.capturedPiece = &capturedPieceType
	return m
}

// NewEpCapture creates a new enpassant capture move
func NewEpCapture(col colour.Colour, from, to square.Square) Move {
	var shift int
	if col == colour.White {
		shift = -8 // e.g. black pawn is on E5, enpassant square is E6
	} else {
		shift = 8
	}
	enpassantPawnRealLocation := bitset.NewFromSquares(square.Square(int(to) + shift))
	m := NewCapture(col, from, to, piece.PAWN, piece.PAWN)
	m.info |= enpassantMask
	m.enpassantPawnRealLocation = &enpassantPawnRealLocation
	return m
}

// NewPromotion creates a new promotion move
func NewPromotion(col colour.Colour, from, to square.Square, toPiece piece.Piece) Move {
	m := New(col, from, to, piece.PAWN)
	m.info |= promotionMask
	m.info |= (uint32(toPiece) - 1) << 19
	return m
}

// NewPromotionCapture creates a new promotion move with capture
func NewPromotionCapture(col colour.Colour, from, to square.Square, toPiece piece.Piece, capturedPieceType piece.Piece) Move {
	m := NewCapture(col, from, to, piece.PAWN, capturedPieceType)
	m.info |= promotionMask
	m.info |= (uint32(toPiece) - 1) << 19
	return m
}

// CastleKingsSide creates O-O
func CastleKingsSide(col colour.Colour) Move {
	var m Move
	if col == colour.White {
		m = New(col, square.E1, square.G1, piece.KING)
	} else {
		m = New(col, square.E8, square.G8, piece.KING)
	}
	m.info |= castlingKingssideMask
	return m
}

// CastleQueensSide creates O-O-O
func CastleQueensSide(col colour.Colour) Move {
	var m Move
	if col == colour.White {
		m = New(col, square.E1, square.C1, piece.KING)
	} else {
		m = New(col, square.E8, square.C8, piece.KING)
	}
	m.info |= castlingQueenssideMask
	return m
}

// IsKingsMove returns true if this move involves the king (castling excluded)
func (m Move) IsKingsMove() bool { return m.PieceType() == piece.KING }

// IsCastles returns true if this move was "castles"
func (m Move) IsCastles() bool { return m.info&castlingMask != 0 }

// IsKingsSideCastles returns true if this move was kings-side "castles"
func (m Move) IsKingsSideCastles() bool {
	return m.info&castlingKingssideMask == castlingKingssideMask
}

// IsQueensSideCastles returns true if this move was queens-side "castles"
func (m Move) IsQueensSideCastles() bool {
	return m.info&castlingQueenssideMask == castlingQueenssideMask
}

// IsCapture returns true if this move was a capture
func (m Move) IsCapture() bool { return m.capturedPiece != nil }

// IsPromotion returns true if this move was a pawn promotion
func (m Move) IsPromotion() bool { return m.info&promotionMask == promotionMask }

// IsEnpassant returns true if this move was an enpassant capture
func (m Move) IsEnpassant() bool { return m.info&enpassantMask == enpassantMask }

// EnpassantPawnRealLocation returns a bitset containing the square where the opponents pawn really was for an enpassant capture
func (m Move) EnpassantPawnRealLocation() bitset.BitSet { return *m.enpassantPawnRealLocation }

// HasEnpassantSquare returns true if enpassantSquare is set
func (m Move) HasEnpassantSquare() bool { return m.enpassantSquare != nil }

// EnpassantSquare returns the enpassant square
func (m Move) EnpassantSquare() square.Square { return *m.enpassantSquare }

// CapturedPiece returns the captured piece (only call if IsCapture()==true)
func (m Move) CapturedPiece() piece.Piece { return *m.capturedPiece }

// PromotedPiece returns the piece which the pawn has promoted to (only call if IsPromotion()==true)
func (m Move) PromotedPiece() piece.Piece { return piece.Piece((m.info&promotionPieceMask)>>19 + 1) }

// From returns the move's source square
func (m Move) From() square.Square { return square.Square(m.info&fromSquareMask + 1) } // mask bits 1..6

// To returns the move's target square
func (m Move) To() square.Square { return square.Square((m.info&toSquareMask)>>6 + 1) } // mask bits 7..12

// PieceType returns the move's piece
func (m Move) PieceType() piece.Piece { return piece.Piece((m.info & movingPieceMask) >> 14) }

// CouldCastleBeforeMove returns true if it was possible to castle before this move
func (m Move) CouldCastleBeforeMove(kingsside bool) bool {
	if kingsside {
		return m.castlingInfo&myColourKingssideMask == myColourKingssideMask
	}
	return m.castlingInfo&myColourQueenssideMask == myColourQueenssideMask
}

// SetCastleBeforeMove sets the flag indicating whether it was possible to castle before this move
func (m *Move) SetCastleBeforeMove(kingsside bool) {
	if kingsside {
		m.castlingInfo |= myColourKingssideMask
	} else {
		m.castlingInfo |= myColourQueenssideMask
	}
}

// OpponentCouldCastleBeforeMove returns true if it was possible FOR THE OTHER SIDE to castle before this move
func (m Move) OpponentCouldCastleBeforeMove(kingsside bool) bool {
	if kingsside {
		return m.castlingInfo&opponentsColourKingssideMask == opponentsColourKingssideMask
	}
	return m.castlingInfo&opponentsColourQueenssideMask == opponentsColourQueenssideMask
}

// SetOpponentCastleBeforeMove sets the flag indicating whether it was possible FOR THE OTHER SIDE to castle before this move
func (m *Move) SetOpponentCastleBeforeMove(kingsside bool) {
	if kingsside {
		m.castlingInfo |= opponentsColourKingssideMask
	} else {
		m.castlingInfo |= opponentsColourQueenssideMask
	}
}

func (m Move) String() string {
	if m.IsKingsSideCastles() {
		return "O-O"
	} else if m.IsQueensSideCastles() {
		return "O-O-O"
	}
	var promotion string
	if m.IsPromotion() {
		promotion = fmt.Sprintf("=%s", m.PromotedPiece().String(colour.White))
	}
	if m.capturedPiece != nil {
		return fmt.Sprintf("%sx%s%s", m.From().String(), m.To().String(), promotion)
	}
	return fmt.Sprintf("%s%s%s", m.From().String(), m.To().String(), promotion)
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
