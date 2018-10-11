package ray

import (
	"github.com/rjo67/chess/bitset"
	"github.com/rjo67/chess/piece/colour"
	"github.com/rjo67/chess/square"
)

// Direction is a ray direction
type Direction uint32

// the ray directions
const (
	NORTH Direction = iota
	NORTHEAST
	EAST
	SOUTHEAST
	SOUTH
	SOUTHWEST
	WEST
	NORTHWEST
)

// AllDirections to iterate over the ray types
var AllDirections = []Direction{NORTH, NORTHEAST, EAST, SOUTHEAST, SOUTH, SOUTHWEST, WEST, NORTHWEST}

// AllBishopDirections to iterate over the ray types applicable for bishop moves
var AllBishopDirections = []Direction{NORTHWEST, NORTHEAST, SOUTHWEST, SOUTHEAST}

// AllRookDirections to iterate over the ray types applicable for rook moves
var AllRookDirections = []Direction{NORTH, EAST, SOUTH, WEST}

// NextSetBit returns the next set bit in the given direction, or 99 if there wasn't one
func (dir Direction) NextSetBit(bs bitset.BitSet, start int) int {
	switch dir {
	case NORTH:
		for bit := start + 8; bit < 65; bit += 8 {
			if bs.IsSet(uint(bit)) {
				return bit
			}
		}
	case NORTHEAST:
		if !onRightHandSide(start) {
			for bit := start + 7; bit < 65; bit += 7 {
				if bs.IsSet(uint(bit)) {
					return bit
				}
				if onRightHandSide(bit) {
					break
				}
			}
		}
	case EAST:
		if !onRightHandSide(start) {
			for bit := start - 1; bit > 0; bit-- {
				if bs.IsSet(uint(bit)) {
					return bit
				}
			}
		}
	case SOUTHEAST:
		if !onRightHandSide(start) {
			for bit := start - 9; bit > 0; bit -= 9 {
				if bs.IsSet(uint(bit)) {
					return bit
				}
				if onRightHandSide(bit) {
					break
				}
			}
		}
	case SOUTH:
		for bit := start - 8; bit > 0; bit -= 8 {
			if bs.IsSet(uint(bit)) {
				return bit
			}
		}
	case SOUTHWEST:
		if !onLeftHandSide(start) {
			for bit := start - 7; bit > 0; bit -= 7 {
				if bs.IsSet(uint(bit)) {
					return bit
				}
				if onLeftHandSide(bit) {
					break
				}
			}
		}
	case WEST:
		if !onLeftHandSide(start) {
			for bit := start + 1; bit < 65; bit++ {
				if bs.IsSet(uint(bit)) {
					return bit
				}
				if onLeftHandSide(bit) {
					break
				}
			}
		}
	case NORTHWEST:
		if !onLeftHandSide(start) {
			for bit := start + 9; bit < 65; bit += 9 {
				if bs.IsSet(uint(bit)) {
					return bit
				}
				if onLeftHandSide(bit) {
					break
				}
			}
		}
	default:
		panic("oops!?")
	}
	return 99
}

// AttackRay stores, for each square, the squares which are potentially attacked by a sliding piece on that square
type AttackRay [8]bitset.BitSet

// internal constructor
func newKnightAttackBitset(start int) bitset.BitSet {
	bs := bitset.New(0)
	startSq := square.Square(start)
	rank := startSq.Rank()
	file := startSq.File()

	if file-2 > 0 && rank+1 < 9 {
		bs.SetSquare(square.FromRankAndFile(rank+1, file-2))
	}
	if file-1 > 0 && rank+2 < 9 {
		bs.SetSquare(square.FromRankAndFile(rank+2, file-1))
	}
	if file+1 < 9 && rank+2 < 9 {
		bs.SetSquare(square.FromRankAndFile(rank+2, file+1))
	}
	if file+2 < 9 && rank+1 < 9 {
		bs.SetSquare(square.FromRankAndFile(rank+1, file+2))
	}
	if file+2 < 9 && rank-1 > 0 {
		bs.SetSquare(square.FromRankAndFile(rank-1, file+2))
	}
	if file+1 < 9 && rank-2 > 0 {
		bs.SetSquare(square.FromRankAndFile(rank-2, file+1))
	}
	if file-1 > 0 && rank-2 > 0 {
		bs.SetSquare(square.FromRankAndFile(rank-2, file-1))
	}
	if file-2 > 0 && rank-1 > 0 {
		bs.SetSquare(square.FromRankAndFile(rank-1, file-2))
	}

	return bs
}

