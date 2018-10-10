package move

import (
	"testing"

	"github.com/rjo67/chess/bitset"
	"github.com/rjo67/chess/ray"
	"github.com/rjo67/chess/square"
)

func TestSearch(t *testing.T) {
	occupiedSquares := bitset.NewFromByteArray([8]byte{0x00, 0x00, 0x40, 0x00, 0x20, 0x80, 0x02, 0x10})
	/*
	 00010000
	 00000010
	 10000000
	 00100000
	 00000000
	 01000000
	 00000000
	 00000000
	*/

	data := []struct {
		startSquare     int
		direction       ray.Direction
		expectedSetBits []int
	}{
		{int(square.A4), ray.NORTH, []int{40, 48}},
		{int(square.A4), ray.NORTH, []int{40, 48}},
		{int(square.D1), ray.NORTHWEST, []int{14, 23}},
		{int(square.G5), ray.WEST, []int{35, 36, 37, 38}},
		{int(square.G8), ray.SOUTHWEST, []int{51, 44, 37, 30, 23}},
		{int(square.C8), ray.SOUTH, []int{54, 46, 38}},
		{int(square.A7), ray.SOUTHEAST, []int{47, 38}},
		{int(square.A7), ray.EAST, []int{55, 54, 53, 52, 51, 50}},
		{int(square.E5), ray.NORTHEAST, []int{43, 50}},
	}

	for testNbr, i := range data {
		bs, _ := Search2(i.startSquare, i.direction, occupiedSquares)
		if len(i.expectedSetBits) != bs.Cardinality() {
			t.Errorf("test %d: expected %d set-bits, got %d for bitset:\n%s", testNbr, len(i.expectedSetBits), bs.Cardinality(), bs.String())
		} else {
			checkBits(testNbr, bs, i.expectedSetBits, t)
		}
	}
}

func checkBits(testNbr int, bs bitset.BitSet, expectedSetBits []int, t *testing.T) {
	errors := make([]int, 0, 5)
	for _, bit := range expectedSetBits {
		if !bs.IsSet(uint(bit)) {
			errors = append(errors, bit)
		}
	}
	if len(errors) != 0 {
		t.Errorf("test %d: found %d errors (%v) for bitset:\n%s", testNbr, len(errors), errors, bs.String())
	}
}
