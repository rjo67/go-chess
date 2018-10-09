package position

import (
	"strings"

	"github.com/rjo67/chess/bitset"
	"github.com/rjo67/chess/piece"
	"github.com/rjo67/chess/piece/colour"
	"github.com/rjo67/chess/ray"
	"github.com/rjo67/chess/square"
)

// Position represents a chess position
type Position struct {
	pieces                         []map[piece.Piece]bitset.BitSet // array of map of piece bitsets, array-dim = colour
	allPieces                      []bitset.BitSet                 // all pieces of a particular colour
	occupiedSquares                bitset.BitSet                   // all occupied squares
	activeColour                   colour.Colour                   // whose move
	castlingAvailabilityKingsSide  []bool                          // whether white/black can castle kingsside
	castlingAvailabilityQueensSide []bool                          // whether white/black can castle queensside
	enpassantSquare                *square.Square
	halfmoveClock                  int
	fullmoveNbr                    int
}

// NewPosition creates a new position
// The bitset arrays are in the order as given by the piece constants
func NewPosition(whitePieces, blackPieces map[piece.Piece]bitset.BitSet) Position {
	var p Position
	p.pieces = make([]map[piece.Piece]bitset.BitSet, 2)
	p.allPieces = make([]bitset.BitSet, 2)
	p.castlingAvailabilityKingsSide = make([]bool, 2)
	p.castlingAvailabilityQueensSide = make([]bool, 2)

	p.pieces[colour.White] = whitePieces
	p.pieces[colour.Black] = blackPieces

	// calculate further bitmaps:
	for _, col := range colour.AllColours {
		p.allPieces[col] = bitset.BitSet{}
		for _, pieceType := range piece.AllPieces {
			p.allPieces[col] = p.allPieces[col].Or(p.pieces[col][pieceType])
		}
	}
	p.occupiedSquares = p.allPieces[colour.White].Or(p.allPieces[colour.Black])

	return p
}

// StartPosition creates a new start position
func StartPosition() Position {
	pieces := make([]map[piece.Piece]bitset.BitSet, 2)

	for _, col := range colour.AllColours {
		pieces[col] = make(map[piece.Piece]bitset.BitSet)
		for _, pieceType := range piece.AllPieces {
			pieces[col][pieceType] = piece.StartPosn[col][pieceType]
		}
	}

	return NewPosition(pieces[colour.White], pieces[colour.Black])
}

// AllPieces returns a bitset with all the occupied squares for the given colour
func (p Position) AllPieces(col colour.Colour) bitset.BitSet {
	return p.allPieces[col]
}

// OccupiedSquares returns the occupied-squares bitset
func (p Position) OccupiedSquares() bitset.BitSet {
	return p.occupiedSquares
}

// Pieces returns the bitset for the given piece type and colour
func (p Position) Pieces(col colour.Colour, pieceType piece.Piece) bitset.BitSet {
	return p.pieces[col][pieceType]
}

// EnpassantSquare returns the current enpassant square (or nil)
func (p Position) EnpassantSquare() *square.Square {
	return p.enpassantSquare
}

// HalfmoveClock returns the current halfmove clock
func (p Position) HalfmoveClock() int {
	return p.halfmoveClock
}

// FullmoveNbr returns the current fullmove nbr
func (p Position) FullmoveNbr() int {
	return p.fullmoveNbr
}

// CastlingAvailability returns the castling availabilty
func (p Position) CastlingAvailability(col colour.Colour, kingsside bool) bool {
	if kingsside {
		return p.castlingAvailabilityKingsSide[col]
	}
	return p.castlingAvailabilityQueensSide[col]
}

// BitSetFor returns the bitset for the given piece and colour
func (p Position) BitSetFor(col colour.Colour, piece piece.Piece) bitset.BitSet {
	return p.pieces[col][piece]
}

