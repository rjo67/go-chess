package position

import (
	"fmt"
	"strings"

	"github.com/rjo67/chess/bitset"
	"github.com/rjo67/chess/move"
	"github.com/rjo67/chess/piece"
	"github.com/rjo67/chess/piece/colour"
	"github.com/rjo67/chess/ray"
	"github.com/rjo67/chess/square"
)

// various masks to access bits in 'castlingAvailability'
var whiteKingssideMask uint32 = 0x1
var whiteQueenssideMask = whiteKingssideMask << 1
var blackKingssideMask uint32 = 0x4
var blackQueenssideMask = blackKingssideMask << 1

// Position represents a chess position
type Position struct {
	pieces                  []map[piece.Piece]bitset.BitSet // array of map of piece bitsets, array-dim = colour
	allPieces               []bitset.BitSet                 // all pieces of a particular colour
	occupiedSquares         bitset.BitSet                   // all occupied squares
	activeColour            colour.Colour                   // whose move
	castlingAvailability    uint32                          // whether white/black can castle kingsside/queensside (see mask values above)
	enpassantSquare         *square.Square                  // enpassant square of current move
	previousEnpassantSquare *square.Square                  // enpassant square if any in previous move
	halfmoveClock           int
	fullmoveNbr             int
}

// NewPosition creates a new position
// The bitset arrays are in the order as given by the piece constants
func NewPosition(whitePieces, blackPieces map[piece.Piece]bitset.BitSet) Position {
	var p Position
	p.pieces = make([]map[piece.Piece]bitset.BitSet, 2)
	p.allPieces = make([]bitset.BitSet, 2)

	p.pieces[colour.White] = whitePieces
	p.pieces[colour.Black] = blackPieces

	// calculate further bitmaps:
	for _, col := range colour.AllColours {
		p.allPieces[col] = bitset.New(0)
		for _, pieceType := range piece.AllPieces {
			p.allPieces[col] = p.allPieces[col].Or(p.pieces[col][pieceType])
		}
	}
	p.occupiedSquares = p.allPieces[colour.White].Or(p.allPieces[colour.Black])

	return p
}

