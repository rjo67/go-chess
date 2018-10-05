package bitset

import (
	"testing"

	"github.com/rjo67/chess/square"
)

func TestBitset(t *testing.T) {
	var b [8]byte
	b[0] = 0x80
	b[1] = 0x40
	b[2] = 0x20
	b[3] = 0x10
	b[4] = 0x08
	b[5] = 0x04
	b[6] = 0x02
	b[7] = 0x01
	bs := NewFromByteArray(b)
	if bs.Val != 72624976668147840 {
		t.Errorf("expected 72624976668147840 but got %d for bitset\n%s", bs.Val, bs.ToString())
	}
}

func TestToString(t *testing.T) {
	bs := BitSet{255}
	str := bs.ToString()
	if str != "00000000\n00000000\n00000000\n00000000\n00000000\n00000000\n00000000\n11111111\n" {
		t.Errorf("got bad string: %s for hex value %x\n", str, bs.Val)
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
		bs := BitSet{Val: 0} // start with empty bitset
		for _, setBit := range setBits {
			bs.Set(setBit)
		}
		checkBits(t, bs, setBits, true)
	}
	bs := BitSet{Val: 0}
	bs.Set(4).Set(25).Set(44)
	val1 := bs.Val
	bs2 := BitSet{Val: 0}
	bs2.SetSquare(square.E1).SetSquare(square.H4).SetSquare(square.E6)
	val2 := bs2.Val
	if val1 != val2 {
		t.Errorf("got different values for the bitsets. Bitset 1:\n%s\n, Bitset 2:\n%s", bs.ToString(), bs2.ToString())
	}
	bs3 := NewFromSquares(square.E1, square.H4, square.E6)
	val3 := bs3.Val
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
		bs := BitSet{Val: 0xFFFFFFFFFFFFFFFF} // start with full bitset
		for _, setBit := range setBits {
			bs.Clear(setBit)
		}
		checkBits(t, bs, setBits, false)
	}
	bs := BitSet{Val: 0xFFFFFFFFFFFFFFFF}
	bs.Clear(4).Clear(25).Clear(44)
	val1 := bs.Val
	bs2 := BitSet{Val: 0xFFFFFFFFFFFFFFFF}
	bs2.ClearSquare(square.E1).ClearSquare(square.H4).ClearSquare(square.E6)
	if val1 != bs2.Val {
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