// internal constructor. Pawns on which squares attack square 'target'?
func newPawnAttackBitset(col colour.Colour, target int) bitset.BitSet {
	switch col {
	case colour.White:
		if target < 16 {
			return bitset.New(0)
		}
		// check for LHS or RHS
		bs := bitset.New(0)
		if target%8 == 0 {
			return *bs.Set(uint(target - 9))
		} else if target%8 == 1 {
			return *bs.Set(uint(target - 7))
		}
		return *bs.Set(uint(target - 7)).Set(uint(target - 9))
	case colour.Black:
		if target > 48 {
			return bitset.New(0)
		}
		// check for LHS or RHS
		bs := bitset.New(0)
		if target%8 == 0 {
			return *bs.Set(uint(target + 7))
		} else if target%8 == 1 {
			return *bs.Set(uint(target + 9))
		}
		return *bs.Set(uint(target + 7)).Set(uint(target + 9))
	default:
		panic("oops")
	}
}

func newKingAttackBitset(start int) bitset.BitSet {
	bs := bitset.New(0)
	startSq := square.Square(start)
	rank := startSq.Rank()
	file := startSq.File()

	if file-1 > 0 && rank-1 > 0 {
		bs.SetSquare(square.FromRankAndFile(rank-1, file-1))
	}
	if file-1 > 0 {
		bs.SetSquare(square.FromRankAndFile(rank, file-1))
	}
	if file-1 > 0 && rank+1 < 9 {
		bs.SetSquare(square.FromRankAndFile(rank+1, file-1))
	}
	if rank+1 < 9 {
		bs.SetSquare(square.FromRankAndFile(rank+1, file))
	}
	if file+1 < 9 && rank+1 < 9 {
		bs.SetSquare(square.FromRankAndFile(rank+1, file+1))
	}
	if file+1 < 9 {
		bs.SetSquare(square.FromRankAndFile(rank, file+1))
	}
	if file+1 < 9 && rank-1 > 0 {
		bs.SetSquare(square.FromRankAndFile(rank-1, file+1))
	}
	if rank-1 > 0 {
		bs.SetSquare(square.FromRankAndFile(rank-1, file))
	}

	return bs
}

// internal constructor
func newAttackRay(start int) AttackRay {
	ar := AttackRay{}
	for i, rayType := range AllDirections {
		var bs = bitset.New(0)
		switch rayType {
		case NORTH:
			for bit := start + 8; bit < 65; bit += 8 {
				bs.Set(uint(bit))
			}
		case NORTHEAST:
			if !onRightHandSide(start) {
				for bit := start + 7; bit < 65; bit += 7 {
					bs.Set(uint(bit))
					if onRightHandSide(bit) {
						break
					}
				}
			}
		case EAST:
			if !onRightHandSide(start) {
				for bit := start - 1; bit > 0; bit-- {
					bs.Set(uint(bit))
					if onRightHandSide(bit) {
						break
					}
				}
			}
		case SOUTHEAST:
			if !onRightHandSide(start) {
				for bit := start - 9; bit > 0; bit -= 9 {
					bs.Set(uint(bit))
					if onRightHandSide(bit) {
						break
					}
				}
			}
		case SOUTH:
			for bit := start - 8; bit > 0; bit -= 8 {
				bs.Set(uint(bit))
			}
		case SOUTHWEST:
			if !onLeftHandSide(start) {
				for bit := start - 7; bit > 0; bit -= 7 {
					bs.Set(uint(bit))
					if onLeftHandSide(bit) {
						break
					}
				}
			}
		case WEST:
			if !onLeftHandSide(start) {
				for bit := start + 1; bit < 65; bit++ {
					bs.Set(uint(bit))
					if onLeftHandSide(bit) {
						break
					}
				}
			}
		case NORTHWEST:
			if !onLeftHandSide(start) {
				for bit := start + 9; bit < 65; bit += 9 {
					bs.Set(uint(bit))
					if onLeftHandSide(bit) {
						break
					}
				}
			}
		default:
			panic("oops!?")
		}
		ar[i] = bs
	}
	return ar
}

func onRightHandSide(bit int) bool { return bit%8 == 1 }
func onLeftHandSide(bit int) bool  { return bit%8 == 0 }