// MakeMove updates the position with the given move
// NB: the move will store the new castling rights if necessary, therefore takes a pointer object
func (p *Position) MakeMove(m *move.Move) {
	myColour := p.activeColour
	otherColour := myColour.Other()
	if m.IsCastles() {
		// move rook (the king's move will be taken care of later)
		var rooksMove bitset.BitSet
		if m.IsKingsSideCastles() {
			rooksMove = kingssideCastlingsRookMove[myColour]
			// castlingAvailability...Side set lower down
		} else if m.IsQueensSideCastles() {
			rooksMove = queenssideCastlingsRookMove[myColour]
		}
		p.pieces[myColour][piece.ROOK] = p.pieces[myColour][piece.ROOK].Xor(rooksMove)
		p.allPieces[myColour] = p.allPieces[myColour].Xor(rooksMove)
		p.occupiedSquares = p.occupiedSquares.Xor(rooksMove)
	} else if m.IsEnpassant() {
		// remove other-coloured piece, which is not at m.To(), but rather m.EnpassantPawnReallyOn()
		p.pieces[otherColour][m.CapturedPiece()] = p.pieces[otherColour][m.CapturedPiece()].Xor(m.EnpassantPawnRealLocation())
		p.allPieces[otherColour] = p.allPieces[otherColour].Xor(m.EnpassantPawnRealLocation())
	} else if m.IsPromotion() {
		if m.IsCapture() {
			BothBs := bitset.NewFromSquares(m.From(), m.To())
			FromBs := bitset.NewFromSquares(m.From())
			ToBs := bitset.NewFromSquares(m.To())
			p.pieces[myColour][piece.PAWN] = p.pieces[myColour][piece.PAWN].Xor(FromBs)
			p.pieces[myColour][m.PromotedPiece()] = p.pieces[myColour][m.PromotedPiece()].Xor(ToBs)
			p.pieces[otherColour][m.CapturedPiece()] = p.pieces[otherColour][m.CapturedPiece()].Xor(ToBs)
			p.allPieces[myColour] = p.allPieces[myColour].Xor(BothBs)
			p.allPieces[otherColour] = p.allPieces[otherColour].Xor(ToBs)

			// just remove the piece at m.From() -- there is a (new) piece at m.To()
			p.occupiedSquares = p.occupiedSquares.Xor(FromBs)
		} else {
			// remove pawn, add promoted piece
			BothBs := bitset.NewFromSquares(m.From(), m.To())
			FromBs := bitset.NewFromSquares(m.From())
			ToBs := bitset.NewFromSquares(m.To())
			p.pieces[myColour][piece.PAWN] = p.pieces[myColour][piece.PAWN].Xor(FromBs)
			p.pieces[myColour][m.PromotedPiece()] = p.pieces[myColour][m.PromotedPiece()].Xor(ToBs)
			p.allPieces[myColour] = p.allPieces[myColour].Xor(BothBs)

			p.occupiedSquares = p.occupiedSquares.Xor(BothBs)
		}
	} else if m.IsCapture() {
		// remove other-coloured piece at m.To()
		targetBs := bitset.NewFromSquares(m.To())
		p.pieces[otherColour][m.CapturedPiece()] = p.pieces[otherColour][m.CapturedPiece()].Xor(targetBs)
		p.allPieces[otherColour] = p.allPieces[otherColour].Xor(targetBs)
	}
	if !m.IsPromotion() {
		// move our colour piece from m.From() to m.To()
		bs := bitset.NewFromSquares(m.From(), m.To())
		p.pieces[myColour][m.PieceType()] = p.pieces[myColour][m.PieceType()].Xor(bs)
		p.allPieces[myColour] = p.allPieces[myColour].Xor(bs)

		if m.IsEnpassant() {
			// must also clear m.EnpassantPawnRealLocation()
			p.occupiedSquares = p.occupiedSquares.Xor(bs).Xor(m.EnpassantPawnRealLocation())
		} else if m.IsCapture() {
			// must clear m.From(), but m.To() is still occupied
			p.occupiedSquares.Clear(uint(m.From()))
		} else {
			p.occupiedSquares = p.occupiedSquares.Xor(bs)
		}
		// remove castling rights on king move
		if m.IsKingsMove() {
			if p.CastlingAvailabilityKingsSide(myColour) {
				m.SetCastleBeforeMove(true)
				p.ToggleCastlingAvailabilityKingsSide(myColour)
			}
			if p.CastlingAvailabilityQueensSide(myColour) {
				m.SetCastleBeforeMove(false)
				p.ToggleCastlingAvailabilityQueensSide(myColour)
			}
		} else if m.PieceType() == piece.ROOK {
			// remove castling rights on rook move
			if p.CastlingAvailabilityQueensSide(myColour) &&
				((myColour == colour.White && m.From() == square.A1) || (myColour == colour.Black && m.From() == square.A8)) {
				m.SetCastleBeforeMove(false)
				p.ToggleCastlingAvailabilityQueensSide(myColour)
			} else if p.CastlingAvailabilityKingsSide(myColour) &&
				((myColour == colour.White && m.From() == square.H1) || (myColour == colour.Black && m.From() == square.H8)) {
				m.SetCastleBeforeMove(true)
				p.ToggleCastlingAvailabilityKingsSide(myColour)
			}
		}
		// remove castling rights FOR OTHER SIDE if necessary
		if p.CastlingAvailabilityQueensSide(otherColour) &&
			((myColour == colour.Black && m.To() == square.A1) || (myColour == colour.White && m.To() == square.A8)) {
			m.SetOpponentCastleBeforeMove(false)
			p.ToggleCastlingAvailabilityQueensSide(otherColour)
		} else if p.CastlingAvailabilityKingsSide(otherColour) &&
			((myColour == colour.Black && m.To() == square.H1) || (myColour == colour.White && m.To() == square.H8)) {
			m.SetOpponentCastleBeforeMove(true)
			p.ToggleCastlingAvailabilityKingsSide(otherColour)
		}
	}
	p.activeColour = p.activeColour.Other()
	if m.HasEnpassantSquare() {
		sq := m.EnpassantSquare()
		p.previousEnpassantSquare = p.enpassantSquare
		p.enpassantSquare = &sq
	} else {
		p.previousEnpassantSquare = p.enpassantSquare
		p.enpassantSquare = nil
	}
}