// Attacks returns a bitset of pieces which attack the given square.
// If requiredColour is specified, only the pieces of this colour will be returned.
// (If requiredColour is AnyColour, all attacks are returned)
func (p Position) Attacks(sq square.Square, requiredColour colour.Colour) bitset.BitSet {
	bs := bitset.BitSet{}

	var bishops, rooks, queens, knights, pawns, kings bitset.BitSet
	diagonals := p._find2(int(sq), ray.AllBishopDirections)
	rankfiles := p._find2(int(sq), ray.AllRookDirections)

	if requiredColour == colour.AnyColour {
		bishops = diagonals.And(p.Pieces(colour.White, piece.BISHOP).Or(p.Pieces(colour.Black, piece.BISHOP)))
		rooks = rankfiles.And(p.Pieces(colour.White, piece.ROOK).Or(p.Pieces(colour.Black, piece.ROOK)))
		knights = ray.KnightAttackBitSets[sq].And(p.Pieces(colour.White, piece.KNIGHT).Or(p.Pieces(colour.Black, piece.KNIGHT)))
		pawns = (ray.PawnAttackBitSets[colour.White][sq].And(p.Pieces(colour.White, piece.PAWN))).Or(ray.PawnAttackBitSets[colour.Black][sq].And(p.Pieces(colour.Black, piece.PAWN)))
		kings = ray.KingAttackBitSets[sq].And(p.Pieces(colour.White, piece.KING).Or(p.Pieces(colour.Black, piece.KING)))

		allQueens := p.Pieces(colour.White, piece.QUEEN).Or(p.Pieces(colour.Black, piece.QUEEN))
		queens = diagonals.And(allQueens).Or(rankfiles.And(allQueens))
	} else {
		// filter on required colour
		bishops = diagonals.And(p.Pieces(requiredColour, piece.BISHOP))
		rooks = rankfiles.And(p.Pieces(requiredColour, piece.ROOK))
		knights = ray.KnightAttackBitSets[sq].And(p.Pieces(requiredColour, piece.KNIGHT))
		queens = diagonals.And(p.Pieces(requiredColour, piece.QUEEN)).Or(rankfiles.And(p.Pieces(requiredColour, piece.QUEEN)))
		pawns = ray.PawnAttackBitSets[requiredColour][sq].And(p.Pieces(requiredColour, piece.PAWN))
		kings = ray.KingAttackBitSets[sq].And(p.Pieces(requiredColour, piece.KING))
	}

	return bs.Or(bishops).Or(rooks).Or(queens).Or(knights).Or(pawns).Or(kings)
}

// String delivers a string representation of the current position
func (p Position) String() string {
	var squares [64]string

	// populate squares with the bitset contents
	for _, col := range colour.AllColours {
		for _, pieceType := range piece.AllPieces {
			bits := p.pieces[col][pieceType].SetBits()
			for _, sq := range bits {
				squares[sq-1] = pieceType.ToString(col)
			}
		}
	}

	var sb strings.Builder
	sb.WriteString("+--------+\n")
	for i := 63; i >= 0; i-- {
		if i%8 == 7 {
			sb.WriteString("|")
		}
		if squares[i] == "" {
			sb.WriteString(".")
		} else {
			sb.WriteString(squares[i])
		}
		if i%8 == 0 {
			sb.WriteString("|\n")
		}
	}
	sb.WriteString("+--------+\n")

	return sb.String()
}

// 'package private', only called by Builder
func (p *Position) setEnpassantSquare(sq *square.Square) {
	p.enpassantSquare = sq
}

// 'package private', only called by Builder
func (p *Position) setHalfmoveClock(clock int) {
	p.halfmoveClock = clock
}

// 'package private', only called by Builder
func (p *Position) setFullmoveNbr(move int) {
	p.fullmoveNbr = move
}

// 'package private', only called by Builder
func (p *Position) setCastlingAvailability(col colour.Colour, kingsside bool, canCastle bool) {
	if kingsside {
		p.castlingAvailabilityKingsSide[col] = canCastle
	}
	p.castlingAvailabilityQueensSide[col] = canCastle
}
