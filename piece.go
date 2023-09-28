package main

type Piece struct {
	x     int
	y     int
	kind  string
	color int
	worth int
}

// starting duck is put to -1,-1 so that it doesn't show
func generate_default_duck() Piece {
	return Piece{x: -1, y: -1, kind: "duck", color: 0, worth: 0}
}

func piece_to_unicode(p Piece) string {
	if p.kind == "empty" {
		return "·"

	} else if p.kind == "duck" {
		return "Θ"
	} else if p.color == 1 {
		if p.kind == "pawn" {
			return "\u2659"
		} else if p.kind == "knight" {
			return "\u2658"
		} else if p.kind == "bishop" {
			return "\u2657"
		} else if p.kind == "rook" {
			return "\u2656"
		} else if p.kind == "queen" {
			return "\u2655"
		} else if p.kind == "king" {
			return "\u2654"
		}
	} else if p.color == -1 {
		if p.kind == "pawn" {
			return "\u265F"
		} else if p.kind == "knight" {
			return "\u265E"
		} else if p.kind == "bishop" {
			return "\u265D"
		} else if p.kind == "rook" {
			return "\u265C"
		} else if p.kind == "queen" {
			return "\u265B"
		} else if p.kind == "king" {
			return "\u265A"
		}
	}

	panic("Could not display unicode piece")
}

func generate_empty_piece(x, y int) Piece {
	return Piece{x: x, y: y, kind: "empty", color: 0, worth: 0}
}

func translate_to_PGN(p *Piece) string {
	if p.kind == "pawn" {
		return ""
	} else if p.kind == "duck" {
		return "Θ"
	} else if p.kind == "queen" {
		return "Q"
	} else if p.kind == "knight" {
		return "N"
	} else if p.kind == "bishop" {
		return "B"
	} else if p.kind == "rook" {
		return "R"
	} else {

		panic("could not translate to PGN")
	}

}