// AttackRays stores the squares that a bishop/queen on square x attacks
var AttackRays = [65]AttackRay{
	AttackRay{}, // first is empty
	// H1
	newAttackRay(1), newAttackRay(2), newAttackRay(3), newAttackRay(4), newAttackRay(5), newAttackRay(6), newAttackRay(7), newAttackRay(8),
	newAttackRay(9), newAttackRay(10), newAttackRay(11), newAttackRay(12), newAttackRay(13), newAttackRay(14), newAttackRay(15), newAttackRay(16),
	newAttackRay(17), newAttackRay(18), newAttackRay(19), newAttackRay(20), newAttackRay(21), newAttackRay(22), newAttackRay(23), newAttackRay(24),
	newAttackRay(25), newAttackRay(26), newAttackRay(27), newAttackRay(28), newAttackRay(29), newAttackRay(30), newAttackRay(31), newAttackRay(32),
	newAttackRay(33), newAttackRay(34), newAttackRay(35), newAttackRay(36), newAttackRay(37), newAttackRay(38), newAttackRay(39), newAttackRay(40),
	newAttackRay(41), newAttackRay(42), newAttackRay(43), newAttackRay(44), newAttackRay(45), newAttackRay(46), newAttackRay(47), newAttackRay(48),
	newAttackRay(49), newAttackRay(50), newAttackRay(51), newAttackRay(52), newAttackRay(53), newAttackRay(54), newAttackRay(55), newAttackRay(56),
	newAttackRay(57), newAttackRay(58), newAttackRay(59), newAttackRay(60), newAttackRay(61), newAttackRay(62), newAttackRay(63), newAttackRay(64),
}

// KnightAttackBitSets stores the squares that a knight on square x attacks
var KnightAttackBitSets = [65]bitset.BitSet{
	bitset.New(0), // first is empty
	// H1
	newKnightAttackBitset(1), newKnightAttackBitset(2), newKnightAttackBitset(3), newKnightAttackBitset(4), newKnightAttackBitset(5), newKnightAttackBitset(6),
	newKnightAttackBitset(7), newKnightAttackBitset(8), newKnightAttackBitset(9), newKnightAttackBitset(10), newKnightAttackBitset(11), newKnightAttackBitset(12),
	newKnightAttackBitset(13), newKnightAttackBitset(14), newKnightAttackBitset(15), newKnightAttackBitset(16), newKnightAttackBitset(17), newKnightAttackBitset(18),
	newKnightAttackBitset(19), newKnightAttackBitset(20), newKnightAttackBitset(21), newKnightAttackBitset(22), newKnightAttackBitset(23), newKnightAttackBitset(24),
	newKnightAttackBitset(25), newKnightAttackBitset(26), newKnightAttackBitset(27), newKnightAttackBitset(28), newKnightAttackBitset(29), newKnightAttackBitset(30),
	newKnightAttackBitset(31), newKnightAttackBitset(32), newKnightAttackBitset(33), newKnightAttackBitset(34), newKnightAttackBitset(35), newKnightAttackBitset(36),
	newKnightAttackBitset(37), newKnightAttackBitset(38), newKnightAttackBitset(39), newKnightAttackBitset(40), newKnightAttackBitset(41), newKnightAttackBitset(42),
	newKnightAttackBitset(43), newKnightAttackBitset(44), newKnightAttackBitset(45), newKnightAttackBitset(46), newKnightAttackBitset(47), newKnightAttackBitset(48),
	newKnightAttackBitset(49), newKnightAttackBitset(50), newKnightAttackBitset(51), newKnightAttackBitset(52), newKnightAttackBitset(53), newKnightAttackBitset(54),
	newKnightAttackBitset(55), newKnightAttackBitset(56), newKnightAttackBitset(57), newKnightAttackBitset(58), newKnightAttackBitset(59), newKnightAttackBitset(60),
	newKnightAttackBitset(61), newKnightAttackBitset(62), newKnightAttackBitset(63), newKnightAttackBitset(64),
}

