package position

import (
	"fmt"

	"github.com/rjo67/chess/bitset"
	"github.com/rjo67/chess/move"
	"github.com/rjo67/chess/piece"
	"github.com/rjo67/chess/piece/colour"
	"github.com/rjo67/chess/ray"
	"github.com/rjo67/chess/square"
)

// FindMoves returns all legal moves in the current position for the given colour.FindMoves
func (p Position) FindMoves(col colour.Colour) []move.Move {
	potentiallyIllegalMoves := p.FindPotentiallyIllegalMoves(col)

	moves := make([]move.Move, 0, len(potentiallyIllegalMoves))
	otherColour := col.Other()
	opponentsKing := square.Square(p.Pieces(otherColour, piece.KING).SetBits()[0])
	myKing := square.Square(p.Pieces(col, piece.KING).SetBits()[0])

	for _, move := range potentiallyIllegalMoves {
		valid := true
		if move.IsCastles() {
			// castle out of check is not allowed
			if p.AnyPieceAttacksSquare(otherColour, myKing) {
				valid = false
			} else {
				// check if castling allowed (relevant squares are free, and not moving into check)
				var bitsetToCheck bitset.BitSet
				var squaresToCheck []square.Square
				if move.IsKingsSideCastles() {
					bitsetToCheck = kingssideCastlingsBitMaps[col]
					squaresToCheck = kingssideCastlingsSquares[col]
				} else {
					bitsetToCheck = queenssideCastlingsBitMaps[col]
					squaresToCheck = queenssideCastlingsSquares[col]
				}
				if p.OccupiedSquares().And(bitsetToCheck).IsEmpty() {
					// anything attacking the relevant squares?
					for _, sq := range squaresToCheck {
						if p.AnyPieceAttacksSquare(otherColour, sq) {
							valid = false
							break
						}
					}
				} else {
					valid = false
				}
			}
		} else {
			if move.IsKingsMove() {
				// is king adjacent to other king
				if opponentsKing.IsAdjacentTo(move.To()) {
					valid = false
				} else {
					// is king moving into check..?
					p.MakeMove(&move)
					if p.AnyPieceAttacksSquare(otherColour, move.To()) {
						valid = false
					}
					p.UnmakeMove(move)
				}
			} else {
				// for other moves: is king now in check
				p.MakeMove(&move)
				if p.AnyPieceAttacksSquare(otherColour, myKing) {
					valid = false
				}
				p.UnmakeMove(move)
			}
		}
		if valid {
			moves = append(moves, move)
		}
	}

	return moves
}

// FindPotentiallyIllegalMoves returns all moves for the given colour in the given position.
// The returned list of moves can contain illegal moves e.g. because of moving into check (see FindMoves)
func (p Position) FindPotentiallyIllegalMoves(col colour.Colour) []move.Move {
	moves := make([]move.Move, 0, 60)

	moves = append(moves, p.findPawnMoves(col)...)
	moves = append(moves, p.findRookMoves(col)...)
	moves = append(moves, p.findKnightMoves(col)...)
	moves = append(moves, p.findBishopMoves(col)...)
	moves = append(moves, p.findQueenMoves(col)...)
	moves = append(moves, p.findKingMoves(col)...)

	return moves
}

