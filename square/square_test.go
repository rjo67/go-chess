package square

import (
	"testing"
)

func TestRank(t *testing.T) {
	data := []struct {
		expected int
		sq       Square
	}{
		{1, H1},
		{2, A2},
		{3, A3},
		{4, D4},
		{5, H5},
		{6, B6},
		{7, C7},
		{8, A8},
	}
	for _, test := range data {
		rank := test.sq.Rank()
		if rank != test.expected {
			t.Errorf("sq: %s, expected %d but got %d", test.sq.String(), test.expected, rank)
		}
	}
}

func TestFile(t *testing.T) {
	data := []struct {
		expected int
		sq       Square
	}{
		{8, H1},
		{1, A2},
		{2, B6},
		{3, C7},
		{4, D4},
		{5, E8},
		{6, F2},
		{7, G7},
		{8, H8},
	}
	for _, test := range data {
		file := test.sq.File()
		if file != test.expected {
			t.Errorf("sq: %s, expected %d but got %d", test.sq.String(), test.expected, file)
		}
	}
}

func TestAdjacent(t *testing.T) {
	data := []struct {
		sq       Square
		expected []int
	}{
		{H1, []int{1, 2, 9, 10}},
		{G1, []int{1, 2, 3, 9, 10, 11}},
		{F1, []int{2, 3, 4, 10, 11, 12}},
		{E1, []int{3, 4, 5, 11, 12, 13}},
		{D1, []int{4, 5, 6, 12, 13, 14}},
		{C1, []int{5, 6, 7, 13, 14, 15}},
		{B1, []int{6, 7, 8, 14, 15, 16}},
		{A1, []int{7, 8, 15, 16}},
		{H3, []int{17, 18, 9, 10, 25, 26}},
		{D3, []int{21, 20, 22, 12, 13, 14, 28, 29, 30}},
		{A3, []int{24, 16, 32, 15, 23, 31}},
	}
	for _, test := range data {
		squares := make([]int, 0, 10)
		for otherSq := 1; otherSq < 65; otherSq++ {
			if test.sq.IsAdjacentTo(Square(otherSq)) {
				squares = append(squares, otherSq)
			}
		}
		if len(squares) != len(test.expected) {
			t.Errorf("sq: %d, expected len %d but got %d (%v)", test.sq, len(test.expected), len(squares), squares)
		} else {
			for _, sq := range test.expected {
				// find sq in squares
				for i, sq2 := range squares {
					if sq2 == sq {
						squares = append(squares[:i], squares[i+1:]...)
						break
					}
				}
			}
			if len(squares) != 0 {
				t.Errorf("sq: %d, squares left over: %v", test.sq, squares)
			}
		}
	}
}

func TestString(t *testing.T) {
	data := []struct {
		expected string
		sq       Square
	}{
		{"H1", H1},
		{"A2", A2},
		{"C7", C7},
		{"D3", D3},
		{"A8", A8},
	}
	for _, test := range data {
		gotString := test.sq.String()
		if gotString != test.expected {
			t.Errorf("expected %s but got %s", test.expected, gotString)
		}
	}
}

func TestFromString(t *testing.T) {
	data := []struct {
		squareAsString string
		expected       Square
	}{
		{"H1", H1},
		{"A2", A2},
		{"C7", C7},
		{"D3", D3},
		{"A8", A8},
		{"H8", H8},
		{"E8", E8},
	}
	for _, test := range data {
		gotSquare, _ := FromString(test.squareAsString)
		if gotSquare != test.expected {
			t.Errorf("string %s: expected %v but got %v", test.squareAsString, test.expected, gotSquare)
		}
	}
}
