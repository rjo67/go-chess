package bitset

import (
	"fmt"
	"strings"

	"github.com/rjo67/chess/square"
)

// BitSet is a bitset for chess purposes.
//
// Can be created by passing an int value directly, using New with an array of bytes,
// or by setting single bits.
//
// Bit 1 is bottom right, Bit 64 is top left
type BitSet struct {
	val uint64
}

// NewFromByteArray is a convenience constructor to create a Bitset from an array of bytes.
// The byte array is processed in reverse order (from [7] down to [0]),
// i.e. input[0] is the bottom row
func NewFromByteArray(input [8]byte) BitSet {
	var bs uint64
	count := uint64(8 * (len(input) - 1))
	for index := len(input) - 1; index >= 0; index-- {
		b := input[index]
		var mult uint64 = 1 << count
		nb := uint64(b) * mult
		bs += nb
		count -= 8
	}
	return BitSet{val: bs}
}

// NewFromSquares is a convenience constructor to create a Bitset from a list of squares.
// The corresponding bits in the Bitset will be set.
func NewFromSquares(squares ...square.Square) BitSet {
	bs := BitSet{0}
	for _, square := range squares {
		bs.SetSquare(square)
	}
	return bs
}

// Val returns the value of the bitset
func (bs BitSet) Val() uint64 {
	return bs.val
}

// And returns a new bitset resulting from the logical AND of the current bitset and the supplied bitset
func (bs BitSet) And(other BitSet) BitSet {
	return BitSet{bs.val & other.val}
}

// Or returns a new bitset resulting from the logical OR of the current bitset and the supplied bitset
func (bs BitSet) Or(other BitSet) BitSet {
	return BitSet{bs.val | other.val}
}

// Cardinality returns the number of set-bits in the bitset
func (bs BitSet) Cardinality() int {
	count := 0
	for i := 1; i < 65; i++ {
		if bs.IsSet(uint(i)) {
			count++
		}
	}
	return count
}

// ToString returns a visual representation of the bitset in 8 rows of 8
func (bs BitSet) ToString() string {
	var posn uint = 65
	var sb strings.Builder
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			posn--
			if bs.IsSet(posn) {
				sb.WriteString("1")
			} else {
				sb.WriteString("0")
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

// SetSquare sets the bit at the given square
func (bs *BitSet) SetSquare(sq square.Square) *BitSet {
	return bs.Set(uint(sq))
}

// Set sets the given bit-position. Returns itself to allow chaining.
// posn identifies the required bit, running from 1 (bottom right) to 64 (top left)
func (bs *BitSet) Set(posn uint) *BitSet {
	if posn < 1 || posn > 64 {
		panic("invalid value for posn: " + fmt.Sprint(posn))
	}
	var mask uint64 = 1 << (posn - 1)
	if bs.val&mask != mask {
		bs.val += mask
	}
	return bs
}

// SetBits returns a slice containing all set-bits
func (bs BitSet) SetBits() []int {
	var squares [64]int
	slot := 0
	for i := 1; i < 65; i++ {
		if bs.IsSet(uint(i)) {
			squares[slot] = i
			slot++
		}
	}
	return squares[0:slot]
}

// ClearSquare clears the bit at the given square
func (bs *BitSet) ClearSquare(sq square.Square) *BitSet {
	return bs.Clear(uint(sq))
}

// Clear unsets the given bit-position. Returns itself to allow chaining.
// posn identifies the required bit, running from 1 (bottom right) to 64 (top left)
func (bs *BitSet) Clear(posn uint) *BitSet {
	if posn < 1 || posn > 64 {
		panic("invalid value for posn: " + fmt.Sprint(posn))
	}
	var mask uint64 = 1 << (posn - 1)
	if bs.val&mask == mask {
		bs.val -= mask
	}
	return bs
}

// IsSet returns true if the given bit-position is set
// posn identifies the required bit, running from 1 (bottom right) to 64 (top left)
func (bs BitSet) IsSet(posn uint) bool {
	if posn < 1 || posn > 64 {
		panic("invalid value for posn: " + fmt.Sprint(posn))
	}
	var mask uint64 = 1 << (posn - 1)
	return bs.val&mask == mask
}

// File creates a bitset containing bits of the nth file (1..8)
func File(n int) BitSet {
	if n < 1 || n > 8 {
		panic(fmt.Sprintf("invalid value for file: %d", n))
	}
	offset := uint(n - 1)
	bs := BitSet{0}
	bs.Set(1 + offset).Set(9 + offset).Set(17 + offset).Set(25 + offset).Set(33 + offset).Set(41 + offset).Set(49 + offset).Set(57 + offset)
	return bs
}

// Rank creates a bitset containing bits of the nth rank (1..8)
func Rank(n int) BitSet {
	if n < 1 || n > 8 {
		panic(fmt.Sprintf("invalid value for rank: %d", n))
	}
	offset := uint((n - 1) * 8)
	bs := BitSet{0}
	bs.Set(1 + offset).Set(2 + offset).Set(3 + offset).Set(4 + offset).Set(5 + offset).Set(6 + offset).Set(7 + offset).Set(8 + offset)
	return bs
}
