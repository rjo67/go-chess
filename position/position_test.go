package position

import (
	"fmt"
	"strings"
	"testing"

	"github.com/rjo67/chess/bitset"
	"github.com/rjo67/chess/move"
	"github.com/rjo67/chess/piece"
	"github.com/rjo67/chess/piece/colour"
	"github.com/rjo67/chess/square"
)

func TestMakeMoveNonCapture(t *testing.T) {
	posn, err := ParseFen("2K2r2/4P3/8/1Q6/8/8/8/3k4 w - - 0 1")
	if err != nil {
		t.Fatalf("error parsing fen: %s", err)
	}
	queenBitset := bitset.NewFromSquares(square.B5)
	movedQueenBitset := bitset.NewFromSquares(square.F1)
	if posn.pieces[colour.White][piece.QUEEN].And(queenBitset).IsEmpty() {
		t.Fatalf("expected white queen at B5\n%s", posn.String())
	}
	if posn.allPieces[colour.White].And(queenBitset).IsEmpty() {
		t.Fatalf("allPieces wrong, no piece at B5\n%s", posn.String())
	}
	if posn.occupiedSquares.And(queenBitset).IsEmpty() {
		t.Fatalf("occupiedSquares wrong, no piece at B5\n%s", posn.String())
	}

	m := move.New(colour.White, square.B5, square.F1, piece.QUEEN)
	posn.MakeMove(m)

	// check after-effects of MakeMove
	if !posn.pieces[colour.White][piece.QUEEN].And(queenBitset).IsEmpty() {
		t.Fatalf("expected no white queen at B5\n%s", posn.String())
	}
	if !posn.allPieces[colour.White].And(queenBitset).IsEmpty() {
		t.Fatalf("allPieces wrong, a piece at B5\n%s", posn.String())
	}
	if !posn.occupiedSquares.And(queenBitset).IsEmpty() {
		t.Fatalf("occupiedSquares wrong, a piece at B5\n%s", posn.String())
	}
	if posn.pieces[colour.White][piece.QUEEN].And(movedQueenBitset).IsEmpty() {
		t.Fatalf("expected white queen at F1\n%s", posn.String())
	}
	if posn.allPieces[colour.White].And(movedQueenBitset).IsEmpty() {
		t.Fatalf("allPieces wrong, no piece at F1\n%s", posn.String())
	}
	if posn.occupiedSquares.And(movedQueenBitset).IsEmpty() {
		t.Fatalf("occupiedSquares wrong, no piece at F1\n%s", posn.String())
	}

	posn.UnmakeMove(m)

	// check after-effects of UnmakeMove
	if posn.pieces[colour.White][piece.QUEEN].And(queenBitset).IsEmpty() {
		t.Fatalf("expected white queen at B5\n%s", posn.String())
	}
	if posn.allPieces[colour.White].And(queenBitset).IsEmpty() {
		t.Fatalf("allPieces wrong, no piece at B5\n%s", posn.String())
	}
	if posn.occupiedSquares.And(queenBitset).IsEmpty() {
		t.Fatalf("occupiedSquares wrong, no piece at B5\n%s", posn.String())
	}
	if !posn.pieces[colour.White][piece.QUEEN].And(movedQueenBitset).IsEmpty() {
		t.Fatalf("expected no white queen at F1\n%s", posn.String())
	}
	if !posn.allPieces[colour.White].And(movedQueenBitset).IsEmpty() {
		t.Fatalf("allPieces wrong, piece at F1\n%s", posn.String())
	}
	if !posn.occupiedSquares.And(movedQueenBitset).IsEmpty() {
		t.Fatalf("occupiedSquares wrong, piece at F1\n%s", posn.String())
	}
}

