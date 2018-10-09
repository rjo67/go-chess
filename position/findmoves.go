package position

import (
	"github.com/rjo67/chess/bitset"
	"github.com/rjo67/chess/move"
	"github.com/rjo67/chess/piece"
	"github.com/rjo67/chess/piece/colour"
	"github.com/rjo67/chess/ray"
	"github.com/rjo67/chess/square"
)

// FindMoves returns all moves for the given colour in the given position
func (p Position) FindMoves(col colour.Colour) []move.Move {
	moves := make([]move.Move, 0, 50) // initially empty, capacity 50

	moves = append(moves, p.findPawnMoves(col)...)
	moves = append(moves, p.findRookMoves(col)...)
	moves = append(moves, p.findKnightMoves(col)...)
	moves = append(moves, p.findBishopMoves(col)...)
	moves = append(moves, p.findQueenMoves(col)...)
	moves = append(moves, p.findKingMoves(col)...)

	return moves
}

func (p Position) findPawnMoves(col colour.Colour) []move.Move {
	moves := make([]move.Move, 0, 50)

	var shift int
	var rankMask bitset.BitSet
	if col == colour.White {
		shift = 8
		rankMask = bitset.Rank2
	} else {
		shift = -8
		rankMask = bitset.Rank7
	}
	// move all pawns up one square, and again for two squares if starting on rank 2
	pawns := p.Pieces(col, piece.PAWN)
	emptySquares := p.OccupiedSquares().Not()
	oneSquare := pawns.Shift(shift).And(emptySquares)
	twoSquares := pawns.And(rankMask).Shift(shift).And(emptySquares).Shift(shift).And(emptySquares)

	for _, bit := range oneSquare.SetBits() {
		moves = append(moves, move.New(col, square.Square(bit-shift), square.Square(bit), piece.PAWN))
	}
	for _, bit := range twoSquares.SetBits() {
		moves = append(moves, move.New(col, square.Square(bit-(2*shift)), square.Square(bit), piece.PAWN))
	}

	//
	// captures...
	//

	otherColour := col.Other()

	if col == colour.White {
		shift = 9
		rankMask = bitset.NotFile1
	} else {
		shift = -9
		rankMask = bitset.NotFile8
	}
	captureLeft := pawns.And(rankMask).Shift(shift).And(p.AllPieces(otherColour))
	for _, bit := range captureLeft.SetBits() {
		moves = append(moves, move.NewCapture(col, square.Square(bit-shift), square.Square(bit), piece.PAWN))
	}

	if col == colour.White {
		shift = 7
		rankMask = bitset.NotFile8
	} else {
		shift = -7
		rankMask = bitset.NotFile1
	}
	captureRight := pawns.And(rankMask).Shift(shift).And(p.OccupiedSquares())
	for _, bit := range captureRight.SetBits() {
		moves = append(moves, move.NewCapture(col, square.Square(bit-shift), square.Square(bit), piece.PAWN))
	}

	//TODO: enpassant and promotion

	return moves
}

func (p Position) findRookMoves(col colour.Colour) []move.Move {
	return p._findForPiece(col, piece.ROOK, ray.AllRookDirections)
}

func (p Position) findKnightMoves(col colour.Colour) []move.Move {
	moves := make([]move.Move, 0, 8)
	otherColour := col.Other()
	for _, startSq := range p.Pieces(col, piece.KNIGHT).SetBits() {
		bs := ray.KnightAttackBitSets[startSq]
		for _, bit := range bs.SetBits() {
			if p.AllPieces(col).IsSet(uint(bit)) {
				// do nothing - square is occupied with a piece of my own colour
			} else if p.AllPieces(otherColour).IsSet(uint(bit)) {
				// capture
				moves = append(moves, move.NewCapture(col, square.Square(startSq), square.Square(bit), piece.KNIGHT))
			} else {
				// empty square
				moves = append(moves, move.New(col, square.Square(startSq), square.Square(bit), piece.KNIGHT))
			}
		}
	}
	return moves
}

func (p Position) findBishopMoves(col colour.Colour) []move.Move {
	return p._findForPiece(col, piece.BISHOP, ray.AllBishopDirections)
}

func (p Position) findQueenMoves(col colour.Colour) []move.Move {
	return p._findForPiece(col, piece.QUEEN, ray.AllDirections)
}

func (p Position) findKingMoves(col colour.Colour) []move.Move {
	moves := make([]move.Move, 0, 10)
	otherColour := col.Other()
	for _, startSq := range p.Pieces(col, piece.KING).SetBits() {
		bs := ray.KingAttackBitSets[startSq]
		for _, bit := range bs.SetBits() {
			if p.AllPieces(col).IsSet(uint(bit)) {
				// do nothing - square is occupied with a piece of my own colour
			} else if p.AllPieces(otherColour).IsSet(uint(bit)) {
				// capture
				moves = append(moves, move.NewCapture(col, square.Square(startSq), square.Square(bit), piece.KING))
			} else {
				// empty square
				moves = append(moves, move.New(col, square.Square(startSq), square.Square(bit), piece.KING))
			}
		}
	}
	if p.CastlingAvailability(col, true) {
		// check if castling allowed
		moves = append(moves, move.CastleKingsSide(col))
	}
	if p.CastlingAvailability(col, false) {
		// check if castling allowed
		moves = append(moves, move.CastleQueensSide(col))
	}
	return moves
}

// Finds all moves for a given pieceType and Colour, in all given directions,
// using the piece squares from the current position
func (p Position) _findForPiece(col colour.Colour, pieceType piece.Piece, directions []ray.Direction) []move.Move {
	moves := make([]move.Move, 0, 15)
	otherColour := col.Other()
	for _, startSq := range p.Pieces(col, pieceType).SetBits() {
		//
		// we iterate over all directions separately (instead of calling _find2)
		// because we only want to use the 'blocker' logic on the last square of each direction
		for _, direction := range directions {
			possibleMoves := move.Search2(startSq, direction, p.OccupiedSquares())
			// the 'blocker' square (if present) will be the last one,
			// only need to check this for the colour of the piece it contains (if any)
			setBits := possibleMoves.SetBits()
			lastSlot := len(setBits) - 1
			for i, bit := range setBits {
				if i == lastSlot {
					if p.AllPieces(col).IsSet(uint(bit)) {
						// do nothing - square is occupied with a piece of my own colour
					} else if p.AllPieces(otherColour).IsSet(uint(bit)) {
						// capture
						moves = append(moves, move.NewCapture(col, square.Square(startSq), square.Square(bit), pieceType))
					} else {
						// no blocker
						moves = append(moves, move.New(col, square.Square(startSq), square.Square(bit), pieceType))
					}
				} else {
					moves = append(moves, move.New(col, square.Square(startSq), square.Square(bit), pieceType))
				}
			}
		}
	}
	return moves
}

// Returns a bitset of all squares which can be attacked in the given directions from the given square.
// Alternatively: returns a bitset of all squares which attack the given square in the given directions.
// TODO rename to AttacksOnSquare
func (p Position) _find2(startSq int, directions []ray.Direction) bitset.BitSet {
	possibleMoves := bitset.BitSet{}
	for _, direction := range directions {
		possibleMoves = possibleMoves.Or(move.Search2(startSq, direction, p.OccupiedSquares()))
	}
	return possibleMoves
}
