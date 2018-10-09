package position

import (
	"testing"

	"github.com/rjo67/chess/bitset"
	"github.com/rjo67/chess/piece/colour"
	"github.com/rjo67/chess/square"
)

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
			t.Errorf("error parsing fen '%s': %s", d.fen, err)
		} else {
			bs := posn.Attacks(d.startSquare, colour.White)
			if len(d.expectedWhiteSetBits) != bs.Cardinality() {
				t.Errorf("test %d (white): got cardinality %d (expected %d) bitset:\n%s", testNbr, bs.Cardinality(), len(d.expectedWhiteSetBits), bs.String())
			} else {
				checkBits(testNbr, bs, d.expectedWhiteSetBits, t)
			}
			bs = posn.Attacks(d.startSquare, colour.Black)
			if len(d.expectedBlackSetBits) != bs.Cardinality() {
				t.Errorf("test %d (black): got cardinality %d (expected %d) bitset:\n%s", testNbr, bs.Cardinality(), len(d.expectedBlackSetBits), bs.String())
			} else {
				checkBits(testNbr, bs, d.expectedBlackSetBits, t)
			}
			anyBits := append(d.expectedWhiteSetBits, d.expectedBlackSetBits...)
			bs = posn.Attacks(d.startSquare, colour.AnyColour)
			if len(anyBits) != bs.Cardinality() {
				t.Errorf("test %d (any colour): got cardinality %d (expected %d) bitset:\n%s", testNbr, bs.Cardinality(), len(anyBits), bs.String())
			} else {
				checkBits(testNbr, bs, anyBits, t)
			}
		}
	}
}