func TestMakePromotionNonCapture(t *testing.T) {
	posn, err := ParseFen("2K2r2/4P3/8/8/8/8/8/3k4 w - - 0 1")
	if err != nil {
		t.Fatalf("error parsing fen: %s", err)
	}
	pawnBitset := bitset.NewFromSquares(square.E7)
	if posn.pieces[colour.White][piece.PAWN].And(pawnBitset).IsEmpty() {
		t.Fatalf("expected white pawn at E7\n%s", posn.String())
	}
	if posn.allPieces[colour.White].And(pawnBitset).IsEmpty() {
		t.Fatalf("allPieces wrong, no piece at E7\n%s", posn.String())
	}
	if posn.occupiedSquares.And(pawnBitset).IsEmpty() {
		t.Fatalf("occupiedSquares wrong, no piece at E7\n%s", posn.String())
	}

	m := move.NewPromotion(colour.White, square.E7, square.E8, piece.QUEEN)
	posn.MakeMove(m)

	promotedQueenBitset := bitset.NewFromSquares(square.E8)
	// check after-effects of MakeMove
	if posn.pieces[colour.White][piece.QUEEN].And(promotedQueenBitset).IsEmpty() {
		t.Fatalf("expected white queen at E8\n%s", posn.String())
	}
	if posn.allPieces[colour.White].And(promotedQueenBitset).IsEmpty() {
		t.Fatalf("allPieces wrong, no piece at E8\n%s", posn.String())
	}
	if posn.occupiedSquares.And(promotedQueenBitset).IsEmpty() {
		t.Fatalf("occupiedSquares wrong, no piece at E8\n%s", posn.String())
	}
	if !posn.pieces[colour.White][piece.PAWN].And(pawnBitset).IsEmpty() {
		t.Fatalf("expected no white pawn at E7\n%s", posn.String())
	}
	if !posn.allPieces[colour.White].And(pawnBitset).IsEmpty() {
		t.Fatalf("allPieces wrong, a piece at E7\n%s", posn.String())
	}
	if !posn.occupiedSquares.And(pawnBitset).IsEmpty() {
		t.Fatalf("occupiedSquares wrong, a piece at E7\n%s", posn.String())
	}

	posn.UnmakeMove(m)

	// check after-effects of UnmakeMove
	if posn.pieces[colour.White][piece.PAWN].And(pawnBitset).IsEmpty() {
		t.Fatalf("expected white queen at E7\n%s", posn.String())
	}
	if posn.allPieces[colour.White].And(pawnBitset).IsEmpty() {
		t.Fatalf("allPieces wrong, no piece at E7\n%s", posn.String())
	}
	if posn.occupiedSquares.And(pawnBitset).IsEmpty() {
		t.Fatalf("occupiedSquares wrong, no piece at E7\n%s", posn.String())
	}
	if !posn.pieces[colour.White][piece.QUEEN].And(promotedQueenBitset).IsEmpty() {
		t.Fatalf("expected no white queen at E8\n%s", posn.String())
	}
	if !posn.allPieces[colour.White].And(promotedQueenBitset).IsEmpty() {
		t.Fatalf("allPieces wrong, piece at E8\n%s", posn.String())
	}
	if !posn.occupiedSquares.And(promotedQueenBitset).IsEmpty() {
		t.Fatalf("occupiedSquares wrong, piece at E8\n%s", posn.String())
	}
}

func TestMakeMoveCapture(t *testing.T) {
	posn, err := ParseFen("2K2r2/4P3/8/1Q6/3R4/3p4/8/3k4 w - - 0 1")
	if err != nil {
		t.Fatalf("error parsing fen: %s", err)
	}
	rookBitset := bitset.NewFromSquares(square.D4)
	capturedPieceBitset := bitset.NewFromSquares(square.D3)

	// check pieces are where they should be
	if posn.pieces[colour.White][piece.ROOK].And(rookBitset).IsEmpty() {
		t.Fatalf("expected white rook at D4\n%s", posn.String())
	}
	if posn.allPieces[colour.White].And(rookBitset).IsEmpty() {
		t.Fatalf("allPieces wrong, no piece at D4\n%s", posn.String())
	}
	if posn.occupiedSquares.And(rookBitset).IsEmpty() {
		t.Fatalf("occupiedSquares wrong, no piece at D4\n%s", posn.String())
	}
	if posn.pieces[colour.Black][piece.PAWN].And(capturedPieceBitset).IsEmpty() {
		t.Fatalf("expected black pawn at D3\n%s", posn.String())
	}
	if posn.allPieces[colour.Black].And(capturedPieceBitset).IsEmpty() {
		t.Fatalf("allPieces wrong, no piece at D3\n%s", posn.String())
	}
	if posn.occupiedSquares.And(capturedPieceBitset).IsEmpty() {
		t.Fatalf("occupiedSquares wrong, no piece at D3\n%s", posn.String())
	}

	m := move.NewCapture(colour.White, square.D4, square.D3, piece.ROOK, piece.PAWN)
	posn.MakeMove(m)

	// check after-effects of MakeMove
	if !posn.pieces[colour.White][piece.ROOK].And(rookBitset).IsEmpty() {
		t.Fatalf("expected no white rook at D4\n%s", posn.String())
	}
	if !posn.allPieces[colour.White].And(rookBitset).IsEmpty() {
		t.Fatalf("allPieces wrong, a piece at D4\n%s", posn.String())
	}
	if !posn.occupiedSquares.And(rookBitset).IsEmpty() {
		t.Fatalf("occupiedSquares wrong, a piece at D4\n%s", posn.String())
	}
	// target square...
	if posn.pieces[colour.White][piece.ROOK].And(capturedPieceBitset).IsEmpty() {
		t.Fatalf("expected white rook at D3\n%s", posn.String())
	}
	if posn.allPieces[colour.White].And(capturedPieceBitset).IsEmpty() {
		t.Fatalf("allPieces wrong, no piece at D3\n%s", posn.String())
	}
	if !posn.pieces[colour.Black][piece.PAWN].And(capturedPieceBitset).IsEmpty() {
		t.Fatalf("expected no black pawn at D3\n%s", posn.String())
	}
	if !posn.allPieces[colour.Black].And(capturedPieceBitset).IsEmpty() {
		t.Fatalf("allPieces wrong, no piece at D3\n%s", posn.String())
	}
	if posn.occupiedSquares.And(capturedPieceBitset).IsEmpty() {
		t.Fatalf("occupiedSquares wrong, no piece at D3\n%s", posn.String())
	}

	posn.UnmakeMove(m)

	// check after-effects of UnmakeMove
	if posn.pieces[colour.White][piece.ROOK].And(rookBitset).IsEmpty() {
		t.Fatalf("expected white rook at D4\n%s", posn.String())
	}
	if posn.allPieces[colour.White].And(rookBitset).IsEmpty() {
		t.Fatalf("allPieces wrong, no piece at D4\n%s", posn.String())
	}
	if posn.occupiedSquares.And(rookBitset).IsEmpty() {
		t.Fatalf("occupiedSquares wrong, no piece at D4\n%s", posn.String())
	}
	if !posn.pieces[colour.White][piece.ROOK].And(capturedPieceBitset).IsEmpty() {
		t.Fatalf("expected no white rook at D3\n%s", posn.String())
	}
	if !posn.allPieces[colour.White].And(capturedPieceBitset).IsEmpty() {
		t.Fatalf("allPieces wrong, piece at D3\n%s", posn.String())
	}
	if posn.pieces[colour.Black][piece.PAWN].And(capturedPieceBitset).IsEmpty() {
		t.Fatalf("expected black pawn at D3\n%s", posn.String())
	}
	if posn.allPieces[colour.Black].And(capturedPieceBitset).IsEmpty() {
		t.Fatalf("allPieces wrong, no piece at D3\n%s", posn.String())
	}
	if posn.occupiedSquares.And(capturedPieceBitset).IsEmpty() {
		t.Fatalf("occupiedSquares wrong, no piece at D3\n%s", posn.String())
	}
}

