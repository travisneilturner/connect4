package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"travisneilturner/connect4/internal/game"

	termutil "github.com/andrew-d/go-termutil"
)

func main() {
	board, err := game.NewBoard(7, 6)
	if err != nil {
		fmt.Println("couldn't init board: " + err.Error())
		os.Exit(-1)
	}

	if termutil.Isatty(os.Stdin.Fd()) {
		doInteractiveGame(board)
		os.Exit(0)
	}

	doNonInteractiveGame(board)
	os.Exit(0)
}

func doNonInteractiveGame(board *game.Board) {
	reader := bufio.NewScanner(os.Stdin)

	// non-interactive mode expects input on SStdin so if we get EOF it's an error
	if !reader.Scan() {
		usage()
		os.Exit(-1)
	}

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

	result, err := board.AddAllPieces(intMoves)
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
		return
	}

	fmt.Println("WINNER: Player " + strconv.Itoa(result.Winner))
}

func doInteractiveGame(board *game.Board) {
	player := 1
	reader := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Player " + strconv.Itoa(player) + " enter a column: ")
		if !reader.Scan() {
			usage()
			os.Exit(-1)
		}

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
				return
			}

			fmt.Println("WINNER: Player " + strconv.Itoa(player))
			return
		}

		if player == 1 {
			player = 2
		} else {
			player = 1
		}
	}
}

func usage() {
	fmt.Println("For interactive mode, do:")
	fmt.Println("> docker run -it travisneilturner/connect4")
	fmt.Println()
	fmt.Println("For non-interactive mode (no TTY attached), do:")
	fmt.Println("> echo [list of space-delimited columns] | docker run -i travisneilturner/connect4")
	fmt.Println()
	fmt.Println("Example: ")
	fmt.Println("> echo 0 1 0 1 0 1 0 | docker run -i travisneilturner/connect4")
	fmt.Println()
	fmt.Println("You can also run the image then enter a list of moves, like so:")
	fmt.Println("> docker run -i travisneilturner/connect4")
	fmt.Println("> 0 1 0 1 0 1 0")
	fmt.Println()
}