// PawnAttackBitSets stores the squares where a pawn would attack square x (per colour)
// !!! this is different to all the others, it stores the squares which attack square x, not v.v. !!!
var PawnAttackBitSets = [][65]bitset.BitSet{
	{
		bitset.New(0), // first is empty
		// H1
		newPawnAttackBitset(colour.White, 1), newPawnAttackBitset(colour.White, 2), newPawnAttackBitset(colour.White, 3), newPawnAttackBitset(colour.White, 4), newPawnAttackBitset(colour.White, 5), newPawnAttackBitset(colour.White, 6),
		newPawnAttackBitset(colour.White, 7), newPawnAttackBitset(colour.White, 8), newPawnAttackBitset(colour.White, 9), newPawnAttackBitset(colour.White, 10), newPawnAttackBitset(colour.White, 11), newPawnAttackBitset(colour.White, 12),
		newPawnAttackBitset(colour.White, 13), newPawnAttackBitset(colour.White, 14), newPawnAttackBitset(colour.White, 15), newPawnAttackBitset(colour.White, 16), newPawnAttackBitset(colour.White, 17), newPawnAttackBitset(colour.White, 18),
		newPawnAttackBitset(colour.White, 19), newPawnAttackBitset(colour.White, 20), newPawnAttackBitset(colour.White, 21), newPawnAttackBitset(colour.White, 22), newPawnAttackBitset(colour.White, 23), newPawnAttackBitset(colour.White, 24),
		newPawnAttackBitset(colour.White, 25), newPawnAttackBitset(colour.White, 26), newPawnAttackBitset(colour.White, 27), newPawnAttackBitset(colour.White, 28), newPawnAttackBitset(colour.White, 29), newPawnAttackBitset(colour.White, 30),
		newPawnAttackBitset(colour.White, 31), newPawnAttackBitset(colour.White, 32), newPawnAttackBitset(colour.White, 33), newPawnAttackBitset(colour.White, 34), newPawnAttackBitset(colour.White, 35), newPawnAttackBitset(colour.White, 36),
		newPawnAttackBitset(colour.White, 37), newPawnAttackBitset(colour.White, 38), newPawnAttackBitset(colour.White, 39), newPawnAttackBitset(colour.White, 40), newPawnAttackBitset(colour.White, 41), newPawnAttackBitset(colour.White, 42),
		newPawnAttackBitset(colour.White, 43), newPawnAttackBitset(colour.White, 44), newPawnAttackBitset(colour.White, 45), newPawnAttackBitset(colour.White, 46), newPawnAttackBitset(colour.White, 47), newPawnAttackBitset(colour.White, 48),
		newPawnAttackBitset(colour.White, 49), newPawnAttackBitset(colour.White, 50), newPawnAttackBitset(colour.White, 51), newPawnAttackBitset(colour.White, 52), newPawnAttackBitset(colour.White, 53), newPawnAttackBitset(colour.White, 54),
		newPawnAttackBitset(colour.White, 55), newPawnAttackBitset(colour.White, 56), newPawnAttackBitset(colour.White, 57), newPawnAttackBitset(colour.White, 58), newPawnAttackBitset(colour.White, 59), newPawnAttackBitset(colour.White, 60),
		newPawnAttackBitset(colour.White, 61), newPawnAttackBitset(colour.White, 62), newPawnAttackBitset(colour.White, 63), newPawnAttackBitset(colour.White, 64),
	},
	{
		bitset.New(0), // first is empty
		// H1
		newPawnAttackBitset(colour.Black, 1), newPawnAttackBitset(colour.Black, 2), newPawnAttackBitset(colour.Black, 3), newPawnAttackBitset(colour.Black, 4), newPawnAttackBitset(colour.Black, 5), newPawnAttackBitset(colour.Black, 6),
		newPawnAttackBitset(colour.Black, 7), newPawnAttackBitset(colour.Black, 8), newPawnAttackBitset(colour.Black, 9), newPawnAttackBitset(colour.Black, 10), newPawnAttackBitset(colour.Black, 11), newPawnAttackBitset(colour.Black, 12),
		newPawnAttackBitset(colour.Black, 13), newPawnAttackBitset(colour.Black, 14), newPawnAttackBitset(colour.Black, 15), newPawnAttackBitset(colour.Black, 16), newPawnAttackBitset(colour.Black, 17), newPawnAttackBitset(colour.Black, 18),
		newPawnAttackBitset(colour.Black, 19), newPawnAttackBitset(colour.Black, 20), newPawnAttackBitset(colour.Black, 21), newPawnAttackBitset(colour.Black, 22), newPawnAttackBitset(colour.Black, 23), newPawnAttackBitset(colour.Black, 24),
		newPawnAttackBitset(colour.Black, 25), newPawnAttackBitset(colour.Black, 26), newPawnAttackBitset(colour.Black, 27), newPawnAttackBitset(colour.Black, 28), newPawnAttackBitset(colour.Black, 29), newPawnAttackBitset(colour.Black, 30),
		newPawnAttackBitset(colour.Black, 31), newPawnAttackBitset(colour.Black, 32), newPawnAttackBitset(colour.Black, 33), newPawnAttackBitset(colour.Black, 34), newPawnAttackBitset(colour.Black, 35), newPawnAttackBitset(colour.Black, 36),
		newPawnAttackBitset(colour.Black, 37), newPawnAttackBitset(colour.Black, 38), newPawnAttackBitset(colour.Black, 39), newPawnAttackBitset(colour.Black, 40), newPawnAttackBitset(colour.Black, 41), newPawnAttackBitset(colour.Black, 42),
		newPawnAttackBitset(colour.Black, 43), newPawnAttackBitset(colour.Black, 44), newPawnAttackBitset(colour.Black, 45), newPawnAttackBitset(colour.Black, 46), newPawnAttackBitset(colour.Black, 47), newPawnAttackBitset(colour.Black, 48),
		newPawnAttackBitset(colour.Black, 49), newPawnAttackBitset(colour.Black, 50), newPawnAttackBitset(colour.Black, 51), newPawnAttackBitset(colour.Black, 52), newPawnAttackBitset(colour.Black, 53), newPawnAttackBitset(colour.Black, 54),
		newPawnAttackBitset(colour.Black, 55), newPawnAttackBitset(colour.Black, 56), newPawnAttackBitset(colour.Black, 57), newPawnAttackBitset(colour.Black, 58), newPawnAttackBitset(colour.Black, 59), newPawnAttackBitset(colour.Black, 60),
		newPawnAttackBitset(colour.Black, 61), newPawnAttackBitset(colour.Black, 62), newPawnAttackBitset(colour.Black, 63), newPawnAttackBitset(colour.Black, 64),
	},
}

