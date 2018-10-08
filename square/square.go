package square

import (
	"fmt"
	"strings"
)

// Square represents a square of the chessboard
type Square uint32

// the squares of the board, from H1..A8. Numerically 1..64
const (
	_ Square = iota // skip 0-value
	H1
	G1
	F1
	E1
	D1
	C1
	B1
	A1
	H2
	G2
	F2
	E2
	D2
	C2
	B2
	A2
	H3
	G3
	F3
	E3
	D3
	C3
	B3
	A3
	H4
	G4
	F4
	E4
	D4
	C4
	B4
	A4
	H5
	G5
	F5
	E5
	D5
	C5
	B5
	A5
	H6
	G6
	F6
	E6
	D6
	C6
	B6
	A6
	H7
	G7
	F7
	E7
	D7
	C7
	B7
	A7
	H8
	G8
	F8
	E8
	D8
	C8
	B8
	A8
)

// String returns the algebraic representation of the square
func (sq Square) String() string {
	return convertToFileLetter(sq.File()) + fmt.Sprint(sq.Rank())
}

func convertToFileLetter(file uint32) string {
	return string('A' - 1 + file)
}

// FromString returns the square matching the algebraic representation
func FromString(str string) (Square, error) {
	if len(str) != 2 {
		return Square(0), fmt.Errorf("unrecognised square '%s'", str)
	}
	upperStr := strings.ToUpper(str)
	if upperStr[1] < '1' || upperStr[1] > '8' || upperStr[0] < 'A' || upperStr[0] > 'H' {
		return Square(0), fmt.Errorf("unrecognised square '%s'", str)
	}
	rank := upperStr[1] - '1'
	file := 7 - (upperStr[0] - 'A')
	return Square(rank*8 + file + 1), nil // "+1" since squares are 1..64
}

// Rank of the square (1..8)
// e.g. H1 has rank 1, A4 has rank 4
func (sq Square) Rank() uint32 {
	return uint32((sq-1)/8) + 1
}

// File of the square (1..8). Can be converted with convertToFileLetter() into A..H
// e.g. H1 has file 8, A4 has file 1
func (sq Square) File() uint32 {
	return uint32(8 - (sq-1)%8)
}
