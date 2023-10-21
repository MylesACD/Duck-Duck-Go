package main

const BLACK = -1
const WHITE = 1
const OUTCOME_MULTIPLIER = 1000

// white will always be positive and black negative
// 0,0 is top left corner
type GameState struct {
	white_mat    int
	black_mat    int
	is_end_state bool
	turn_num     int
	curr_player  int
	duck         Piece
	result       int

	board [8][8]Piece
	// if a game has 100 entries of previously reversable states it is a draw
	reversable_previous_state_strings []string
	previous_move                     Move
}

func generate_starting_game_state() GameState {
	new := GameState{white_mat: 39, black_mat: 39, is_end_state: false, turn_num: 1, curr_player: WHITE, duck: generate_default_duck()}

	for y := 2; y < 6; y++ {
		for x := 0; x < 8; x++ {
			new.board[y][x] = generate_empty_piece(x, y)
		}
	}

	// have important pieces first
	new.board[0][3] = Piece{x: 3, y: 0, kind: "queen", color: BLACK, worth: 9}
	new.board[0][4] = Piece{x: 4, y: 0, kind: "king", color: BLACK, worth: 100}
	new.board[0][0] = Piece{x: 0, y: 0, kind: "rook", color: BLACK, worth: 5}
	new.board[0][7] = Piece{x: 7, y: 0, kind: "rook", color: BLACK, worth: 5}
	new.board[0][1] = Piece{x: 1, y: 0, kind: "knight", color: BLACK, worth: 3}
	new.board[0][6] = Piece{x: 6, y: 0, kind: "knight", color: BLACK, worth: 3}
	new.board[0][2] = Piece{x: 2, y: 0, kind: "bishop", color: BLACK, worth: 3}
	new.board[0][5] = Piece{x: 5, y: 0, kind: "bishop", color: BLACK, worth: 3}

	new.board[7][3] = Piece{x: 3, y: 7, kind: "queen", color: WHITE, worth: 9}
	new.board[7][4] = Piece{x: 4, y: 7, kind: "king", color: WHITE, worth: 100}
	new.board[7][0] = Piece{x: 0, y: 7, kind: "rook", color: WHITE, worth: 5}
	new.board[7][7] = Piece{x: 7, y: 7, kind: "rook", color: WHITE, worth: 5}
	new.board[7][1] = Piece{x: 1, y: 7, kind: "knight", color: WHITE, worth: 3}
	new.board[7][6] = Piece{x: 6, y: 7, kind: "knight", color: WHITE, worth: 3}
	new.board[7][2] = Piece{x: 2, y: 7, kind: "bishop", color: WHITE, worth: 3}
	new.board[7][5] = Piece{x: 5, y: 7, kind: "bishop", color: WHITE, worth: 3}

	for i := 0; i < 8; i++ {
		new.board[1][i] = Piece{x: i, y: 1, kind: "pawn", color: BLACK, worth: 1}
		new.board[6][i] = Piece{x: i, y: 6, kind: "pawn", color: WHITE, worth: 1}
	}

	return new
}

// mutates calling state
func (state *GameState) swap_pieces(x1, y1, x2, y2 int) {

	temp := state.board[y2][x2]
	state.board[y2][x2] = state.board[y1][x1]
	state.board[y1][x1] = temp

	// update the piece's internal info, TODO maybe delete later
	state.board[y1][x1].x = x1
	state.board[y1][x1].y = y1
	state.board[y2][x2].x = x2
	state.board[y2][x2].y = y2
}

func is_3fold_rep(s *GameState) bool {
	count := 1
	curr_str := s.String()
	for i := range s.reversable_previous_state_strings {
		if s.reversable_previous_state_strings[i] == curr_str {
			count++
			if count == 3 {
				return true
			}

		}

	}

	return false
}

func is_50_move_limit(s *GameState) bool {
	return len(s.reversable_previous_state_strings) == 100
}

// mutates GameState
func (state *GameState) clear_reversable_state_strings() {
	state.reversable_previous_state_strings = nil
}

func evaluate_board(s *GameState) int {
	// Evaluate the current game state and return a heuristic value
	// Positive values favor the maximizing player (White), negative values favor the minimizing player (Black)
	if s.is_end_state {
		return s.result * OUTCOME_MULTIPLIER
	} else {
		return s.white_mat - s.black_mat
	}
}