// KingAttackBitSets stores the squares that a king on square x attacks
var KingAttackBitSets = [65]bitset.BitSet{
	bitset.New(0), // first is empty
	// H1
	newKingAttackBitset(1), newKingAttackBitset(2), newKingAttackBitset(3), newKingAttackBitset(4), newKingAttackBitset(5), newKingAttackBitset(6),
	newKingAttackBitset(7), newKingAttackBitset(8), newKingAttackBitset(9), newKingAttackBitset(10), newKingAttackBitset(11), newKingAttackBitset(12),
	newKingAttackBitset(13), newKingAttackBitset(14), newKingAttackBitset(15), newKingAttackBitset(16), newKingAttackBitset(17), newKingAttackBitset(18),
	newKingAttackBitset(19), newKingAttackBitset(20), newKingAttackBitset(21), newKingAttackBitset(22), newKingAttackBitset(23), newKingAttackBitset(24),
	newKingAttackBitset(25), newKingAttackBitset(26), newKingAttackBitset(27), newKingAttackBitset(28), newKingAttackBitset(29), newKingAttackBitset(30),
	newKingAttackBitset(31), newKingAttackBitset(32), newKingAttackBitset(33), newKingAttackBitset(34), newKingAttackBitset(35), newKingAttackBitset(36),
	newKingAttackBitset(37), newKingAttackBitset(38), newKingAttackBitset(39), newKingAttackBitset(40), newKingAttackBitset(41), newKingAttackBitset(42),
	newKingAttackBitset(43), newKingAttackBitset(44), newKingAttackBitset(45), newKingAttackBitset(46), newKingAttackBitset(47), newKingAttackBitset(48),
	newKingAttackBitset(49), newKingAttackBitset(50), newKingAttackBitset(51), newKingAttackBitset(52), newKingAttackBitset(53), newKingAttackBitset(54),
	newKingAttackBitset(55), newKingAttackBitset(56), newKingAttackBitset(57), newKingAttackBitset(58), newKingAttackBitset(59), newKingAttackBitset(60),
	newKingAttackBitset(61), newKingAttackBitset(62), newKingAttackBitset(63), newKingAttackBitset(64),
}

// AttacksOnEnpassantSquares stores bitsets of the squares which attack the 8 enpassant squares (A6..H6 for white moves, A3..H3 for black
// Usage: AttacksOnEnpassantSquares[colour.White][4] returns a bitset of the (two) squares which attack D6
// Runs from file 1 to file 8 ([1] to [8]
var AttacksOnEnpassantSquares = [2][9]bitset.BitSet{{bitset.New(0), bitset.New(0x4000000000), bitset.New(0xA000000000), bitset.New(0x5000000000),
	bitset.New(0x2800000000), bitset.New(0x1400000000), bitset.New(0x0A00000000), bitset.New(0x0500000000), bitset.New(0x0200000000)},
	{bitset.New(0), bitset.New(0x40000000), bitset.New(0xA0000000), bitset.New(0x50000000), bitset.New(0x28000000),
		bitset.New(0x14000000), bitset.New(0x0A000000), bitset.New(0x05000000), bitset.New(0x02000000)},
}
