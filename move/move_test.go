package move

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/rjo67/chess/bitset"
	"github.com/rjo67/chess/piece"
	"github.com/rjo67/chess/piece/colour"
	"github.com/rjo67/chess/ray"
	"github.com/rjo67/chess/square"
)

func TestSize(t *testing.T) {
	m := New(colour.White, square.A5, square.B7, piece.KNIGHT)
	fmt.Printf("size of move: %d bytes\n", reflect.TypeOf(m).Size())
}

func TestMove(t *testing.T) {
	m := New(colour.White, square.A5, square.B7, piece.KNIGHT)
	if m.From() != square.A5 {
		t.Fatalf("'from' square incorrect: " + m.String())
	}
	if m.To() != square.B7 {
		t.Fatalf("'to' square incorrect: " + m.String())
	}
	m = New(colour.White, square.H1, square.A8, piece.ROOK)
	if m.From() != square.H1 {
		t.Fatalf("'from' square incorrect: " + m.String())
	}
	if m.To() != square.A8 {
		t.Fatalf("'to' square incorrect: " + m.String())
	}
}

func TestPromotion(t *testing.T) {
	m := NewPromotion(colour.White, square.A7, square.A8, piece.KNIGHT)
	if m.PromotedPiece() != piece.KNIGHT {
		t.Fatalf("promoted piece incorrect: " + m.String())
	}
	m = NewPromotion(colour.White, square.A7, square.A8, piece.QUEEN)
	if m.PromotedPiece() != piece.QUEEN {
		t.Fatalf("promoted piece incorrect: " + m.String())
	}
	m = NewPromotionCapture(colour.White, square.A7, square.A8, piece.ROOK, piece.QUEEN)
	if m.PromotedPiece() != piece.ROOK {
		t.Fatalf("promoted piece incorrect: " + m.String())
	}
	m = NewPromotionCapture(colour.White, square.A7, square.A8, piece.BISHOP, piece.QUEEN)
	if m.PromotedPiece() != piece.BISHOP {
		t.Fatalf("promoted piece incorrect: " + m.String())
	}
}
func TestPieceType(t *testing.T) {
	m := New(colour.White, square.A7, square.A8, piece.PAWN)
	if m.PieceType() != piece.PAWN {
		t.Fatalf("piece incorrect: " + m.String())
	}
	m = New(colour.White, square.A7, square.A8, piece.ROOK)
	if m.PieceType() != piece.ROOK {
		t.Fatalf("piece incorrect: " + m.String())
	}
	m = New(colour.White, square.A7, square.A8, piece.KNIGHT)
	if m.PieceType() != piece.KNIGHT {
		t.Fatalf("piece incorrect: " + m.String())
	}
	m = New(colour.White, square.A7, square.A8, piece.BISHOP)
	if m.PieceType() != piece.BISHOP {
		t.Fatalf("piece incorrect: " + m.String())
	}
	m = New(colour.White, square.A7, square.A8, piece.QUEEN)
	if m.PieceType() != piece.QUEEN {
		t.Fatalf("piece incorrect: " + m.String())
	}
	m = New(colour.White, square.A7, square.A8, piece.KING)
	if m.PieceType() != piece.KING {
		t.Fatalf("piece incorrect: " + m.String())
	}
}

func TestCastlingBits(t *testing.T) {
	m := New(colour.White, square.A5, square.B7, piece.KNIGHT)
	if m.CouldCastleBeforeMove(true) {
		t.Fatalf("castling not possible")
	}
	m.SetCastleBeforeMove(true)
	if !m.CouldCastleBeforeMove(true) {
		t.Fatalf("castling should be possible")
	}
	if m.CouldCastleBeforeMove(false) {
		t.Fatalf("castling not possible")
	}
	m.SetCastleBeforeMove(false)
	if !m.CouldCastleBeforeMove(false) {
		t.Fatalf("castling should be possible")
	}

	// opponent info

	if m.OpponentCouldCastleBeforeMove(true) {
		t.Fatalf("opponent castling not possible")
	}
	if m.OpponentCouldCastleBeforeMove(false) {
		t.Fatalf("opponent castling not possible")
	}

	m.SetOpponentCastleBeforeMove(true)
	if !m.OpponentCouldCastleBeforeMove(true) {
		t.Fatalf("opponent castling should be possible")
	}
	if m.OpponentCouldCastleBeforeMove(false) {
		t.Fatalf("opponent castling not possible")
	}
	m.SetOpponentCastleBeforeMove(false)
	if !m.OpponentCouldCastleBeforeMove(true) {
		t.Fatalf("opponent castling should be possible")
	}
	if !m.OpponentCouldCastleBeforeMove(false) {
		t.Fatalf("opponent castling should be possible")
	}
}

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