func TestFindMoves(t *testing.T) {
	data := []struct {
		name             string
		fen              string
		col              colour.Colour
		expectedNbrMoves []int
	}{
		{"initialPosition", "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 0", colour.White, []int{20, 400, 8902, 197281, 4865609}},
		{"posn2", "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", colour.White, []int{48, 2039, 97862, 4085603}},
		{"posn3", "8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 0", colour.White, []int{14, 191, 2812, 43238, 674624}},
		{"posn5", "rnbqkb1r/pp1p1ppp/2p5/4P3/2B5/8/PPP1NnPP/RNBQK2R w KQkq - 0 6", colour.White, []int{42, 1352, 53392}},
		{"posn6", "r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10", colour.White, []int{46, 2079, 89890, 3894594 /*, 164075551*/}},
		{"numpty2", "8/p7/8/1P6/K1k3p1/6P1/7P/8 w - - 0 10", colour.White, []int{5, 39, 237, 2002, 14062, 120995, 966152}},
		{"numpty3", "r3k2r/p6p/8/B7/1pp1p3/3b4/P6P/R3K2R w KQkq - 0 10", colour.White, []int{17, 341, 6666, 150072, 3186478}},
		{"numpty4", "8/5p2/8/2k3P1/p3K3/8/1P6/8 b - - 0 10", colour.Black, []int{9, 85, 795, 7658, 72120, 703851}},
		{"numpty5", "r3k2r/pb3p2/5npp/n2p4/1p1PPB2/6P1/P2N1PBP/R3K2R b KQkq - 0 10", colour.Black, []int{29, 953, 27990, 909807}},
		/*
			{"illegalEpMove1W", "8/8/8/8/k1p4R/8/3P4/3K4 w - - 0 1", colour.White, []int{-1, -1, -1, -1, -1, 1134888}},
			{"illegalEpMove1B", "3k4/3p4/8/K1P4r/8/8/8/8 b - - 0 1", colour.Black, []int{-1, -1, -1, -1, -1, 1134888}},
			{"illegalEpMove2W", "8/8/4k3/8/2p5/8/B2P2K1/8 w - - 0 1", colour.White, []int{-1, -1, -1, -1, -1, 1015133}},
			{"illegalEpMove2B", "8/b2p2k1/8/2P5/8/4K3/8/8 b - - 0 1", colour.Black, []int{-1, -1, -1, -1, -1, 1015133}},
			{"epCaptureChecksOpponentW", "8/5k2/8/2Pp4/2B5/1K6/8/8 w - d6 0 1", colour.White, []int{-1, -1, -1, -1, -1, 1440467}},
			{"epCaptureChecksOpponentB", "8/8/1k6/2b5/2pP4/8/5K2/8 b - d3 0 1", colour.Black, []int{-1, -1, -1, -1, -1, 1440467}},
			{"shortCastlingChecksOpponentW", "5k2/8/8/8/8/8/8/4K2R w K - 0 1", colour.White, []int{-1, -1, -1, -1, -1, 661072}},
			{"shortCastlingChecksOpponentB", "4k2r/8/8/8/8/8/8/5K2 b k - 0 1", colour.Black, []int{-1, -1, -1, -1, -1, 661072}},
			{"longCastlingChecksOpponentW", "3k4/8/8/8/8/8/8/R3K3 w Q - 0 1", colour.White, []int{-1, -1, -1, -1, -1, 803711}},
			{"longCastlingChecksOpponentB", "r3k3/8/8/8/8/8/8/3K4 b q - 0 1", colour.Black, []int{-1, -1, -1, -1, -1, 803711}},
			{"castlingIncludingLosingOrRookCaptureW", "r3k2r/1b4bq/8/8/8/8/7B/R3K2R w KQkq - 0 1", colour.White, []int{-1, -1, -1, 1274206}},
			{"castlingIncludingLosingOrRookCaptureB", "r3k2r/7b/8/8/8/8/1B4BQ/R3K2R b KQkq - 0 1", colour.Black, []int{-1, -1, -1, 1274206}},
			{"castlingPreventedW", "r3k2r/8/5Q2/8/8/3q4/8/R3K2R w KQkq - 0 1", colour.White, []int{-1, -1, -1, 1720476}},
			{"castlingPreventedB", "r3k2r/8/3Q4/8/8/5q2/8/R3K2R b KQkq - 0 1", colour.Black, []int{-1, -1, -1, 1720476}},
			{"promoteOutOfCheckW", "2K2r2/4P3/8/8/8/8/8/3k4 w - - 0 1", colour.White, []int{-1, -1, -1, -1, -1, 3821001}},
			{"promoteOutOfCheckB", "3K4/8/8/8/8/8/4p3/2k2R2 b - - 0 1", colour.Black, []int{-1, -1, -1, -1, -1, 3821001}},
		*/
		{"discoveredCheck", "8/8/8/2k3PR/8/1p2K3/2P2B2/2Q5 w - - 0 10", colour.White, []int{31}},
		/*
			{"discoveredCheck2W", "5K2/8/1Q6/2N5/8/1p2k3/8/8 w - - 0 1", colour.White, []int{-1, -1, -1, -1, 1004658}},
			{"discoveredCheck2B", "8/8/1P2K3/8/2n5/1q6/8/5k2 b - - 0 1", colour.Black, []int{-1, -1, -1, -1, 1004658}},
			{"selfStalemateCheckmate", "8/k1P5/8/1K6/8/8/8/8 w - - 0 1", colour.White, []int{-1, -1, -1, -1, -1, -1, 567584}},
			{"selfStalemateCheckmate2", "8/8/8/8/1k6/8/K1p5/8 b - - 0 1", colour.Black, []int{-1, -1, -1, -1, -1, -1, 567584}},
			{"selfStalemateW", "K1k5/8/P7/8/8/8/8/8 w - - 0 1", colour.White, []int{-1, -1, -1, -1, -1, 2217}},
			{"selfStalemateB", "8/8/8/8/8/p7/8/k1K5 b - - 0 1", colour.Black, []int{-1, -1, -1, -1, -1, 2217}},
			{"promotion www.rocechess.ch/perft.html", "n1n5/PPPk4/8/8/8/8/4Kppp/5N1N b - - 0 1", colour.Black, []int{24, 496, 9483, 182838, 3605103// , 71179139 }},
			{"promoteToGiveCheckW", "4k3/1P6/8/8/8/8/K7/8 w - - 0 1", colour.White, []int{-1, -1, -1, -1, -1, 217342}},
			{"promoteToGiveCheckB", "8/k7/8/8/8/8/1p6/4K3 b - - 0 1", colour.Black, []int{-1, -1, -1, -1, -1, 217342}},
			{"underpromoteToGiveCheckW", "8/P1k5/K7/8/8/8/8/8 w - - 0 1", colour.White, []int{-1, -1, -1, -1, -1, 92683}},
			{"underpromoteToGiveCheckB", "8/8/8/8/8/k7/p1K5/8 b - - 0 1", colour.Black, []int{-1, -1, -1, -1, -1, 92683}},
			{"doubleCheckW", "8/5k2/8/5N2/5Q2/2K5/8/8 w - - 0 1", colour.White, []int{-1, -1, -1, 23527}},
			{"doubleCheckB", "8/8/2k5/5q2/5n2/8/5K2/8 b - - 0 1", colour.Black, []int{-1, -1, -1, 23527}},
		*/
	}
	for i, d := range data {
		testNbr := i + 1
		position, err := ParseFen(d.fen)
		if err != nil {
			t.Errorf("test %d (%s): could not parse fen: %s, err: %s", testNbr, d.name, d.fen, err.Error())
		} else {
			moves := position.FindMoves(d.col)
			if len(moves) != d.expectedNbrMoves[0] {
				t.Errorf("test %d (%s): expected %d moves but got %d: %v", testNbr, d.name, d.expectedNbrMoves[0], len(moves), moves)
			}
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
