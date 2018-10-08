package move

import (
	"fmt"

	"github.com/rjo67/chess/bitset"
	"github.com/rjo67/chess/piece"
	"github.com/rjo67/chess/position"
	"github.com/rjo67/chess/ray"
	"github.com/rjo67/chess/square"
)

// Move stores information about a move
type Move struct {
	from, to  square.Square
	pieceType piece.Piece
	capture   bool
}

func (m Move) String() string {
	if m.capture {
		return fmt.Sprintf("%sx%s", m.from.String(), m.to.String())
	}
	return fmt.Sprintf("%s-%s", m.from.String(), m.to.String())
}

// FindMoves returns all moves for the given colour in the given position
func FindMoves(posn position.Position, colour piece.Colour) []Move {
	moves := make([]Move, 0, 50) // initially empty, capacity 50

	moves = append(moves, findPawnMoves(posn, colour)...)
	moves = append(moves, findRookMoves(posn, colour)...)
	moves = append(moves, findKnightMoves()...)
	moves = append(moves, findBishopMoves(posn, colour)...)
	moves = append(moves, findQueenMoves(posn, colour)...)
	moves = append(moves, findKingMoves()...)

	return moves
}

func findPawnMoves(posn position.Position, colour piece.Colour) []Move {
	moves := make([]Move, 0, 35)
	otherColour := colour.Other()

	var shift int
	var rankMask bitset.BitSet
	if colour == piece.WHITE {
		shift = 8
		rankMask = bitset.Rank2
	} else {
		shift = -8
		rankMask = bitset.Rank7
	}
	// move all pawns up one square, and again for two squares if starting on rank 2
	pawns := posn.Pieces(colour, piece.PAWN)
	emptySquares := posn.OccupiedSquares().Not()
	oneSquare := pawns.Shift(shift).And(emptySquares)
	twoSquares := pawns.And(rankMask).Shift(shift).And(emptySquares).Shift(shift).And(emptySquares)

	for _, bit := range oneSquare.SetBits() {
		moves = append(moves, Move{pieceType: piece.PAWN, from: square.Square(bit - shift), to: square.Square(bit)})
	}
	for _, bit := range twoSquares.SetBits() {
		moves = append(moves, Move{pieceType: piece.PAWN, from: square.Square(bit - (2 * shift)), to: square.Square(bit)})
	}

	_ = otherColour
	/*
		for _, startSq := range posn.Pieces(colour, piece.P).SetBits() {
			for _, direction := range directions {
				possibleMoves := search2(startSq, direction, posn.OccupiedSquares())
				// the 'blocker' square (if present) will be the last one,
				// only need to check this for the colour of the piece it contains (if any)
				setBits := possibleMoves.SetBits()
				lastSlot := len(setBits) - 1
				for i, bit := range setBits {
					if i == lastSlot {
						if posn.AllPieces(colour).IsSet(uint(bit)) {
							// do nothing - square is occupied with a piece of my own colour
						} else if posn.AllPieces(otherColour).IsSet(uint(bit)) {
							// capture
							moves = append(moves, Move{pieceType: pieceType, from: square.Square(startSq), to: square.Square(bit), capture: true})
						} else {
							// no blocker
							moves = append(moves, Move{pieceType: pieceType, from: square.Square(startSq), to: square.Square(bit)})
						}
					} else {
						moves = append(moves, Move{pieceType: pieceType, from: square.Square(startSq), to: square.Square(bit)})
					}
				}
			}
		}
	*/
	return moves
}

func findRookMoves(posn position.Position, colour piece.Colour) []Move {
	return _find(posn, colour, piece.ROOK, ray.AllRookDirections)
}

func findKnightMoves() []Move {
	moves := make([]Move, 0, 50)
	return moves
}

func findBishopMoves(posn position.Position, colour piece.Colour) []Move {
	return _find(posn, colour, piece.BISHOP, ray.AllBishopDirections)
}

func findQueenMoves(posn position.Position, colour piece.Colour) []Move {
	return _find(posn, colour, piece.QUEEN, ray.AllDirections)
}

func findKingMoves() []Move {
	moves := make([]Move, 0, 50)
	return moves
}

func _find(posn position.Position, colour piece.Colour, pieceType piece.Piece, directions []ray.Direction) []Move {
	moves := make([]Move, 0, 15)
	otherColour := colour.Other()
	for _, startSq := range posn.Pieces(colour, pieceType).SetBits() {
		for _, direction := range directions {
			possibleMoves := search2(startSq, direction, posn.OccupiedSquares())
			// the 'blocker' square (if present) will be the last one,
			// only need to check this for the colour of the piece it contains (if any)
			setBits := possibleMoves.SetBits()
			lastSlot := len(setBits) - 1
			for i, bit := range setBits {
				if i == lastSlot {
					if posn.AllPieces(colour).IsSet(uint(bit)) {
						// do nothing - square is occupied with a piece of my own colour
					} else if posn.AllPieces(otherColour).IsSet(uint(bit)) {
						// capture
						moves = append(moves, Move{pieceType: pieceType, from: square.Square(startSq), to: square.Square(bit), capture: true})
					} else {
						// no blocker
						moves = append(moves, Move{pieceType: pieceType, from: square.Square(startSq), to: square.Square(bit)})
					}
				} else {
					moves = append(moves, Move{pieceType: pieceType, from: square.Square(startSq), to: square.Square(bit)})
				}
			}
		}
	}
	return moves
}

// search2 implements the algorithm as described in secition 2 of http://www.craftychess.com/hyatt/bitmaps.html
// i.e. using 'normal' bitmaps.
// The returned bitset contains all possible squares which can be moved to in the given direction.
// NB: it is up to the caller to determine whether the 'blocker' square (if present, the last 'set-bit' in the bitset) is
//     occupied by an opponent's piece (==capture) or our own colour (==no move)
func search2(startSquare int, direction ray.Direction, occupiedSquares bitset.BitSet) bitset.BitSet {
	attackRay := ray.AttackRays[startSquare][direction]           // get the squares attacked in the given direction
	blockers := attackRay.And(occupiedSquares)                    // returns blockers along the direction
	blockingSquare := direction.NextSetBit(blockers, startSquare) // find the first blocking square
	if blockingSquare != 99 {
		attackRay = attackRay.Xor(ray.AttackRays[blockingSquare][direction]) // remove squares 'after' the blocking square
	}
	return attackRay
}
