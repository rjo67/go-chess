package bitset

import (
	"testing"

	"github.com/rjo67/chess/square"
)

func TestBitset(t *testing.T) {
	bs := NewFromByteArray([8]byte{0x80, 0x40, 0x20, 0x10, 0x08, 0x04, 0x02, 0x01})
	if bs.Val() != 72624976668147840 {
		t.Errorf("expected 72624976668147840 but got %d for bitset\n%s", bs.Val(), bs.ToString())
	}
}

func TestToString(t *testing.T) {
	bs := BitSet{255}
	str := bs.ToString()
	if str != "00000000\n00000000\n00000000\n00000000\n00000000\n00000000\n00000000\n11111111\n" {
		t.Errorf("got bad string: %s for hex value %x\n", str, bs.Val())
	}

}
func TestIsSet(t *testing.T) {
	var b [8]byte
	b[7] = 0x01
	b[6] = 0x02
	b[5] = 0x04
	b[4] = 0x08
	b[3] = 0x10
	b[2] = 0x20
	b[1] = 0x40
	b[0] = 0x80
	bs := NewFromByteArray(b)
	checkBits(t, bs, []uint{8, 15, 22, 29, 36, 43, 50, 57}, true)
}

func TestSet(t *testing.T) {
	data := [][]uint{
		{1, 7},
		{5, 16, 21, 35},
		{3, 22, 44, 64},
	}
	for _, setBits := range data {
		bs := BitSet{} // start with empty bitset
		for _, setBit := range setBits {
			bs.Set(setBit)
		}
		checkBits(t, bs, setBits, true)
	}
	bs := BitSet{}
	bs.Set(4).Set(25).Set(44)
	val1 := bs.Val()
	bs2 := BitSet{}
	bs2.SetSquare(square.E1).SetSquare(square.H4).SetSquare(square.E6)
	val2 := bs2.Val()
	if val1 != val2 {
		t.Errorf("got different values for the bitsets. Bitset 1:\n%s\n, Bitset 2:\n%s", bs.ToString(), bs2.ToString())
	}
	bs3 := NewFromSquares(square.E1, square.H4, square.E6)
	val3 := bs3.Val()
	if val2 != val3 {
		t.Errorf("got different values for the bitsets. Bitset 2:\n%s\n, Bitset 3:\n%s", bs2.ToString(), bs3.ToString())
	}
}

func TestClear(t *testing.T) {
	data := [][]uint{
		{1, 7},
		{5, 16, 21, 35},
		{3, 22, 44, 64},
	}
	for _, setBits := range data {
		bs := BitSet{val: 0xFFFFFFFFFFFFFFFF} // start with full bitset
		for _, setBit := range setBits {
			bs.Clear(setBit)
		}
		checkBits(t, bs, setBits, false)
	}
	bs := BitSet{val: 0xFFFFFFFFFFFFFFFF}
	bs.Clear(4).Clear(25).Clear(44)
	val1 := bs.Val()
	bs2 := BitSet{val: 0xFFFFFFFFFFFFFFFF}
	bs2.ClearSquare(square.E1).ClearSquare(square.H4).ClearSquare(square.E6)
	if val1 != bs2.Val() {
		t.Errorf("got different values for the bitsets. Bitset 1:\n%s\n, Bitset 2:\n%s", bs.ToString(), bs2.ToString())
	}
}

func TestOr(t *testing.T) {
	bs1 := NewFromSquares(square.E1, square.H8)
	bs2 := NewFromSquares(square.A8, square.F2)
	bs3 := bs1.Or(bs2)
	// make sure bs1 and bs2 weren't affected
	checkBits(t, bs1, []uint{uint(square.E1), uint(square.H8)}, true)
	checkBits(t, bs2, []uint{uint(square.A8), uint(square.F2)}, true)
	// and check result of tne 'or'
	checkBits(t, bs3, []uint{uint(square.E1), uint(square.H8), uint(square.A8), uint(square.F2)}, true)
}

func TestAnd(t *testing.T) {
	bs1 := NewFromSquares(square.E1, square.H8)
	bs2 := NewFromSquares(square.E1, square.F2)
	bs3 := bs1.And(bs2)
	// make sure bs1 and bs2 weren't affected
	checkBits(t, bs1, []uint{uint(square.E1), uint(square.H8)}, true)
	checkBits(t, bs2, []uint{uint(square.E1), uint(square.F2)}, true)
	// and check result of tne 'and'
	checkBits(t, bs3, []uint{uint(square.E1)}, true)
}