func TestEnpassantCapture(t *testing.T) {
	posn, err := ParseFen("8/p7/8/1P6/K1k3p1/6P1/7P/8 b - - 0 10")
	if err != nil {
		t.Fatalf("error parsing fen: %s", err)
	}

	// first: black moves a7-a5 (setting enpassant square)
	m := move.New(colour.Black, square.A7, square.A5, piece.PAWN)
	if !m.HasEnpassantSquare() {
		t.Fatalf("move should have enpassant square")
	}
	if m.EnpassantSquare() != square.A6 {
		t.Fatalf("enpassant square should be A6 but was: %s", m.EnpassantSquare().String())
	}
	posn.MakeMove(m)

	blackPawnBitsetBeforeMove := bitset.NewFromSquares(square.A7)
	blackPawnBitsetAfterMove := bitset.NewFromSquares(square.A5)

	if posn.pieces[colour.Black][piece.PAWN].And(blackPawnBitsetAfterMove).IsEmpty() {
		t.Fatalf("expected pawn at A5\n%s", posn.String())
	}
	if posn.allPieces[colour.Black].And(blackPawnBitsetAfterMove).IsEmpty() {
		t.Fatalf("allPieces wrong, no piece at A5\n%s", posn.String())
	}
	if posn.occupiedSquares.And(blackPawnBitsetAfterMove).IsEmpty() {
		t.Fatalf("occupiedSquares wrong, no piece at A5\n%s", posn.String())
	}
	if posn.EnpassantSquare() == nil {
		t.Fatalf("posn should have enpassant square")
	}
	if *posn.EnpassantSquare() != square.A6 {
		t.Fatalf("enpassant square should be A6 but was: %s", (*posn.EnpassantSquare()).String())
	}

	posn.UnmakeMove(m)

	// check after-effects of UnmakeMove
	if posn.pieces[colour.Black][piece.PAWN].And(blackPawnBitsetBeforeMove).IsEmpty() {
		t.Fatalf("expected pawn at A7\n%s", posn.String())
	}
	if posn.allPieces[colour.Black].And(blackPawnBitsetBeforeMove).IsEmpty() {
		t.Fatalf("allPieces wrong, no piece at A7\n%s", posn.String())
	}
	if posn.occupiedSquares.And(blackPawnBitsetBeforeMove).IsEmpty() {
		t.Fatalf("occupiedSquares wrong, no piece at A7\n%s", posn.String())
	}
	if posn.EnpassantSquare() != nil {
		t.Fatalf("posn should not have enpassant square")
	}

	// redo the first (black) move
	posn.MakeMove(m)

	// second: white takes enpassant
	m = move.NewEpCapture(colour.White, square.B5, square.A6)

	posn.MakeMove(m)

	bothBlackPawnSquares := blackPawnBitsetBeforeMove.Or(blackPawnBitsetAfterMove)
	if !posn.pieces[colour.Black][piece.PAWN].And(bothBlackPawnSquares).IsEmpty() {
		t.Fatalf("expected no black pawns at A5/A7\n%s", posn.String())
	}
	if !posn.allPieces[colour.Black].And(bothBlackPawnSquares).IsEmpty() {
		t.Fatalf("allPieces wrong, should be no pieces at A5/A7\n%s", posn.String())
	}
	if !posn.occupiedSquares.And(bothBlackPawnSquares).IsEmpty() {
		t.Fatalf("occupiedSquares wrong, should be no pieces at A5/A7\n%s", posn.String())
	}
	if posn.EnpassantSquare() != nil {
		t.Fatalf("posn should not have enpassant square")
	}
	whitePawnSquareBeforeMove := bitset.NewFromSquares(square.B5)
	whitePawnSquareAfterMove := bitset.NewFromSquares(square.A6)
	if posn.pieces[colour.White][piece.PAWN].And(whitePawnSquareAfterMove).IsEmpty() {
		t.Fatalf("expected white pawn at A6\n%s", posn.String())
	}
	if posn.allPieces[colour.White].And(whitePawnSquareAfterMove).IsEmpty() {
		t.Fatalf("allPieces wrong, should be piece at A6\n%s", posn.String())
	}
	if posn.occupiedSquares.And(whitePawnSquareAfterMove).IsEmpty() {
		t.Fatalf("occupiedSquares wrong, should be piece at A6\n%s", posn.String())
	}

	posn.UnmakeMove(m)
	if !posn.pieces[colour.White][piece.PAWN].And(whitePawnSquareAfterMove).IsEmpty() {
		t.Fatalf("expected no white pawn at A6\n%s", posn.String())
	}
	if !posn.allPieces[colour.White].And(whitePawnSquareAfterMove).IsEmpty() {
		t.Fatalf("allPieces wrong, should be no piece at A6\n%s", posn.String())
	}
	if !posn.occupiedSquares.And(whitePawnSquareAfterMove).IsEmpty() {
		t.Fatalf("occupiedSquares wrong, should be no piece at A6\n%s", posn.String())
	}
	if posn.pieces[colour.White][piece.PAWN].And(whitePawnSquareBeforeMove).IsEmpty() {
		t.Fatalf("expected white pawn at B5\n%s", posn.String())
	}
	if posn.allPieces[colour.White].And(whitePawnSquareBeforeMove).IsEmpty() {
		t.Fatalf("allPieces wrong, should be piece at B5\n%s", posn.String())
	}
	if posn.occupiedSquares.And(whitePawnSquareBeforeMove).IsEmpty() {
		t.Fatalf("occupiedSquares wrong, should be piece at B5\n%s", posn.String())
	}
	if posn.EnpassantSquare() == nil {
		t.Fatalf("posn should have enpassant square")
	}
	if *posn.EnpassantSquare() != square.A6 {
		t.Fatalf("enpassant square should be A6 but was: %s", (*posn.EnpassantSquare()).String())
	}
}

