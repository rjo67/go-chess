package position

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rjo67/chess/bitset"
	"github.com/rjo67/chess/piece"
	"github.com/rjo67/chess/piece/colour"
	"github.com/rjo67/chess/square"
)

const (
	badNbrFields               string = "Wrong number of fields"
	noWhiteKing                string = "White king not defined"
	noBlackKing                string = "Black king not defined"
	castlingAvailabilitySyntax string = "castling availability syntax error"
)

// ParseError encapsulates errors found whilst parsing
type ParseError struct {
	msg   string // description of error
	field int    // field position in input where error was found
}

func (e ParseError) Error() string {
	return fmt.Sprintf("%s in field %d", e.msg, e.field)
}

// ParseFen creates a position from a FEN string
// https://en.wikipedia.org/wiki/Forsyth%E2%80%93Edwards_Notation
func ParseFen(fen string) (Position, error) {
	fields := strings.Split(fen, " ")
	if len(fields) != 6 {
		return Position{}, ParseError{badNbrFields, 0}
	}

	// process the first field
	pieceMap, err := processField1(fields[0])
	if err != nil {
		return Position{}, err
	}
	// sanity check
	if pieceMap[piece.WhiteKingString].Cardinality() != 1 {
		return Position{}, ParseError{noWhiteKing, 1}
	}
	if pieceMap[piece.BlackKingString].Cardinality() != 1 {
		return Position{}, ParseError{noBlackKing, 1}
	}

	builder := NewBuilder()
	for _, col := range colour.AllColours {
		for _, pieceStr := range piece.PieceMapping[col] {
			//			fmt.Printf("adding colour %d, piece %s, bitset:\n%s", col, pieceStr, pieceMap[pieceStr].ToString())
			builder.AddPiece(col, piece.FromString(col, pieceStr), pieceMap[pieceStr])
		}
	}

	// second field: activeColour
	activeColour, err := processField2(fields[1])
	if err != nil {
		return Position{}, err
	}
	builder.ActiveColour(activeColour)

	// third field: castling rights
	castlingAvailabilityKingsSide, castlingAvailabilityQueensSide, err := processField3(fields[2])
	if err != nil {
		return Position{}, err
	}
	builder.CastlingAvailability(colour.White, true, castlingAvailabilityKingsSide[colour.White])
	builder.CastlingAvailability(colour.White, false, castlingAvailabilityQueensSide[colour.White])
	builder.CastlingAvailability(colour.Black, true, castlingAvailabilityKingsSide[colour.Black])
	builder.CastlingAvailability(colour.Black, false, castlingAvailabilityQueensSide[colour.Black])

	// fourth field: enpassant square
	enpassantSquare, err := processField4(fields[3], activeColour)
	if err != nil {
		return Position{}, err
	}
	builder.EnpassantSquare(enpassantSquare)

	// fifth field: halfmove clock
	halfmoveClock, err := processField5(fields[4])
	if err != nil {
		return Position{}, err
	}
	builder.HalfmoveClock(halfmoveClock)

	// sixth field: fullmove nbr
	fullmoveNbr, err := processField6(fields[5])
	if err != nil {
		return Position{}, err
	}
	builder.FullmoveNbr(fullmoveNbr)

	posn := builder.Build()
	return posn, err
}