// UnmakeMove updates the position with the reverse of the given move
// (does not need to update the move, therefore not a pointer)
func (p *Position) UnmakeMove(m move.Move) {
	myColour := p.activeColour.Other() // position has played the move 'm' which will have changed activeColor to the other side
	otherColour := myColour.Other()
	if m.IsCastles() {
		// move rook back (the king's move will be taken care of later)
		var rooksMove bitset.BitSet
		if m.IsKingsSideCastles() {
			rooksMove = kingssideCastlingsRookMove[myColour]
		} else if m.IsQueensSideCastles() {
			rooksMove = queenssideCastlingsRookMove[myColour]
		}
		p.pieces[myColour][piece.ROOK] = p.pieces[myColour][piece.ROOK].Xor(rooksMove)
		p.allPieces[myColour] = p.allPieces[myColour].Xor(rooksMove)
		p.occupiedSquares = p.occupiedSquares.Xor(rooksMove)
	}
	var enpassantPawnRealLocation bitset.BitSet
	if m.IsEnpassant() {
		// restore other-coloured piece -- not at m.To(), but rather at m.EnpassantPawnReallyOn()
		enpassantPawnRealLocation = m.EnpassantPawnRealLocation()
		p.pieces[otherColour][m.CapturedPiece()] = p.pieces[otherColour][m.CapturedPiece()].Xor(enpassantPawnRealLocation)
		p.allPieces[otherColour] = p.allPieces[otherColour].Xor(enpassantPawnRealLocation)
	} else if m.IsPromotion() {
		if m.IsCapture() {
			BothBs := bitset.NewFromSquares(m.From(), m.To())
			FromBs := bitset.NewFromSquares(m.From())
			ToBs := bitset.NewFromSquares(m.To())
			p.pieces[myColour][piece.PAWN] = p.pieces[myColour][piece.PAWN].Xor(FromBs)
			p.pieces[myColour][m.PromotedPiece()] = p.pieces[myColour][m.PromotedPiece()].Xor(ToBs)
			p.pieces[otherColour][m.CapturedPiece()] = p.pieces[otherColour][m.CapturedPiece()].Xor(ToBs)
			p.allPieces[myColour] = p.allPieces[myColour].Xor(BothBs)
			p.allPieces[otherColour] = p.allPieces[otherColour].Xor(ToBs)

			// just restore the piece at m.From() -- there is already a (new) piece at m.To()
			p.occupiedSquares = p.occupiedSquares.Xor(FromBs)
		} else {
			// restore pawn, remove promoted piece
			BothBs := bitset.NewFromSquares(m.From(), m.To())
			FromBs := bitset.NewFromSquares(m.From())
			ToBs := bitset.NewFromSquares(m.To())
			p.pieces[myColour][piece.PAWN] = p.pieces[myColour][piece.PAWN].Xor(FromBs)
			p.pieces[myColour][m.PromotedPiece()] = p.pieces[myColour][m.PromotedPiece()].Xor(ToBs)
			p.allPieces[myColour] = p.allPieces[myColour].Xor(BothBs)
			p.occupiedSquares = p.occupiedSquares.Xor(BothBs)
		}
	} else if m.IsCapture() {
		// restore other-coloured piece at m.To()
		targetBs := bitset.NewFromSquares(m.To())
		p.pieces[otherColour][m.CapturedPiece()] = p.pieces[otherColour][m.CapturedPiece()].Or(targetBs)
		p.allPieces[otherColour] = p.allPieces[otherColour].Or(targetBs)
	}

	if !m.IsPromotion() {
		bs := bitset.NewFromSquares(m.From(), m.To())
		p.pieces[myColour][m.PieceType()] = p.pieces[myColour][m.PieceType()].Xor(bs)
		p.allPieces[myColour] = p.allPieces[myColour].Xor(bs)

		if m.IsEnpassant() {
			// must set m.EnpassantPawnRealLocation()
			p.occupiedSquares = p.occupiedSquares.Xor(bs).Xor(enpassantPawnRealLocation)
		} else if m.IsCapture() {
			// must set m.From() again, but m.To() is still occupied
			p.occupiedSquares.Set(uint(m.From()))
		} else {
			p.occupiedSquares = p.occupiedSquares.Xor(bs)
		}
		// restore castling rights on king or rook move
		if m.IsKingsMove() {
			if m.CouldCastleBeforeMove(true) {
				p.SetCastlingAvailabilityKingsSide(myColour)
			}
			if m.CouldCastleBeforeMove(false) {
				p.SetCastlingAvailabilityQueensSide(myColour)
			}
		} else if m.PieceType() == piece.ROOK {
			// restore castling rights on rook move
			if m.CouldCastleBeforeMove(false) &&
				((myColour == colour.White && m.From() == square.A1) || (myColour == colour.Black && m.From() == square.A8)) {
				p.SetCastlingAvailabilityQueensSide(myColour)
			} else if m.CouldCastleBeforeMove(true) &&
				((myColour == colour.White && m.From() == square.H1) || (myColour == colour.Black && m.From() == square.H8)) {
				p.SetCastlingAvailabilityKingsSide(myColour)
			}
		}

		// restore castling rights FOR OTHER SIDE if necessary
		if m.OpponentCouldCastleBeforeMove(false) &&
			((myColour == colour.Black && m.To() == square.A1) || (myColour == colour.White && m.To() == square.A8)) {
			p.SetCastlingAvailabilityQueensSide(otherColour)
		} else if m.OpponentCouldCastleBeforeMove(true) &&
			((myColour == colour.Black && m.To() == square.H1) || (myColour == colour.White && m.To() == square.H8)) {
			p.SetCastlingAvailabilityKingsSide(otherColour)
		}
	}

	p.activeColour = p.activeColour.Other()
	p.enpassantSquare = p.previousEnpassantSquare
}