func TestMakePromotionCapture(t *testing.T) {
	posn, err := ParseFen("2K2r2/4P3/8/8/8/8/8/3k4 w - - 0 1")
	if err != nil {
		t.Fatalf("error parsing fen: %s", err)
	}

	pawnBitset := bitset.NewFromSquares(square.E7)
	if posn.pieces[colour.White][piece.PAWN].And(pawnBitset).IsEmpty() {
		t.Fatalf("expected white pawn at E7\n%s", posn.String())
	}
	if posn.allPieces[colour.White].And(pawnBitset).IsEmpty() {
		t.Fatalf("allPieces wrong, no piece at E7\n%s", posn.String())
	}
	if posn.occupiedSquares.And(pawnBitset).IsEmpty() {
		t.Fatalf("occupiedSquares wrong, no piece at E7\n%s", posn.String())
	}

	m := move.NewPromotionCapture(colour.White, square.E7, square.F8, piece.QUEEN, piece.ROOK)
	posn.MakeMove(m)

	// check after-effects of MakeMove
	promotedQueenBitset := bitset.NewFromSquares(square.F8)
	capturedPieceBitset := promotedQueenBitset

	if posn.pieces[colour.White][piece.QUEEN].And(promotedQueenBitset).IsEmpty() {
		t.Fatalf("expected white queen at F8\n%s", posn.String())
	}
	if posn.allPieces[colour.White].And(promotedQueenBitset).IsEmpty() {
		t.Fatalf("allPieces wrong, no piece at F8\n%s", posn.String())
	}
	if !posn.pieces[colour.White][piece.PAWN].And(pawnBitset).IsEmpty() {
		t.Fatalf("expected no white pawn at E7\n%s", posn.String())
	}
	if !posn.allPieces[colour.White].And(pawnBitset).IsEmpty() {
		t.Fatalf("allPieces wrong, a piece at E7\n%s", posn.String())
	}
	if !posn.occupiedSquares.And(pawnBitset).IsEmpty() {
		t.Fatalf("occupiedSquares wrong, a piece at E7\n%s", posn.String())
	}
	if !posn.pieces[colour.Black][piece.ROOK].And(capturedPieceBitset).IsEmpty() {
		t.Fatalf("expected no black rook at F8\n%s", posn.String())
	}
	if !posn.allPieces[colour.Black].And(capturedPieceBitset).IsEmpty() {
		t.Fatalf("allPieces wrong, expected no piece at F8\n%s", posn.String())
	}
	if posn.occupiedSquares.And(promotedQueenBitset).IsEmpty() {
		t.Fatalf("occupiedSquares wrong, no piece at F8\n%s", posn.String())
	}

	posn.UnmakeMove(m)

	// check after-effects of UnmakeMove
	if posn.pieces[colour.White][piece.PAWN].And(pawnBitset).IsEmpty() {
		t.Fatalf("expected white pawn at E7\n%s", posn.String())
	}
	if posn.allPieces[colour.White].And(pawnBitset).IsEmpty() {
		t.Fatalf("allPieces wrong, no piece at E7\n%s", posn.String())
	}
	if posn.occupiedSquares.And(pawnBitset).IsEmpty() {
		t.Fatalf("occupiedSquares wrong, no piece at E7\n%s", posn.String())
	}
	if !posn.pieces[colour.White][piece.QUEEN].And(promotedQueenBitset).IsEmpty() {
		t.Fatalf("expected no white queen at F8\n%s", posn.String())
	}
	if !posn.allPieces[colour.White].And(promotedQueenBitset).IsEmpty() {
		t.Fatalf("allPieces wrong, piece at F8\n%s", posn.String())
	}
	if posn.pieces[colour.Black][piece.ROOK].And(capturedPieceBitset).IsEmpty() {
		t.Fatalf("expected black rook at F8\n%s", posn.String())
	}
	if posn.allPieces[colour.Black].And(capturedPieceBitset).IsEmpty() {
		t.Fatalf("allPieces wrong, expected a piece at F8\n%s", posn.String())
	}
	if posn.occupiedSquares.And(capturedPieceBitset).IsEmpty() {
		t.Fatalf("occupiedSquares wrong, expected piece at F8\n%s", posn.String())
	}
}

