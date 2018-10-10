package ray

import (
	"testing"

	"github.com/rjo67/chess/bitset"
)

func TestAttackRays(t *testing.T) {
	ar := AttackRays[1]
	checkBitSet(ar[NORTH], []uint{9, 17, 25, 33, 41, 49, 57}, t)
	checkBitSet(ar[EAST], []uint{}, t)
	ar = AttackRays[12]
	checkBitSet(ar[WEST], []uint{13, 14, 15, 16}, t)
	checkBitSet(ar[NORTHWEST], []uint{21, 30, 39, 48}, t)
	checkBitSet(ar[SOUTHWEST], []uint{5}, t)
	checkBitSet(ar[SOUTHEAST], []uint{3}, t)
	ar = AttackRays[56]
	checkBitSet(ar[WEST], []uint{}, t)
	checkBitSet(ar[NORTH], []uint{64}, t)
	checkBitSet(ar[SOUTH], []uint{48, 40, 32, 24, 16, 8}, t)
	checkBitSet(ar[NORTHEAST], []uint{63}, t)
}

func TestKnightAttackBitSet(t *testing.T) {
	checkBitSet(KnightAttackBitSets[1], []uint{11, 18}, t)
	checkBitSet(KnightAttackBitSets[64], []uint{54, 47}, t)
}

func TestKingAttackBitSet(t *testing.T) {
	checkBitSet(KingAttackBitSets[1], []uint{2, 10, 9}, t)
	checkBitSet(KingAttackBitSets[64], []uint{63, 56, 55}, t)
}

// tests the given bitset by creating a new one with setBits,
// and ORing the two together. The result should be the same bitset value
func checkBitSet(bs bitset.BitSet, setBits []uint, t *testing.T) {
	newBitset := bitset.New(0)
	for _, bit := range setBits {
		newBitset.Set(bit)
	}
	result := bs.Or(newBitset)
	if bs.Val() != result.Val() {
		t.Errorf("differing bitsets. BS1:\n%s\n, BS2:\n%s", bs.String(), newBitset.String())
	}
}
