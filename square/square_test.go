package square

import (
	"testing"
)

func TestRank(t *testing.T) {
	data := []struct {
		expected uint32
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
			t.Errorf("sq: %s, expected %d but got %d", test.sq.ToString(), test.expected, rank)
		}
	}
}

func TestFile(t *testing.T) {
	data := []struct {
		expected uint32
		sq       Square
	}{
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
			t.Errorf("sq: %s, expected %d but got %d", test.sq.ToString(), test.expected, file)
		}
	}
}

func TestToString(t *testing.T) {
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
		gotString := test.sq.ToString()
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
		gotSquare := FromString(test.squareAsString)
		if gotSquare != test.expected {
			t.Errorf("string %s: expected %v but got %v", test.squareAsString, test.expected, gotSquare)
		}
	}
}
