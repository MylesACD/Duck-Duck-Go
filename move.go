package main

import "strconv"

type Move struct {
	piece   Piece
	sx      int
	sy      int
	cap     bool
	ex      int
	ey      int
	extra   string
	passant bool
	castle  bool
}

var files = map[int]string{
	0: "a",
	1: "b",
	2: "c",
	3: "d",
	4: "e",
	5: "f",
	6: "g",
	7: "h",
}

func (m Move) String() string {
	str := translate_to_PGN(&m.piece)

	// castle case
	if m.castle {
		return m.extra

		// pawn case
	} else if str == "" {
		if !m.cap {
			return files[m.ex] + strconv.Itoa(8-m.ey) + m.extra
		}

		return files[m.sx] + "x" + files[m.ex] + strconv.Itoa(8-m.ey) + m.extra

	} else if !m.cap {
		return str + files[m.sx] + strconv.Itoa(8-m.sy) + files[m.ex] + strconv.Itoa(8-m.ey) + m.extra
	}

	return str + files[m.sx] + strconv.Itoa(8-m.sy) + "x" + files[m.ex] + strconv.Itoa(8-m.ey) + m.extra
}