func TestXor(t *testing.T) {
	bs1 := NewFromSquares(square.E1, square.H8)
	bs2 := NewFromSquares(square.E1, square.F2)
	bs3 := bs1.Xor(bs2)
	// make sure bs1 and bs2 weren't affected
	checkBits(t, bs1, []uint{uint(square.E1), uint(square.H8)}, true)
	checkBits(t, bs2, []uint{uint(square.E1), uint(square.F2)}, true)
	// and check result of tne operation
	checkBits(t, bs3, []uint{uint(square.F2), uint(square.H8)}, true)
}

func TestCardinality(t *testing.T) {
	data := [][]uint{
		{},
		{1, 7},
		{5, 16, 21, 35, 47, 48, 49, 55},
	}
	for _, setBits := range data {
		bs := BitSet{} // start with empty bitset
		for _, setBit := range setBits {
			bs.Set(setBit)
		}
		if bs.Cardinality() != len(setBits) {
			t.Errorf("expected cardinality %d but got %d for bitset:\n%s", len(setBits), bs.Cardinality(), bs.ToString())
		}
	}
}

func TestRank(t *testing.T) {
	checkBits(t, Rank(1), []uint{1, 2, 3, 4, 5, 6, 7, 8}, true)
	checkBits(t, Rank(2), []uint{9, 10, 11, 12, 13, 14, 15, 16}, true)
	checkBits(t, Rank(3), []uint{17, 18, 19, 20, 21, 22, 23, 24}, true)
	checkBits(t, Rank(4), []uint{25, 26, 27, 28, 29, 30, 31, 32}, true)
	checkBits(t, Rank(5), []uint{33, 34, 35, 36, 37, 38, 39, 40}, true)
	checkBits(t, Rank(6), []uint{41, 42, 43, 44, 45, 46, 47, 48}, true)
	checkBits(t, Rank(7), []uint{49, 50, 51, 52, 53, 54, 55, 56}, true)
	checkBits(t, Rank(8), []uint{57, 58, 59, 60, 61, 62, 63, 64}, true)
}

func TestFile(t *testing.T) {
	checkBits(t, File(1), []uint{1, 9, 17, 25, 33, 41, 49, 57}, true)
	checkBits(t, File(2), []uint{2, 10, 18, 26, 34, 42, 50, 58}, true)
	checkBits(t, File(3), []uint{3, 11, 19, 27, 35, 43, 51, 59}, true)
	checkBits(t, File(4), []uint{4, 12, 20, 28, 36, 44, 52, 60}, true)
	checkBits(t, File(5), []uint{5, 13, 21, 29, 37, 45, 53, 61}, true)
	checkBits(t, File(6), []uint{6, 14, 22, 30, 38, 46, 54, 62}, true)
	checkBits(t, File(7), []uint{7, 15, 23, 31, 39, 47, 55, 63}, true)
	checkBits(t, File(8), []uint{8, 16, 24, 32, 40, 48, 56, 64}, true)
}

// helper routine. Checks that all required bits are set, and all others are not set
func checkBits(t *testing.T, bs BitSet, setBits []uint, checkIfSet bool) {
	for _, bit := range setBits {
		if checkIfSet {
			if !bs.IsSet(uint(bit)) {
				t.Errorf("bit %d should be set for bitset\n%s", bit, bs.ToString())
			}
		} else {
			if bs.IsSet(uint(bit)) {
				t.Errorf("bit %d should not be set for bitset\n%s", bit, bs.ToString())
			}
		}
	}
	for i := 1; i < 65; i++ {
		// test all bit positions that aren't in 'setBits'
		var found = false
		for _, bit := range setBits {
			if uint(i) == bit {
				found = true
				break
			}
		}
		if !found {
			if checkIfSet {
				if bs.IsSet(uint(i)) {
					t.Errorf("bit %d should not be set for bitset\n%s", i, bs.ToString())
				}
			} else {
				if !bs.IsSet(uint(i)) {
					t.Errorf("bit %d should be set for bitset\n%s", i, bs.ToString())
				}
			}
		}
	}
}
