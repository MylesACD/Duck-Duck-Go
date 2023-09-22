package main

import (
	"fmt"
	_ "net/http/pprof"
)

func minimax(state GameState, depth int, maximizingPlayer bool) int {
	/*
		if isTerminalNode(state, depth) {
			return evaluateBoard(state)
		}

		if maximizingPlayer {
			bestValue := math.MinInt
			possibleMoves := generatePossibleMoves(state)

			for _, move := range possibleMoves {
				// Apply the move to the current state
				// newGameState := applyMove(state, move)

				value := minimax(newGameState, depth-1, false)
				bestValue = max(bestValue, value)
			}

			return bestValue
		} else {
			bestValue := math.MaxInt
			possibleMoves := generatePossibleMoves(state)

			for _, move := range possibleMoves {
				// Apply the move to the current state
				// newGameState := applyMove(state, move)

				value := minimax(newGameState, depth-1, true)
				bestValue = min(bestValue, value)
			}

			return bestValue
		}
	*/
	return 0
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	// Initialize the game state
	initialState := get_starting_game_state() // Initialize with the starting position
	fmt.Println(initialState)

}
