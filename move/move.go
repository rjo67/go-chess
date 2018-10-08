package move

import (
	"github.com/rjo67/chess/bitset"
	"github.com/rjo67/chess/piece"
	"github.com/rjo67/chess/position"
	"github.com/rjo67/chess/ray"
	"github.com/rjo67/chess/square"
)

// for each square on the board, a bitset containing the possible moves
// keyed by square 1..64 (H1..A8)
// TODO maybe not necessary, instead use 'Rays' ?
//var rookMoves = [65]bitset.BitSet{
//	bitset.BitSet{}, // first is empty
//	bitset.File(1).Or(bitset.Rank(1)),
//	bitset.File(2).Or(bitset.Rank(1)),
//	bitset.File(3).Or(bitset.Rank(1)),
//	bitset.File(4).Or(bitset.Rank(1)),
//	bitset.File(5).Or(bitset.Rank(1)),
//	bitset.File(6).Or(bitset.Rank(1)),
//	bitset.File(7).Or(bitset.Rank(1)),
//	bitset.File(8).Or(bitset.Rank(1)),
//
//	bitset.File(1).Or(bitset.Rank(2)),
//	bitset.File(2).Or(bitset.Rank(2)),
//	bitset.File(3).Or(bitset.Rank(2)),
//	bitset.File(4).Or(bitset.Rank(2)),
//	bitset.File(5).Or(bitset.Rank(2)),
//	bitset.File(6).Or(bitset.Rank(2)),
//	bitset.File(7).Or(bitset.Rank(2)),
//	bitset.File(8).Or(bitset.Rank(2)),
//
//	bitset.File(1).Or(bitset.Rank(3)),
//	bitset.File(2).Or(bitset.Rank(3)),
//	bitset.File(3).Or(bitset.Rank(3)),
//	bitset.File(4).Or(bitset.Rank(3)),
//	bitset.File(5).Or(bitset.Rank(3)),
//	bitset.File(6).Or(bitset.Rank(3)),
//	bitset.File(7).Or(bitset.Rank(3)),
//	bitset.File(8).Or(bitset.Rank(3)),
//
//	bitset.File(1).Or(bitset.Rank(4)),
//	bitset.File(2).Or(bitset.Rank(4)),
//	bitset.File(3).Or(bitset.Rank(4)),
//	bitset.File(4).Or(bitset.Rank(4)),
//	bitset.File(5).Or(bitset.Rank(4)),
//	bitset.File(6).Or(bitset.Rank(4)),
//	bitset.File(7).Or(bitset.Rank(4)),
//	bitset.File(8).Or(bitset.Rank(4)),
//
//	bitset.File(1).Or(bitset.Rank(5)),
//	bitset.File(2).Or(bitset.Rank(5)),
//	bitset.File(3).Or(bitset.Rank(5)),
//	bitset.File(4).Or(bitset.Rank(5)),
//	bitset.File(5).Or(bitset.Rank(5)),
//	bitset.File(6).Or(bitset.Rank(5)),
//	bitset.File(7).Or(bitset.Rank(5)),
//	bitset.File(8).Or(bitset.Rank(5)),
//
//	bitset.File(1).Or(bitset.Rank(6)),
//	bitset.File(2).Or(bitset.Rank(6)),
//	bitset.File(3).Or(bitset.Rank(6)),
//	bitset.File(4).Or(bitset.Rank(6)),
//	bitset.File(5).Or(bitset.Rank(6)),
//	bitset.File(6).Or(bitset.Rank(6)),
//	bitset.File(7).Or(bitset.Rank(6)),
//	bitset.File(8).Or(bitset.Rank(6)),
//
//	bitset.File(1).Or(bitset.Rank(7)),
//	bitset.File(2).Or(bitset.Rank(7)),
//	bitset.File(3).Or(bitset.Rank(7)),
//	bitset.File(4).Or(bitset.Rank(7)),
//	bitset.File(5).Or(bitset.Rank(7)),
//	bitset.File(6).Or(bitset.Rank(7)),
//	bitset.File(7).Or(bitset.Rank(7)),
//	bitset.File(8).Or(bitset.Rank(7)),
//
//	bitset.File(1).Or(bitset.Rank(8)),
//	bitset.File(2).Or(bitset.Rank(8)),
//	bitset.File(3).Or(bitset.Rank(8)),
//	bitset.File(4).Or(bitset.Rank(8)),
//	bitset.File(5).Or(bitset.Rank(8)),
//	bitset.File(6).Or(bitset.Rank(8)),
//	bitset.File(7).Or(bitset.Rank(8)),
//	bitset.File(8).Or(bitset.Rank(8)),
//}

// Move stores information about a move
type Move struct {
	from, to square.Square
}

// FindMoves returns all moves for the given colour in the given position
func FindMoves(posn position.Position, colour piece.Colour) []Move {
	moves := make([]Move, 0, 50) // initially empty, capacity 50

	moves = append(moves, findPawnMoves()...)
	moves = append(moves, findRookMoves(posn, colour)...)
	moves = append(moves, findKnightMoves()...)
	moves = append(moves, findBishopMoves()...)
	moves = append(moves, findQueenMoves()...)
	moves = append(moves, findKingMoves()...)

	return moves
}

func findPawnMoves() []Move {
	moves := make([]Move, 0, 50)
	return moves
}

func findRookMoves(posn position.Position, colour piece.Colour) []Move {
	moves := make([]Move, 0, 15)
	for _, rook := range posn.Pieces(colour, piece.ROOK).SetBits() {
		for _, direction := range ray.AllDirections {
			possibleMoves := search2(rook, direction, posn.OccupiedSquares())
			_ = possibleMoves
			// process possibleMoves
		}

	}
	return moves
}

func findKnightMoves() []Move {
	moves := make([]Move, 0, 50)
	return moves
}

func findBishopMoves() []Move {
	moves := make([]Move, 0, 50)
	return moves
}

func findQueenMoves() []Move {
	moves := make([]Move, 0, 50)
	return moves
}

func findKingMoves() []Move {
	moves := make([]Move, 0, 50)
	return moves
}

// this implements the algorithm as described in part 2 of http://www.craftychess.com/hyatt/bitmaps.html
// (using 'normal' bitmaps)
func search2(startSquare int, direction ray.Direction, occupiedSquares bitset.BitSet) bitset.BitSet {
	attackRay := ray.AttackRays[startSquare][direction]           // get the squares attacked in the given direction
	blockers := attackRay.And(occupiedSquares)                    // returns blockers along the direction
	blockingSquare := direction.NextSetBit(blockers, startSquare) // find the first blocking square
	if blockingSquare != 99 {
		attackRay = attackRay.Xor(ray.AttackRays[blockingSquare][direction]) // remove squares 'after' the blocking square
	}
	return attackRay
}
