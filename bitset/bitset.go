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
	Val uint64
}

// New is a convenience constructor to create a Bitset from an array of bytes.
// The byte array is processed in reverse order (from [7] down to [0]),
// i.e. input[0] is the bottom row
func New(input [8]byte) BitSet {
	var bs uint64
	count := uint64(8 * (len(input) - 1))
	for index := len(input) - 1; index >= 0; index-- {
		b := input[index]
		var mult uint64 = 1 << count
		nb := uint64(b) * mult
		bs += nb
		count -= 8
	}
	return BitSet{Val: bs}
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
	return bs.Set(uint(sq) + 1) // squares run 0..63
}

// Set sets the given bit-position. Returns itself to allow chaining.
// posn identifies the required bit, running from 1 (bottom right) to 64 (top left)
func (bs *BitSet) Set(posn uint) *BitSet {
	if posn < 1 || posn > 64 {
		panic("invalid value for posn: " + fmt.Sprint(posn))
	}
	var mask uint64 = 1 << (posn - 1)
	if bs.Val&mask != mask {
		bs.Val += mask
	}
	return bs
}

// ClearSquare clears the bit at the given square
func (bs *BitSet) ClearSquare(sq square.Square) *BitSet {
	return bs.Clear(uint(sq) + 1) // squares run 0..63
}

// Clear unsets the given bit-position. Returns itself to allow chaining.
// posn identifies the required bit, running from 1 (bottom right) to 64 (top left)
func (bs *BitSet) Clear(posn uint) *BitSet {
	if posn < 1 || posn > 64 {
		panic("invalid value for posn: " + fmt.Sprint(posn))
	}
	var mask uint64 = 1 << (posn - 1)
	if bs.Val&mask == mask {
		bs.Val -= mask
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
	return bs.Val&mask == mask
}