// first field -- piece information in 8 subfields separated by '/'
func processField1(field1 string) (map[string]*bitset.BitSet, error) {
	subFields := strings.Split(field1, "/")
	if len(subFields) != 8 {
		return nil, ParseError{badNbrFields, 1}
	}

	// set up a map of bitsets corresponding to the piece identifiers
	pieceMap := make(map[string]*bitset.BitSet)
	for _, col := range colour.AllColours {
		for _, piece := range piece.PieceMapping[col] {
			pieceMap[piece] = &bitset.BitSet{}
		}
	}
	rankOffset := 72 // "points to" the left-hand file of the chessboard. Starts at rank 8 and works down
	for fieldNbr := 0; fieldNbr < 8; fieldNbr++ {
		rankOffset -= 8
		fileOffset := 1 // 1..8 corresponding to the values in the fen subfields
		for i := 0; i < len(subFields[fieldNbr]); i++ {
			str := string(subFields[fieldNbr][i])
			switch str {
			case piece.WhitePawnString, piece.WhiteRookString, piece.WhiteKnightString, piece.WhiteBishopString, piece.WhiteQueenString, piece.WhiteKingString,
				piece.BlackPawnString, piece.BlackRookString, piece.BlackKnightString, piece.BlackBishopString, piece.BlackQueenString, piece.BlackKingString:
				if fileOffset > 8 {
					return nil, ParseError{fmt.Sprintf("subfield %d too long at position %d", fieldNbr+1, i+1), 1}
				}
				pieceMap[str].Set(uint(rankOffset - (fileOffset - 1)))
				fileOffset++
			case "1", "2", "3", "4", "5", "6", "7", "8":
				val, _ := strconv.Atoi(str)
				if fileOffset+val > 9 { // since fileOffset starts at 1
					return nil, ParseError{fmt.Sprintf("subfield %d too long at position %d", fieldNbr+1, i+1), 1}
				}
				fileOffset += val
			default:
				return nil, ParseError{fmt.Sprintf("unrecognised: '%s' at position %d of subfield %d", str, i+1, fieldNbr+1), 1}
			}
		}
		// at the end of the subfields, 'fileOffset' must be 9, otherwise the definition was too short
		if fileOffset != 9 {
			return nil, ParseError{fmt.Sprintf("subfield %d too short", fieldNbr+1), 1}
		}

	}
	return pieceMap, nil
}

// second field: activeColour
func processField2(field string) (colour.Colour, error) {
	var activeColour colour.Colour
	switch field {
	case "w":
		activeColour = colour.White
	case "b":
		activeColour = colour.Black
	default:
		return 0, ParseError{fmt.Sprintf("unrecognised colour: '%s'", field), 2}
	}
	return activeColour, nil
}

// third field: castling rights
func processField3(field string) ([]bool, []bool, error) {
	kingsSide := make([]bool, 2)
	queensSide := make([]bool, 2)
	if field == "-" {
		return kingsSide, queensSide, nil
	}
	var kingsSideWhite, kingsSideBlack, queensSideWhite, queensSideBlack bool
	for i := 0; i < len(field); i++ {
		str := string(field[i])
		switch str {
		case "K":
			if kingsSideWhite {
				return nil, nil, ParseError{fmt.Sprintf("%s %s", castlingAvailabilitySyntax, "(multiple 'K')"), 3}
			}
			kingsSideWhite = true
			kingsSide[colour.White] = true
		case "k":
			if kingsSideBlack {
				return nil, nil, ParseError{fmt.Sprintf("%s %s", castlingAvailabilitySyntax, "(multiple 'k')"), 3}
			}
			kingsSideBlack = true
			kingsSide[colour.Black] = true
		case "Q":
			if queensSideWhite {
				return nil, nil, ParseError{fmt.Sprintf("%s %s", castlingAvailabilitySyntax, "(multiple 'Q')"), 3}
			}
			queensSideWhite = true
			queensSide[colour.White] = true
		case "q":
			if queensSideBlack {
				return nil, nil, ParseError{fmt.Sprintf("%s %s", castlingAvailabilitySyntax, "(multiple 'q')"), 3}
			}
			queensSideBlack = true
			queensSide[colour.Black] = true
		default:
			return nil, nil, ParseError{castlingAvailabilitySyntax, 3}
		}
	}
	return kingsSide, queensSide, nil
}

// fourth field: enpassant square
func processField4(field string, activeColour colour.Colour) (*square.Square, error) {
	if field == "-" {
		return nil, nil
	}
	sq, err := square.FromString(field)
	if err != nil {
		return nil, err
	}
	if (activeColour == colour.White && sq.Rank() != 6) || (activeColour == colour.Black && sq.Rank() != 3) {
		return nil, ParseError{fmt.Sprintf("invalid e.p. square '%s' for active colour: %s", sq.String(), activeColour.String()), 3}
	}
	return &sq, nil
}

// fifth field: halfmove clock
func processField5(field string) (int, error) {
	i, err := strconv.Atoi(field)
	if err != nil {
		return 0, fmt.Errorf("could not parse halfmove clock: '%s'", field)
	}
	if i < 0 {
		return 0, fmt.Errorf("invalid value for halfmove clock: '%d'", i)
	}
	return i, nil
}

// sixth field: fullmove nbr
func processField6(field string) (int, error) {
	i, err := strconv.Atoi(field)
	if err != nil {
		return 0, fmt.Errorf("could not parse fullmove number: '%s'", field)
	}
	if i < 0 {
		return 0, fmt.Errorf("invalid value for fullmove number: '%d'", i)
	}
	return i, nil
}
