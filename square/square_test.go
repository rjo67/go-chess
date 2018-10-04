package square

import (
	"testing"
)

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