func TestAttacksSquare(t *testing.T) {
	data := []struct {
		fen           string
		targetSquares []square.Square
		col           colour.Colour
		pieceType     piece.Piece
		expected      []bool
	}{
		{"1Q3k1B/ppBp4/4Q3/R5NR/2nB1P2/3N1N2/P2PPPP1/6K1 w KQkq - 0 0",
			[]square.Square{square.E5, square.H2, square.C6}, colour.White, piece.ROOK, []bool{true, true, false}},
		{"1Q3k1B/ppBp4/4Q3/R5NR/2nB1P2/3N1N2/P2PPPP1/6K1 w KQkq - 0 0",
			[]square.Square{square.D6, square.C7, square.A7}, colour.White, piece.BISHOP, []bool{true, false, true}},
		{"1Q3k1B/ppBp4/4Q3/R5NR/2nB1P2/3N1N2/P2PPPP1/6K1 w KQkq - 0 0",
			[]square.Square{square.E3, square.B5, square.B3}, colour.White, piece.QUEEN, []bool{true, false, false}},
		{"1Q3k1B/ppBp4/4Q3/R5NR/2nB1P2/3N1N2/P2PPPP1/6K1 w KQkq - 0 0",
			[]square.Square{square.E7, square.E8, square.F5}, colour.Black, piece.KING, []bool{true, true, false}},
		{"1Q3k1B/ppBp4/4Q3/R5NR/2nB1P2/3N1N2/P2PPPP1/6K1 w KQkq - 0 0",
			[]square.Square{square.F2, square.F1, square.G3}, colour.White, piece.KING, []bool{true, true, false}},
		{"1Q3k1B/ppBp4/4Q3/R5NR/2nB1P2/3N1N2/P2PPPP1/6K1 w KQkq - 0 0",
			[]square.Square{square.B2, square.F7, square.F5}, colour.White, piece.KNIGHT, []bool{true, true, false}},
		{"1Q3k1B/ppBp4/4Q3/R5NR/2nB1P2/3N1N2/P2PPPP1/6K1 w KQkq - 0 0",
			[]square.Square{square.E5, square.F5, square.G5}, colour.White, piece.PAWN, []bool{true, false, true}},
		{"1Q3k1B/ppBp4/4Q3/R5NR/2nB1P2/3N1N2/P2PPPP1/6K1 w KQkq - 0 0",
			[]square.Square{square.C6, square.D6, square.E6}, colour.Black, piece.PAWN, []bool{true, false, true}},
	}

	for testNbr, d := range data {
		posn, err := ParseFen(d.fen)
		if err != nil {
			t.Fatalf("error parsing fen '%s': %s", d.fen, err)
		}
		for i, sq := range d.targetSquares {
			if d.expected[i] != posn.PieceAttacksSquare(d.col, d.pieceType, sq) {
				t.Fatalf("wrong result for test#%d.%d", testNbr, i)
			}
		}
	}
}

