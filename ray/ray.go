package ray

import (
	"github.com/rjo67/chess/bitset"
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

// AttackRay stores, for each square, the squares which are potentially attacked by a piece on that square
type AttackRay [8]bitset.BitSet

// internal constructor
func newAttackRay(start int) AttackRay {
	ar := AttackRay{}
	for i, rayType := range AllDirections {
		var bs bitset.BitSet
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

// AttackRays stores the attack rays for each square of the board
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
