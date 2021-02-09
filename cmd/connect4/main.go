package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"travisneilturner/connect4/internal/game"
)

func main() {
	board, err := game.NewBoard(7, 6)
	if err != nil {
		fmt.Println("couldn't init board: " + err.Error())
		os.Exit(-1)
	}

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		doNonInteractiveGame(board)
	}

	doInteractiveGame(board)
	os.Exit(0)
}

func doNonInteractiveGame(board *game.Board) {
	player := 1
	reader := bufio.NewScanner(os.Stdin)
	reader.Scan()
	moves := strings.Split(reader.Text(), " ")

	intMoves := make([]int, len(moves))
	for i, mv := range moves {
		val, err := strconv.Atoi(mv)
		if err != nil {
			fmt.Println("move in pos " + strconv.Itoa(i) + " is invalid, must be a number from 0 to " + strconv.Itoa(len(board.State[0])-1))
			os.Exit(-1)
		}
		intMoves[i] = val
	}

	result, err := board.AddAllPieces(intMoves, player)
	if err != nil {
		fmt.Println("error playing game: " + err.Error())
		os.Exit(-1)
	}

	if result == nil {
		fmt.Println("incomplete game -- did not result in a winner or a draw")
		os.Exit(-1)
	}

	if result.Winner == -1 {
		fmt.Println("DRAW")
		os.Exit(0)
	}

	fmt.Println("WINNER: Player " + strconv.Itoa(result.Winner))
	os.Exit(0)
}

func doInteractiveGame(board *game.Board) {
	player := 1
	reader := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Player " + strconv.Itoa(player) + " enter a column: ")
		reader.Scan()

		if len(reader.Text()) == 0 {
			fmt.Println("Please enter a number between 0 and " + strconv.Itoa(len(board.State[0])-1))
			continue
		}

		col, err := strconv.Atoi(reader.Text())
		if err != nil {
			fmt.Println("Error with input, please enter a number between 0 and " + strconv.Itoa(len(board.State[0])-1))
			continue
		}

		outcome, err := board.AddPiece(col, player)
		if err != nil {
			fmt.Println("Error processing move: " + err.Error())
			continue
		}

		if outcome != nil {
			if outcome.Winner == -1 {
				fmt.Println("DRAW")
				os.Exit(0)
			}

			fmt.Println("WINNER: Player " + strconv.Itoa(player))
			os.Exit(0)
		}

		if player == 1 {
			player = 2
		} else {
			player = 1
		}
	}
}