func TestAttacks(t *testing.T) {
	data := []struct {
		fen                  string
		startSquare          square.Square
		expectedWhiteSetBits []int
		expectedBlackSetBits []int
	}{
		{"1Q3k1B/ppBp4/4Q3/R5NR/2nB1P2/3N1N2/P2PPPP1/6K1 w KQkq - 0 0", square.E5, []int{19, 21, 27, 29, 40, 44, 54, 57}, []int{30}},
		{"1Q3k1B/ppBp4/4Q3/R5NR/2nB1P2/3N1N2/P2PPPP1/6K1 w KQkq - 0 0", square.H2, []int{2, 19, 33}, []int{}},
		{"1Q3k1B/ppBp4/4Q3/R5NR/2nB1P2/3N1N2/P2PPPP1/6K1 w KQkq - 0 0", square.E3, []int{11, 13, 29, 44}, []int{30}},
		{"1Q3k1B/ppBp4/4Q3/R5NR/2nB1P2/3N1N2/P2PPPP1/6K1 w KQkq - 0 0", square.B6, []int{54, 44, 29}, []int{30, 56}},
		{"1Q3k1B/ppBp4/4Q3/R5NR/2nB1P2/3N1N2/P2PPPP1/6K1 w KQkq - 0 0", square.A6, []int{40, 44}, []int{55}},
		{"1Q3k1B/ppBp4/4Q3/R5NR/2nB1P2/3N1N2/P2PPPP1/6K1 w KQkq - 0 0", square.F7, []int{34, 44}, []int{59}},
	}

	for testNbr, d := range data {
		posn, err := ParseFen(d.fen)
		if err != nil {
			t.Fatalf("error parsing fen '%s': %s", d.fen, err)
		}

		bs := posn.Attacks(d.startSquare, colour.White)
		if len(d.expectedWhiteSetBits) != bs.Cardinality() {
			t.Fatalf("test %d (white): got cardinality %d (expected %d) bitset:\n%s", testNbr, bs.Cardinality(), len(d.expectedWhiteSetBits), bs.String())
		} else {
			checkBits(testNbr, bs, d.expectedWhiteSetBits, t)
		}
		bs = posn.Attacks(d.startSquare, colour.Black)
		if len(d.expectedBlackSetBits) != bs.Cardinality() {
			t.Fatalf("test %d (black): got cardinality %d (expected %d) bitset:\n%s", testNbr, bs.Cardinality(), len(d.expectedBlackSetBits), bs.String())
		}
		checkBits(testNbr, bs, d.expectedBlackSetBits, t)

		anyBits := append(d.expectedWhiteSetBits, d.expectedBlackSetBits...)
		bs = posn.Attacks(d.startSquare, colour.AnyColour)
		if len(anyBits) != bs.Cardinality() {
			t.Fatalf("test %d (any colour): got cardinality %d (expected %d) bitset:\n%s", testNbr, bs.Cardinality(), len(anyBits), bs.String())
		}
		checkBits(testNbr, bs, anyBits, t)
	}
}

type moveData struct {
	fen              string
	expectedNbrMoves []int
}