func generate_possible_moves(s *GameState) []Move {
	// Generate and return a list of possible moves for the current game state
	var possible_moves []Move
	for y := range s.board {
		for x := range s.board[y] {
			piece := s.board[y][x]
			kind := piece.kind
			tempx := -1
			tempy := -1
			if piece.worth != 0 && piece.color == s.curr_player {

				if kind == "pawn" {
					// move 1 forward
					if y-s.curr_player >= 0 && is_empty(s, x, y-s.curr_player) {
						possible_moves = append(possible_moves, new_move(&piece, x, y, false, x, y-s.curr_player))
					}

					// move 2 forward
					var pawn_rank int
					if s.curr_player == BLACK {
						pawn_rank = 1
					} else {
						pawn_rank = 6
					}
					if (y == pawn_rank) && is_empty(s, x, y-s.curr_player) && is_empty(s, x, y-2*s.curr_player) {
						possible_moves = append(possible_moves, new_move(&piece, x, y, false, x, y-2*s.curr_player))
					}

					// take left
					if x > 0 && !is_empty(s, x-1, y-s.curr_player) && s.board[y-s.curr_player][x-1].color != s.curr_player {
						possible_moves = add_move(possible_moves, s, new_move(&piece, x, y, true, x-1, y-s.curr_player))
					}

					// take right
					if x < 7 && !is_empty(s, x+1, y-s.curr_player) && s.board[y-s.curr_player][x+1].color != s.curr_player {

						possible_moves = add_move(possible_moves, s, new_move(&piece, x, y, true, x+1, y-s.curr_player))

					}

					//en passant
					if can_passant(s, x, y) {
						possible_moves = append(possible_moves, new_special_move(&piece, x, y, true, s.previous_move.ex, y-s.curr_player, "", true, false))
					}

				} else if kind == "knight" {
					// north jump, left
					tempx = x - 1
					tempy = y - 2
					if in_bounds(tempx, tempy) && s.board[tempy][tempx].color != s.curr_player {
						possible_moves = add_move(possible_moves, s, new_move(&piece, x, y, false, tempx, tempy))
					}
					// north jump, right
					tempx = x + 1
					tempy = y - 2
					if in_bounds(tempx, tempy) && s.board[tempy][tempx].color != s.curr_player {
						possible_moves = add_move(possible_moves, s, new_move(&piece, x, y, false, tempx, tempy))
					}

					// east jump, north
					tempx = x + 2
					tempy = y - 1
					if in_bounds(tempx, tempy) && s.board[tempy][tempx].color != s.curr_player {
						possible_moves = add_move(possible_moves, s, new_move(&piece, x, y, false, tempx, tempy))
					}
					// east jump, south
					tempx = x + 2
					tempy = y + 1
					if in_bounds(tempx, tempy) && s.board[tempy][tempx].color != s.curr_player {
						possible_moves = add_move(possible_moves, s, new_move(&piece, x, y, false, tempx, tempy))
					}
					// south jump, east
					tempx = x + 1
					tempy = y + 2
					if in_bounds(tempx, tempy) && s.board[tempy][tempx].color != s.curr_player {
						possible_moves = add_move(possible_moves, s, new_move(&piece, x, y, false, tempx, tempy))
					}
					// south jump, west
					tempx = x - 1
					tempy = y + 2
					if in_bounds(tempx, tempy) && s.board[tempy][tempx].color != s.curr_player {
						possible_moves = add_move(possible_moves, s, new_move(&piece, x, y, false, tempx, tempy))
					}
					// west jump, south
					tempx = x - 2
					tempy = y + 1
					if in_bounds(tempx, tempy) && s.board[tempy][tempx].color != s.curr_player {
						possible_moves = add_move(possible_moves, s, new_move(&piece, x, y, false, tempx, tempy))
					}
					// west jump, north
					tempx = x - 2
					tempy = y - 1
					if in_bounds(tempx, tempy) && s.board[tempy][tempx].color != s.curr_player {
						possible_moves = add_move(possible_moves, s, new_move(&piece, x, y, false, tempx, tempy))
					}

				} else if kind == "bishop" {
					// move nw
					tempx = x - 1
					tempy = y - 1
					for in_bounds(tempx, tempy) {
						if s.board[tempy][tempx].kind == "empty" {
							possible_moves = add_move(possible_moves, s, new_move(&piece, x, y, false, tempx, tempy))

						} else if s.board[tempy][tempx].color != piece.color {
							possible_moves = add_move(possible_moves, s, new_move(&piece, x, y, true, tempx, tempy))
							break
						} else {
							break
						}
						tempx -= 1
						tempy -= 1
					}
					// move ne
					tempx = x + 1
					tempy = y - 1
					for in_bounds(tempx, tempy) {
						if s.board[tempy][tempx].kind == "empty" {
							possible_moves = add_move(possible_moves, s, new_move(&piece, x, y, false, tempx, tempy))

						} else if s.board[tempy][tempx].color != piece.color {
							possible_moves = add_move(possible_moves, s, new_move(&piece, x, y, true, tempx, tempy))
							break
						} else {
							break
						}
						tempx += 1
						tempy -= 1
					}
					// move se
					tempx = x + 1
					tempy = y + 1
					for in_bounds(tempx, tempy) {
						if s.board[tempy][tempx].kind == "empty" {
							possible_moves = add_move(possible_moves, s, new_move(&piece, x, y, false, tempx, tempy))

						} else if s.board[tempy][tempx].color != piece.color {
							possible_moves = add_move(possible_moves, s, new_move(&piece, x, y, true, tempx, tempy))
							break
						} else {
							break
						}
						tempx += 1
						tempy += 1
					}
					// move sw
					tempx = x - 1
					tempy = y + 1
					for in_bounds(tempx, tempy) {
						if s.board[tempy][tempx].kind == "empty" {
							possible_moves = add_move(possible_moves, s, new_move(&piece, x, y, false, tempx, tempy))

						} else if s.board[tempy][tempx].color != piece.color {
							possible_moves = add_move(possible_moves, s, new_move(&piece, x, y, true, tempx, tempy))
							break
						} else {
							break
						}
						tempx -= 1
						tempy += 1
					}

				} else if kind == "rook" {
					// move west
					tempx = x - 1
					tempy = y
					for in_bounds(tempx, tempy) {
						if s.board[tempy][tempx].kind == "empty" {
							possible_moves = add_move(possible_moves, s, new_move(&piece, x, y, false, tempx, tempy))

						} else if s.board[tempy][tempx].color != piece.color {
							possible_moves = add_move(possible_moves, s, new_move(&piece, x, y, true, tempx, tempy))
							break
						} else {
							break
						}
						tempx -= 1
					}
					// move north
					tempx = x
					tempy = y - 1
					for in_bounds(tempx, tempy) {
						if s.board[tempy][tempx].kind == "empty" {
							possible_moves = add_move(possible_moves, s, new_move(&piece, x, y, false, tempx, tempy))

						} else if s.board[tempy][tempx].color != piece.color {
							possible_moves = add_move(possible_moves, s, new_move(&piece, x, y, true, tempx, tempy))
							break
						} else {
							break
						}
						tempy -= 1
					}
					// move east
					tempx = x + 1
					tempy = y
					for in_bounds(tempx, tempy) {
						if s.board[tempy][tempx].kind == "empty" {
							possible_moves = add_move(possible_moves, s, new_move(&piece, x, y, false, tempx, tempy))

						} else if s.board[tempy][tempx].color != piece.color {
							possible_moves = add_move(possible_moves, s, new_move(&piece, x, y, true, tempx, tempy))
							break
						} else {
							break
						}
						tempx += 1
					}
					// move south
					tempx = x
					tempy = y + 1
					for in_bounds(tempx, tempy) {
						if s.board[tempy][tempx].kind == "empty" {
							possible_moves = add_move(possible_moves, s, new_move(&piece, x, y, false, tempx, tempy))

						} else if s.board[tempy][tempx].color != piece.color {
							possible_moves = add_move(possible_moves, s, new_move(&piece, x, y, true, tempx, tempy))
							break
						} else {
							break
						}
						tempy += 1
					}

				} else if kind == "king" {
					//TODO
				} else if kind == "queen" {
					// move west
					tempx = x - 1
					tempy = y
					for in_bounds(tempx, tempy) {
						if s.board[tempy][tempx].kind == "empty" {
							possible_moves = add_move(possible_moves, s, new_move(&piece, x, y, false, tempx, tempy))

						} else if s.board[tempy][tempx].color != piece.color {
							possible_moves = add_move(possible_moves, s, new_move(&piece, x, y, true, tempx, tempy))
							break
						} else {
							break
						}
						tempx -= 1
					}
					// move north
					tempx = x
					tempy = y - 1
					for in_bounds(tempx, tempy) {
						if s.board[tempy][tempx].kind == "empty" {
							possible_moves = add_move(possible_moves, s, new_move(&piece, x, y, false, tempx, tempy))

						} else if s.board[tempy][tempx].color != piece.color {
							possible_moves = add_move(possible_moves, s, new_move(&piece, x, y, true, tempx, tempy))
							break
						} else {
							break
						}
						tempy -= 1
					}
					// move east
					tempx = x + 1
					tempy = y
					for in_bounds(tempx, tempy) {
						if s.board[tempy][tempx].kind == "empty" {
							possible_moves = add_move(possible_moves, s, new_move(&piece, x, y, false, tempx, tempy))

						} else if s.board[tempy][tempx].color != piece.color {
							possible_moves = add_move(possible_moves, s, new_move(&piece, x, y, true, tempx, tempy))
							break
						} else {
							break
						}
						tempx += 1
					}
					// move south
					tempx = x
					tempy = y + 1
					for in_bounds(tempx, tempy) {
						if s.board[tempy][tempx].kind == "empty" {
							possible_moves = add_move(possible_moves, s, new_move(&piece, x, y, false, tempx, tempy))

						} else if s.board[tempy][tempx].color != piece.color {
							possible_moves = add_move(possible_moves, s, new_move(&piece, x, y, true, tempx, tempy))
							break
						} else {
							break
						}
						tempy += 1
					}
					// move nw
					tempx = x - 1
					tempy = y - 1
					for in_bounds(tempx, tempy) {
						if s.board[tempy][tempx].kind == "empty" {
							possible_moves = add_move(possible_moves, s, new_move(&piece, x, y, false, tempx, tempy))

						} else if s.board[tempy][tempx].color != piece.color {
							possible_moves = add_move(possible_moves, s, new_move(&piece, x, y, true, tempx, tempy))
							break
						} else {
							break
						}
						tempx -= 1
						tempy -= 1
					}
					// move ne
					tempx = x + 1
					tempy = y - 1
					for in_bounds(tempx, tempy) {
						if s.board[tempy][tempx].kind == "empty" {
							possible_moves = add_move(possible_moves, s, new_move(&piece, x, y, false, tempx, tempy))

						} else if s.board[tempy][tempx].color != piece.color {
							possible_moves = add_move(possible_moves, s, new_move(&piece, x, y, true, tempx, tempy))
							break
						} else {
							break
						}
						tempx += 1
						tempy -= 1
					}
					// move se
					tempx = x + 1
					tempy = y + 1
					for in_bounds(tempx, tempy) {
						if s.board[tempy][tempx].kind == "empty" {
							possible_moves = add_move(possible_moves, s, new_move(&piece, x, y, false, tempx, tempy))

						} else if s.board[tempy][tempx].color != piece.color {
							possible_moves = add_move(possible_moves, s, new_move(&piece, x, y, true, tempx, tempy))
							break
						} else {
							break
						}
						tempx += 1
						tempy += 1
					}
					// move sw
					tempx = x - 1
					tempy = y + 1
					for in_bounds(tempx, tempy) {
						if s.board[tempy][tempx].kind == "empty" {
							possible_moves = add_move(possible_moves, s, new_move(&piece, x, y, false, tempx, tempy))

						} else if s.board[tempy][tempx].color != piece.color {
							possible_moves = add_move(possible_moves, s, new_move(&piece, x, y, true, tempx, tempy))
							break
						} else {
							break
						}
						tempx -= 1
						tempy += 1
					}

				} else {
					panic("unrecognized piece kind on the board")
				}
			}

		}
	}

	return possible_moves
}

func add_move(list []Move, s *GameState, m Move) []Move {
	if s.board[m.ey][m.ex].kind == "empty" {
		return append(list, m)
	} else {
		m.cap = true
		if s.board[m.ey][m.ex].kind == "king" {
			m.extra += "#"
		}
		return append(list, m)
	}
}

func in_bounds(x, y int) bool {
	return x > -1 && x < 8 && y > -1 && y < 8
}

func is_empty(s *GameState, x, y int) bool {
	return s.board[y][x].kind == "empty"
}

func can_passant(s *GameState, x, y int) bool {
	if s.turn_num > 1 {
		was_pawn := s.previous_move.piece.kind == "pawn"
		was_double := abs(s.previous_move.sy-s.previous_move.ey) == 2
		is_adjacent := abs(s.previous_move.ex-x) == 1
		is_level := s.previous_move.ey == y

		return was_pawn && was_double && is_adjacent && is_level
	}

	return false
}

func (s GameState) String() string {

	var str string
	for y := range s.board {
		for x := range s.board[y] {
			str += piece_to_unicode(s.board[y][x]) + " "
		}
		str += "\n"
	}
	return str
}
