package position

import (
	"strings"
	"testing"

	"github.com/rjo67/chess/bitset"
	"github.com/rjo67/chess/piece"
	"github.com/rjo67/chess/piece/colour"
	"github.com/rjo67/chess/square"
)

func TestBadNbrFields(t *testing.T) {
	_, err := ParseFen("8/8/8/8/8/8/8/8 w KQkq - 0")
	if err == nil {
		t.Error("expected error")
	}
}

func TestNoKingDefined(t *testing.T) {
	_, err := ParseFen("7k/8/8/8/8/8/8/8 w KQkq - 0 0")
	checkErrorMessage(err, "king not defined", t)
	_, err = ParseFen("8/8/8/8/8/8/8/7K w KQkq - 0 0")
	checkErrorMessage(err, "king not defined", t)
}

func TestField1Error(t *testing.T) {
	_, err := ParseFen("7k/8/8/8/8/8/8 w KQkq - 0 0")
	checkErrorMessage(err, "number of fields", t)
	_, err = ParseFen("8/5rpQP/8/8/8/8/8/7K w KQkq - 0 0")
	checkErrorMessage(err, "subfield 2 too long at position 5", t)
	_, err = ParseFen("8/5r21/8/8/8/8/8/7K w KQkq - 0 0")
	checkErrorMessage(err, "subfield 2 too long at position 4", t)
	_, err = ParseFen("7k/8/8/9/8/8/8/7K w KQkq - 0 0")
	checkErrorMessage(err, "unrecognised: '9' at position 1 of subfield 4", t)
	_, err = ParseFen("7k/8/8/7/8/8/8/7K w KQkq - 0 0")
	checkErrorMessage(err, "subfield 4 too short", t)
}

func TestCorrectField1(t *testing.T) {
	fen := "7k/p3r3/3q2q1/pppn3b/8/7K/8/8 w KQkq - 0 1"
	posn, err := ParseFen(fen)
	if err != nil {
		t.Errorf("%s", err.Error())
	}
	data := []struct {
		bits  []int
		piece piece.Piece
	}{
		{[]int{56, 40, 39, 38}, piece.PAWN},
		{[]int{52}, piece.ROOK},
		{[]int{37}, piece.KNIGHT},
		{[]int{33}, piece.BISHOP},
		{[]int{42, 45}, piece.QUEEN},
		{[]int{57}, piece.KING},
	}

	for _, i := range data {
		bs := posn.BitSetFor(colour.Black, i.piece)
		checkBitsForPiece(i.piece, i.bits, bs, t)
	}
}

func TestField2Error(t *testing.T) {
	_, err := ParseFen("7k/8/8/8/8/8/8/7K - KQkq - 0 0")
	checkErrorMessage(err, "unrecognised colour", t)
	_, err = ParseFen("7k/8/8/8/8/8/8/7K W KQkq - 0 0")
	checkErrorMessage(err, "unrecognised colour", t)
	_, err = ParseFen("7k/8/8/8/8/8/8/7K wb KQkq - 0 0")
	checkErrorMessage(err, "unrecognised colour", t)
}

func TestField3(t *testing.T) {
	_, err := ParseFen("7k/8/8/8/8/8/8/7K w K3kq - 0 0")
	checkErrorMessage(err, "castling availability syntax error", t)
	_, err = ParseFen("7k/8/8/8/8/8/8/7K w KKkq - 0 0")
	checkErrorMessage(err, "castling availability syntax error", t)
	_, err = ParseFen("7k/8/8/8/8/8/8/7K w Kqq - 0 0")
	checkErrorMessage(err, "castling availability syntax error", t)
	posn, err := ParseFen("7k/8/8/8/8/8/8/7K w Kq - 0 0")
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	} else {
		if !posn.CastlingAvailabilityKingsSide(colour.White) {
			t.Errorf("expected kings-side castling for white")
		}
		if posn.CastlingAvailabilityQueensSide(colour.White) {
			t.Errorf("did not expect queens-side castling for white")
		}
		if posn.CastlingAvailabilityKingsSide(colour.Black) {
			t.Errorf("did not expect kings-side castling for black")
		}
		if !posn.CastlingAvailabilityQueensSide(colour.Black) {
			t.Errorf("expected queens-side castling for black")
		}
	}
	posn, err = ParseFen("7k/8/8/8/8/8/8/7K w - - 0 0")
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	} else {
		if posn.CastlingAvailabilityKingsSide(colour.White) {
			t.Errorf("did not expect kings-side castling for white")
		}
		if posn.CastlingAvailabilityQueensSide(colour.White) {
			t.Errorf("did not expect queens-side castling for white")
		}
		if posn.CastlingAvailabilityKingsSide(colour.Black) {
			t.Errorf("did not expect kings-side castling for black")
		}
		if posn.CastlingAvailabilityQueensSide(colour.Black) {
			t.Errorf("did not queens-side castling for black")
		}
	}
}

func TestField4(t *testing.T) {
	_, err := ParseFen("7k/8/8/8/8/8/8/7K w KQkq a0 0 0")
	checkErrorMessage(err, "unrecognised square", t)
	posn, err := ParseFen("7k/8/8/8/8/8/8/7K w KQkq a6 0 0")
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	} else {
		if posn.EnpassantSquare() == nil {
			t.Error("enpassant square in posn is nil")
		} else if *posn.EnpassantSquare() != square.A6 {
			t.Errorf("expected A6 but got: %d", *posn.EnpassantSquare())
		}
	}
	posn, err = ParseFen("8/5k2/8/2Pp4/2B5/1K6/8/8 b - d6 0 1")
	checkErrorMessage(err, "invalid e.p. square 'D6' for active colour: B", t)
}

func TestField5(t *testing.T) {
	_, err := ParseFen("7k/8/8/8/8/8/8/7K w KQkq - a 0")
	checkErrorMessage(err, "could not parse halfmove clock", t)
	_, err = ParseFen("7k/8/8/8/8/8/8/7K w KQkq - -2 0")
	checkErrorMessage(err, "invalid value for halfmove clock", t)
	posn, err := ParseFen("7k/8/8/8/8/8/8/7K w KQkq - 2 0")
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if posn.HalfmoveClock() != 2 {
		t.Errorf("expected 2 but got: %d", posn.HalfmoveClock())
	}
}

func TestField6(t *testing.T) {
	_, err := ParseFen("7k/8/8/8/8/8/8/7K w KQkq - 0 a")
	checkErrorMessage(err, "could not parse fullmove number", t)
	_, err = ParseFen("7k/8/8/8/8/8/8/7K w KQkq - 0 -2")
	checkErrorMessage(err, "invalid value for fullmove number", t)
	posn, err := ParseFen("7k/8/8/8/8/8/8/7K w KQkq - 0 4")
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if posn.FullmoveNbr() != 4 {
		t.Errorf("expected 4 but got: %d", posn.FullmoveNbr())
	}
}

func checkErrorMessage(err error, expectedMessage string, t *testing.T) {
	if err == nil {
		t.Errorf("expected error with text '%s'", expectedMessage)
	} else if !strings.Contains(err.Error(), expectedMessage) {
		t.Errorf("expected '%s' but got message: '%s'", expectedMessage, err.Error())
	}
}

func checkBitsForPiece(piece piece.Piece, bits []int, bs bitset.BitSet, t *testing.T) {
	for _, i := range bits {
		if !bs.IsSet(uint(i)) {
			t.Errorf("checking piece %d: bit %d should be set for bitset:\n%s", piece, i, bs.String())
		}
	}
}