func TestInitialPosition(t *testing.T) {
	doTest(moveData{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 0", []int{20, 400, 8902, 197281, 4865609}}, t)
}
func TestPosn2(t *testing.T) {
	doTest(moveData{"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []int{48, 2039, 97862, 4085603}}, t)
}
func TestPosn3(t *testing.T) {
	doTest(moveData{"8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 0", []int{14, 191, 2812, 43238, 674624}}, t)
}
func TestPosn5(t *testing.T) {
	doTest(moveData{"rnbqkb1r/pp1p1ppp/2p5/4P3/2B5/8/PPP1NnPP/RNBQK2R w KQkq - 0 6", []int{42, 1352, 53392}}, t)
}
func TestPosn6(t *testing.T) {
	doTest(moveData{"r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10", []int{46, 2079, 89890, 3894594, 164075551}}, t)
}
func TestNumpty2(t *testing.T) {
	doTest(moveData{"8/p7/8/1P6/K1k3p1/6P1/7P/8 w - - 0 10", []int{5, 39, 237, 2002, 14062, 120995, 966152}}, t)
}
func TestNumpty3(t *testing.T) {
	doTest(moveData{"r3k2r/p6p/8/B7/1pp1p3/3b4/P6P/R3K2R w KQkq - 0 10", []int{17, 341, 6666, 150072, 3186478}}, t)
}
func TestNumpty4(t *testing.T) {
	doTest(moveData{"8/5p2/8/2k3P1/p3K3/8/1P6/8 b - - 0 10", []int{9, 85, 795, 7658, 72120, 703851}}, t)
}
func TestNumpty5(t *testing.T) {
	doTest(moveData{"r3k2r/pb3p2/5npp/n2p4/1p1PPB2/6P1/P2N1PBP/R3K2R b KQkq - 0 10", []int{29, 953, 27990, 909807}}, t)
}
func TestIllegalEpMove1(t *testing.T) {
	doTest(moveData{"8/8/8/8/k1p4R/8/3P4/3K4 w - - 0 1", []int{18, 92, 1670, 10138, 185429, 1134888}}, t)
	doTest(moveData{"8/8/8/8/k1p4R/8/3P4/3K4 b - - 0 1", []int{5, 89, 555, 10094, 61765, 1124950}}, t)
}
func TestIllegalEpMove2(t *testing.T) {
	doTest(moveData{"8/8/4k3/8/2p5/8/B2P2K1/8 w - - 0 1", []int{13, 102, 1266, 10276, 135655, 1015133}}, t)
	doTest(moveData{"8/8/4k3/8/2p5/8/B2P2K1/8 b - - 0 1", []int{8, 104, 872, 11047, 84630, 1139270}}, t)
}

func TestEpResultsInCheck(t *testing.T) {
	// enpassant move C4xD3 is illegal, because of the rook check
	doTest(moveData{"8/8/8/8/1kpP3R/8/B5K1/8 b - d3 0 1", []int{6, 136, 732, 16861, 99272}}, t)
}
func TestEpCaptureChecksOpponent(t *testing.T) {
	doTest(moveData{"8/5k2/8/2Pp4/2B5/1K6/8/8 w - d6 0 1", []int{15, 126, 1928, 13931, 206379, 1440467}}, t)
}
func TestShortCastlingChecksOpponent(t *testing.T) {
	doTest(moveData{"5k2/8/8/8/8/8/8/4K2R w K - 0 1", []int{15, 66, 1198, 6399, 120330, 661072}}, t)
	doTest(moveData{"4k2r/8/8/8/8/8/8/4K2R b k - 0 1", []int{15, 171, 2601, 38779, 621743}}, t)
}
func TestLongCastlingChecksOpponent(t *testing.T) {
	doTest(moveData{"3k4/8/8/8/8/8/8/R3K3 w Q - 0 1", []int{16, 71, 1286, 7418, 141077, 803711}}, t)
	doTest(moveData{"r3k3/8/8/8/8/8/8/3K4 b q - 0 1", []int{16, 71, 1286, 7418, 141077, 803711}}, t)
}
func TestCastlingIncludingLosingOrRookCapture(t *testing.T) {
	doTest(moveData{"r3k2r/1b4bq/8/8/8/8/7B/R3K2R w KQkq - 0 1", []int{26, 1141, 27826, 1274206}}, t)
	doTest(moveData{"r3k2r/1b4bq/8/8/8/8/7B/R3K2R b KQkq - 0 1", []int{47, 1011, 47973, 1105187}}, t)
}
func TestCastlingPrevented(t *testing.T) {
	doTest(moveData{"r3k2r/8/5Q2/8/8/3q4/8/R3K2R w KQkq - 0 1", []int{44, 1494, 50509, 1720476}}, t)
	doTest(moveData{"r3k2r/8/3Q4/8/8/5q2/8/R3K2R b KQkq - 0 1", []int{44, 1494, 50509, 1720476}}, t)
}
func TestPromoteOutOfCheck(t *testing.T) {
	doTest(moveData{"2K2r2/4P3/8/8/8/8/8/3k4 w - - 0 1", []int{11, 133, 1442, 19174, 266199, 3821001}}, t)
	doTest(moveData{"3K4/8/8/8/8/8/4p3/2k2R2 b - - 0 1", []int{11, 133, 1442, 19174, 266199, 3821001}}, t)
}
func TestDiscoveredCheck(t *testing.T) {
	doTest(moveData{"8/8/8/2k3PR/8/1p2K3/2P2B2/2Q5 w - - 0 10", []int{31, 223, 7685, 54476}}, t)
}
func TestDiscoveredCheck2(t *testing.T) {
	doTest(moveData{"5K2/8/1Q6/2N5/8/1p2k3/8/8 w - - 0 1", []int{29, 165, 5160, 31961, 1004658}}, t)
	doTest(moveData{"8/8/1P2K3/8/2n5/1q6/8/5k2 b - - 0 1", []int{29, 165, 5160, 31961, 1004658}}, t)
}
func TestSelfStalemate(t *testing.T) {
	doTest(moveData{"8/k1P5/8/1K6/8/8/8/8 w - - 0 1", []int{10, 25, 268, 926, 10857, 43261, 567584}}, t)
	doTest(moveData{"8/8/8/8/1k6/8/K1p5/8 b - - 0 1", []int{10, 25, 268, 926, 10857, 43261, 567584}}, t)
}
func TestSelfStalemate2(t *testing.T) {
	doTest(moveData{"K1k5/8/P7/8/8/8/8/8 w - - 0 1", []int{2, 6, 13, 63, 382, 2217, 15453}}, t)
	doTest(moveData{"8/8/8/8/8/p7/8/k1K5 b - - 0 1", []int{2, 6, 13, 63, 382, 2217, 15453}}, t)
}
func TestPromotionRocechess(t *testing.T) {
	//www.rocechess.ch/perft.html
	doTest(moveData{"n1n5/PPPk4/8/8/8/8/4Kppp/5N1N b - - 0 1", []int{24, 496, 9483, 182838, 3605103 /* , 71179139 */}}, t)
}
func TestPromotionToGiveCheck(t *testing.T) {
	doTest(moveData{"4k3/1P6/8/8/8/8/K7/8 w - - 0 1", []int{9, 40, 472, 2661, 38983, 217342}}, t)
	doTest(moveData{"8/k7/8/8/8/8/1p6/4K3 b - - 0 1", []int{9, 40, 472, 2661, 38983, 217342}}, t)
}
func TestUnderPromoteToGiveCheck(t *testing.T) {
	doTest(moveData{"8/P1k5/K7/8/8/8/8/8 w - - 0 1", []int{6, 27, 273, 1329, 18135, 92683}}, t)
	doTest(moveData{"8/8/8/8/8/k7/p1K5/8 b - - 0 1", []int{6, 27, 273, 1329, 18135, 92683}}, t)
}
func TestDoubleCheck(t *testing.T) {
	doTest(moveData{"8/5k2/8/5N2/5Q2/2K5/8/8 w - - 0 1", []int{37, 183, 6559, 23527, 811573}}, t)
	doTest(moveData{"8/8/2k5/5q2/5n2/8/5K2/8 b - - 0 1", []int{37, 183, 6559, 23527, 811573}}, t)
}

func doTest(data moveData, t *testing.T) {
	posn, err := ParseFen(data.fen)
	if err != nil {
		t.Fatalf("could not parse fen: %s, err: %s", data.fen, err.Error())
	}
	for depth, expectedNbrMoves := range data.expectedNbrMoves {
		if depth < 99 { // can use this to limit which tests are carried out
			if expectedNbrMoves != -1 {
				var moveMapString strings.Builder // for logging i/c of error
				moveMap := perft(posn, depth+1)
				var moveCount int
				var moveNbr int
				for key, entry := range moveMap {
					moveNbr++
					moveCount += entry
					moveMapString.WriteString(fmt.Sprintf("%3d: %8s\t%3d\n", moveNbr, key, entry))
				}
				if moveCount != expectedNbrMoves {
					t.Fatalf("depth: %d, expected %d moves but got %d, fen: %s, moveMap:\n%s", depth+1, expectedNbrMoves, moveCount, data.fen, moveMapString.String())
				}
			}
		}
	}
}

// entry point for Perft.
// a "move-map" will be returned. The keys of the map are the possible first half-moves in this position.
// The values of the map are the number of LEAF moves from each starting move.
func perft(posn Position, depth int) map[string]int {
	moveMap := make(map[string]int, 0)
	// fill move map with starting moves
	for _, startMove := range posn.FindMoves(posn.activeColour) {
		posn.MakeMove(startMove)
		moveMap[startMove.String()] = len(p2(startMove, posn, depth)) // just store the number of moves, to allow GC of the moves
		posn.UnmakeMove(startMove)
	}
	return moveMap
}

// processes one half-move level.
func p2(origMove move.Move, posn Position, depth int) []move.Move {
	if depth == 1 {
		// at a leaf node, return an array of one element, in order to count the leaf-moves
		return []move.Move{origMove}
	}
	movesAtNextDepth := make([]move.Move, 0, 300)
	for _, m := range posn.FindMoves(posn.activeColour) {
		posn.MakeMove(m)
		movesAtNextDepth = append(movesAtNextDepth, p2(m, posn, depth-1)...)
		posn.UnmakeMove(m)
	}
	return movesAtNextDepth
}

func checkBits(testNbr int, bs bitset.BitSet, expectedSetBits []int, t *testing.T) {
	errors := make([]int, 0, 5)
	for _, bit := range expectedSetBits {
		if !bs.IsSet(uint(bit)) {
			errors = append(errors, bit)
		}
	}
	if len(errors) != 0 {
		t.Fatalf("test %d: found %d errors (%v) for bitset:\n%s", testNbr, len(errors), errors, bs.String())
	}
}
