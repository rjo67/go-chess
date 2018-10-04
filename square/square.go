package square

import (
	"fmt"
	"strings"
)

// Square represents a square of the chessboard
type Square uint32

// the squares of the board, from H1..A8
const (
	H1 Square = iota
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

// ToString returns the algebraic representation of the square
func (sq Square) ToString() string {
	return convertToFileLetter(sq.File()) + fmt.Sprint(sq.Rank()+1)
}

func convertToFileLetter(file uint32) string {
	return string('A' + (7 - file))
}

// FromString returns the square matching the algebraic representation
func FromString(str string) Square {
	upperStr := strings.ToUpper(str)
	rank := upperStr[1] - '1'
	file := 7 - (upperStr[0] - 'A')
	return Square(rank*8 + file)
}

// Rank of the square (0..7) TODO maybe better to have 1..8?
func (sq Square) Rank() uint32 {
	return uint32(sq / 8)
}

// File of the square (0..7). Can be converted with convertToFileLetter() into A..H
func (sq Square) File() uint32 {
	return uint32(sq % 8)
}