// StartPosition creates a new start position
func StartPosition() Position {
	pieces := make([]map[piece.Piece]bitset.BitSet, 2)

	for _, col := range colour.AllColours {
		pieces[col] = make(map[piece.Piece]bitset.BitSet)
		for _, pieceType := range piece.AllPieces {
			pieces[col][pieceType] = piece.StartPosn[col][pieceType]
		}
	}

	return NewPosition(pieces[colour.White], pieces[colour.Black])
}

// PieceAt returns the piece of the specified colour located at sq
// -- panic if there is no such piece
func (p Position) PieceAt(sq uint, requiredColour colour.Colour) piece.Piece {
	for _, pieceType := range piece.AllPieces {
		if p.pieces[requiredColour][pieceType].IsSet(sq) {
			return pieceType
		}
	}
	panic(fmt.Sprintf("no %s piece found on square %d", requiredColour.String(), sq))
}

// AllPieces returns a bitset with all the occupied squares for the given colour
func (p Position) AllPieces(col colour.Colour) bitset.BitSet {
	return p.allPieces[col]
}

// OccupiedSquares returns the occupied-squares bitset
func (p Position) OccupiedSquares() bitset.BitSet {
	return p.occupiedSquares
}

// Pieces returns the bitset for the given piece type and colour
func (p Position) Pieces(col colour.Colour, pieceType piece.Piece) bitset.BitSet {
	return p.pieces[col][pieceType]
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

// CastlingAvailabilityKingsSide returns the castling availabilty on the kingsside for the given colour
func (p Position) CastlingAvailabilityKingsSide(col colour.Colour) bool {
	if col == colour.White {
		return p.castlingAvailability&whiteKingssideMask == whiteKingssideMask
	}
	return p.castlingAvailability&blackKingssideMask == blackKingssideMask
}

// CastlingAvailabilityQueensSide returns the castling availabilty on the queensside for the given colour
func (p Position) CastlingAvailabilityQueensSide(col colour.Colour) bool {
	if col == colour.White {
		return p.castlingAvailability&whiteQueenssideMask == whiteQueenssideMask
	}
	return p.castlingAvailability&blackQueenssideMask == blackQueenssideMask
}

// ToggleCastlingAvailabilityKingsSide toggles the castling availabilty on the kingsside for the given colour
func (p *Position) ToggleCastlingAvailabilityKingsSide(col colour.Colour) {
	if col == colour.White {
		p.castlingAvailability ^= whiteKingssideMask
	} else {
		p.castlingAvailability ^= blackKingssideMask
	}
}

// SetCastlingAvailabilityKingsSide sets the castling availabilty on the kingsside for the given colour
func (p *Position) SetCastlingAvailabilityKingsSide(col colour.Colour) {
	if col == colour.White {
		p.castlingAvailability |= whiteKingssideMask
	} else {
		p.castlingAvailability |= blackKingssideMask
	}
}

// ToggleCastlingAvailabilityQueensSide toggles the castling availabilty on the queensside for the given colour
func (p *Position) ToggleCastlingAvailabilityQueensSide(col colour.Colour) {
	if col == colour.White {
		p.castlingAvailability ^= whiteQueenssideMask
	} else {
		p.castlingAvailability ^= blackQueenssideMask
	}
}

// SetCastlingAvailabilityQueensSide sets the castling availabilty on the queensside for the given colour
func (p *Position) SetCastlingAvailabilityQueensSide(col colour.Colour) {
	if col == colour.White {
		p.castlingAvailability |= whiteQueenssideMask
	} else {
		p.castlingAvailability |= blackQueenssideMask
	}
}

// BitSetFor returns the bitset for the given piece and colour
func (p Position) BitSetFor(col colour.Colour, piece piece.Piece) bitset.BitSet {
	return p.pieces[col][piece]
}

// Attacks returns a bitset of pieces which attack the given square.
// If requiredColour is specified, only the pieces of this colour will be returned.
// (If requiredColour is AnyColour, all attacks are returned)
func (p Position) Attacks(sq square.Square, requiredColour colour.Colour) bitset.BitSet {
	bs := bitset.New(0)
	var bishops, rooks, queens, knights, pawns, kings bitset.BitSet
	diagonals := p._find2(int(sq), ray.AllBishopDirections)
	rankfiles := p._find2(int(sq), ray.AllRookDirections)

	if requiredColour == colour.AnyColour {
		bishops = diagonals.And(p.Pieces(colour.White, piece.BISHOP).Or(p.Pieces(colour.Black, piece.BISHOP)))
		rooks = rankfiles.And(p.Pieces(colour.White, piece.ROOK).Or(p.Pieces(colour.Black, piece.ROOK)))
		knights = ray.KnightAttackBitSets[sq].And(p.Pieces(colour.White, piece.KNIGHT).Or(p.Pieces(colour.Black, piece.KNIGHT)))
		pawns = (ray.PawnAttackBitSets[colour.White][sq].And(p.Pieces(colour.White, piece.PAWN))).Or(ray.PawnAttackBitSets[colour.Black][sq].And(p.Pieces(colour.Black, piece.PAWN)))
		kings = ray.KingAttackBitSets[sq].And(p.Pieces(colour.White, piece.KING).Or(p.Pieces(colour.Black, piece.KING)))

		allQueens := p.Pieces(colour.White, piece.QUEEN).Or(p.Pieces(colour.Black, piece.QUEEN))
		queens = diagonals.And(allQueens).Or(rankfiles.And(allQueens))
	} else {
		// filter on required colour
		bishops = diagonals.And(p.Pieces(requiredColour, piece.BISHOP))
		rooks = rankfiles.And(p.Pieces(requiredColour, piece.ROOK))
		knights = ray.KnightAttackBitSets[sq].And(p.Pieces(requiredColour, piece.KNIGHT))
		queens = diagonals.And(p.Pieces(requiredColour, piece.QUEEN)).Or(rankfiles.And(p.Pieces(requiredColour, piece.QUEEN)))
		pawns = ray.PawnAttackBitSets[requiredColour][sq].And(p.Pieces(requiredColour, piece.PAWN))
		kings = ray.KingAttackBitSets[sq].And(p.Pieces(requiredColour, piece.KING))
	}

	return bs.Or(bishops).Or(rooks).Or(queens).Or(knights).Or(pawns).Or(kings)
}

// String delivers a string representation of the current position
func (p Position) String() string {
	var squares [64]string

	// populate squares with the bitset contents
	for _, col := range colour.AllColours {
		for _, pieceType := range piece.AllPieces {
			bits := p.pieces[col][pieceType].SetBits()
			for _, sq := range bits {
				squares[sq-1] = pieceType.String(col)
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