func (p Position) findPawnMoves(col colour.Colour) []move.Move {
	moves := make([]move.Move, 0, 50)

	var shift int
	var rankMask, promotionMask, everythingExceptPromotion bitset.BitSet
	if col == colour.White {
		shift = 8
		rankMask = bitset.Rank2
		promotionMask = bitset.Rank8
		everythingExceptPromotion = bitset.NotRank8
	} else {
		shift = -8
		rankMask = bitset.Rank7
		promotionMask = bitset.Rank1
		everythingExceptPromotion = bitset.NotRank1
	}
	// move all pawns up one square, and again for two squares if starting on rank 2
	pawns := p.Pieces(col, piece.PAWN)
	emptySquares := p.OccupiedSquares().Not()
	oneSquare := pawns.Shift(shift).And(emptySquares)
	twoSquares := pawns.And(rankMask).Shift(shift).And(emptySquares).Shift(shift).And(emptySquares)

	promotedPawns := oneSquare.And(promotionMask)
	for _, bit := range promotedPawns.SetBits() {
		for _, promotedPiece := range piece.PromotedPawnPieceCandidates {
			moves = append(moves, move.NewPromotion(col, square.Square(bit-shift), square.Square(bit), promotedPiece))
		}
	}
	for _, bit := range oneSquare.And(everythingExceptPromotion).SetBits() {
		moves = append(moves, move.New(col, square.Square(bit-shift), square.Square(bit), piece.PAWN))
	}
	for _, bit := range twoSquares.SetBits() {
		moves = append(moves, move.New(col, square.Square(bit-(2*shift)), square.Square(bit), piece.PAWN))
	}

	// captures...

	if col == colour.White {
		shift = 9
		rankMask = bitset.NotFile1
	} else {
		shift = -9
		rankMask = bitset.NotFile8
	}
	otherColour := col.Other()
	captureLeft := pawns.And(rankMask).Shift(shift).And(p.AllPieces(otherColour))
	promotedPawns = captureLeft.And(promotionMask)
	for _, bit := range promotedPawns.SetBits() {
		for _, promotedPiece := range piece.PromotedPawnPieceCandidates {
			moves = append(moves, move.NewPromotionCapture(col, square.Square(bit-shift), square.Square(bit), promotedPiece, p.PieceAt(uint(bit), otherColour)))
		}
	}
	for _, bit := range captureLeft.And(everythingExceptPromotion).SetBits() {
		moves = append(moves, move.NewCapture(col, square.Square(bit-shift), square.Square(bit), piece.PAWN, p.PieceAt(uint(bit), otherColour)))
	}

	if col == colour.White {
		shift = 7
		rankMask = bitset.NotFile8
	} else {
		shift = -7
		rankMask = bitset.NotFile1
	}
	captureRight := pawns.And(rankMask).Shift(shift).And(p.AllPieces(otherColour))
	promotedPawns = captureRight.And(promotionMask)
	for _, bit := range promotedPawns.SetBits() {
		for _, promotedPiece := range piece.PromotedPawnPieceCandidates {
			moves = append(moves, move.NewPromotionCapture(col, square.Square(bit-shift), square.Square(bit), promotedPiece, p.PieceAt(uint(bit), otherColour)))
		}
	}
	for _, bit := range captureRight.And(everythingExceptPromotion).SetBits() {
		moves = append(moves, move.NewCapture(col, square.Square(bit-shift), square.Square(bit), piece.PAWN, p.PieceAt(uint(bit), otherColour)))
	}

	if p.EnpassantSquare() != nil {
		epAttacks := ray.AttacksOnEnpassantSquares[col][p.EnpassantSquare().File()]
		for _, bit := range epAttacks.And(pawns).SetBits() {
			moves = append(moves, move.NewEpCapture(col, square.Square(bit), *p.EnpassantSquare()))
		}
	}

	return moves
}

func (p Position) findRookMoves(col colour.Colour) []move.Move {
	return p._findForPiece(col, piece.ROOK, ray.AllRookDirections)
}

func (p Position) findKnightMoves(col colour.Colour) []move.Move {
	moves := make([]move.Move, 0, 8)
	otherColour := col.Other()
	for _, startSq := range p.Pieces(col, piece.KNIGHT).SetBits() {
		bs := ray.KnightAttackBitSets[startSq] //TODO and with opponents pieces
		for _, bit := range bs.SetBits() {
			if p.AllPieces(col).IsSet(uint(bit)) {
				// do nothing - square is occupied with a piece of my own colour
			} else if p.AllPieces(otherColour).IsSet(uint(bit)) {
				// capture
				moves = append(moves, move.NewCapture(col, square.Square(startSq), square.Square(bit), piece.KNIGHT, p.PieceAt(uint(bit), otherColour)))
			} else {
				// empty square
				moves = append(moves, move.New(col, square.Square(startSq), square.Square(bit), piece.KNIGHT))
			}
		}
	}
	return moves
}

func (p Position) findBishopMoves(col colour.Colour) []move.Move {
	return p._findForPiece(col, piece.BISHOP, ray.AllBishopDirections)
}

func (p Position) findQueenMoves(col colour.Colour) []move.Move {
	return p._findForPiece(col, piece.QUEEN, ray.AllDirections)

}

// these squares must be empty
var kingssideCastlingsBitMaps = [2]bitset.BitSet{bitset.New(0x06), bitset.New(0x0600000000000000)}

// these squares cannot be attacked
var kingssideCastlingsSquares = [2][]square.Square{{square.F1, square.G1}, {square.F8, square.G8}}

// kingssideCastlingsRookMove is the XOR for rook's move, kingsside castling
var kingssideCastlingsRookMove = [2]bitset.BitSet{bitset.New(0x05), bitset.New(0x0500000000000000)}

// these squares must be empty
var queenssideCastlingsBitMaps = [2]bitset.BitSet{bitset.New(0x70), bitset.New(0x7000000000000000)}

// these squares cannot be attacked
var queenssideCastlingsSquares = [2][]square.Square{{square.C1, square.D1}, {square.C8, square.D8}}

// queenssideCastlingsRookMove is the XOR for rook's move, queenssside castling
var queenssideCastlingsRookMove = [2]bitset.BitSet{bitset.New(0x90), bitset.New(0x9000000000000000)}

// in which directions can the 'castling fields' possibly be attacked?
var castlingAttackDirections = [][]ray.Direction{{ray.NORTHWEST, ray.NORTH, ray.NORTHEAST}, {ray.SOUTHWEST, ray.SOUTH, ray.SOUTHEAST}}

// returns castling-moves if theoretically possible, does not check if they are legal
func (p Position) findKingMoves(col colour.Colour) []move.Move {
	moves := make([]move.Move, 0, 10)
	otherColour := col.Other()
	for _, startSq := range p.Pieces(col, piece.KING).SetBits() {
		bs := ray.KingAttackBitSets[startSq].AndNot(p.AllPieces(col)) // remove my own pieces
		for _, bit := range bs.SetBits() {
			if p.AllPieces(otherColour).IsSet(uint(bit)) {
				// capture
				moves = append(moves, move.NewCapture(col, square.Square(startSq), square.Square(bit), piece.KING, p.PieceAt(uint(bit), otherColour)))
			} else {
				// empty square
				moves = append(moves, move.New(col, square.Square(startSq), square.Square(bit), piece.KING))
			}
		}
	}

	// castling moves are added without checking for legality
	if p.CastlingAvailabilityKingsSide(col) {
		moves = append(moves, move.CastleKingsSide(col))
	}
	if p.CastlingAvailabilityQueensSide(col) {
		moves = append(moves, move.CastleQueensSide(col))
	}
	return moves
}

// AnyPieceAttacksSquare returns true if any piece of the given colour attacks the target square
func (p Position) AnyPieceAttacksSquare(col colour.Colour, targetSq square.Square) bool {
	for _, pieceType := range piece.AllPieces {
		if p.PieceAttacksSquare(col, pieceType, targetSq) {
			return true
		}
	}
	return false
}

// PieceAttacksSquare returns true if a piece of the given type and colour attacks the target square
func (p Position) PieceAttacksSquare(col colour.Colour, pieceType piece.Piece, targetSq square.Square) bool {
	var directions []ray.Direction
	switch pieceType {
	case piece.KNIGHT:
		possibleMoves := ray.KnightAttackBitSets[targetSq].And(p.Pieces(col, pieceType))
		return !possibleMoves.IsEmpty()
	case piece.KING:
		possibleMoves := ray.KingAttackBitSets[targetSq].And(p.Pieces(col, pieceType))
		return !possibleMoves.IsEmpty()
	case piece.PAWN:
		//TODO take enpassant square into account
		if col == colour.White {
			bs := bitset.NewFromSquares(targetSq)
			possibleMoves := (bs.And(bitset.NotFile1).Shift(-7)).Or(bs.And(bitset.NotFile8).Shift(-9))
			return !possibleMoves.And(p.Pieces(col, pieceType)).IsEmpty()
		}
		bs := bitset.NewFromSquares(targetSq)
		possibleMoves := (bs.And(bitset.NotFile1).Shift(9)).Or(bs.And(bitset.NotFile8).Shift(7))
		return !possibleMoves.And(p.Pieces(col, pieceType)).IsEmpty()
	case piece.ROOK:
		directions = ray.AllRookDirections
	case piece.BISHOP:
		directions = ray.AllBishopDirections
	case piece.QUEEN:
		directions = ray.AllDirections
	default:
		panic("bad piece")
	}
	var found bool
	target := int(targetSq)
	for _, direction := range directions {
		possibleMoves, _ := move.Search2(target, direction, p.OccupiedSquares())
		// concentrate on required piece type
		possibleMoves = possibleMoves.And(p.Pieces(col, pieceType))

		if !possibleMoves.IsEmpty() {
			found = true
			break
		}
	}
	return found
}

// Finds all moves for a given pieceType and Colour, in all given directions,
// using the piece squares from the current position.
// Invalid moves are not discarded here, i.e. a returned move may be illegal because of moving into check
func (p Position) _findForPiece(col colour.Colour, pieceType piece.Piece, directions []ray.Direction) []move.Move {
	moves := make([]move.Move, 0, 20)
	otherColour := col.Other()
	for _, startSq := range p.Pieces(col, pieceType).SetBits() {
		// must iterate over all directions separately (instead of calling _find2) because of the 'blocker' logic
		for _, direction := range directions {
			possibleMoves, blockingSquare := move.Search2(startSq, direction, p.OccupiedSquares())
			// only need to check the 'blocker' square (if present) for the colour of the piece it contains (if any)
			for _, bit := range possibleMoves.SetBits() {
				if bit == blockingSquare {
					if p.AllPieces(col).IsSet(uint(bit)) {
						// do nothing - square is occupied with a piece of my own colour
					} else if p.AllPieces(otherColour).IsSet(uint(bit)) {
						// capture
						moves = append(moves, move.NewCapture(col, square.Square(startSq), square.Square(bit), pieceType, p.PieceAt(uint(bit), otherColour)))
					} else {
						panic(fmt.Sprintf("blocking square %d set but no piece present", blockingSquare))
					}
				} else {
					moves = append(moves, move.New(col, square.Square(startSq), square.Square(bit), pieceType))
				}
			}
		}
	}
	return moves
}

// Returns a bitset of all squares which can be attacked in the given directions from the given square.
// Another way of putting it: returns a bitset of all squares which attack the given square in the given directions.
//
// This is irrespective of which pieces are on the given squares, if any.
// TODO rename to AttacksOnSquare?
func (p Position) _find2(startSq int, directions []ray.Direction) bitset.BitSet {
	possibleMoves := bitset.New(0)
	for _, direction := range directions {
		bs, _ := move.Search2(startSq, direction, p.OccupiedSquares())
		possibleMoves = possibleMoves.Or(bs)
	}
	return possibleMoves
}
