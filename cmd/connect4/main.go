package main

import (
	"bufio"
	"errors"
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
		err = doInteractiveGame(board)
	} else {
		err = doNonInteractiveGame(board)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	os.Exit(0)
}

func doNonInteractiveGame(board *game.Board) error {
	reader := bufio.NewScanner(os.Stdin)

	// non-interactive mode expects input on Stdin so if we get false value from scanner
	// we can assume the image wasn't invoked correctly
	if !reader.Scan() {
		usage()
		// only return a more verbose error if something *really* bad happens
		// if we get an EOF it's better to just print the usage and exit
		if reader.Err() != nil {
			return errors.New("error while scanning: " + reader.Err().Error())
		}
		return nil
	}

	moves := strings.Split(reader.Text(), " ")

	intMoves := make([]int, len(moves))
	for i, mv := range moves {
		val, err := strconv.Atoi(mv)
		if err != nil {
			return errors.New("input \"" + mv + "\" in pos " + strconv.Itoa(i+1) + " is invalid, must be a number from 0 to " + strconv.Itoa(len(board.State[0])-1))
		}
		intMoves[i] = val
	}

	result, err := board.AddAllPieces(intMoves)
	if err != nil {
		return errors.New("error playing game: " + err.Error())
	}

	if result == nil {
		return errors.New("incomplete game -- did not result in a winner or a draw")
	}

	if result.Winner == -1 {
		fmt.Println("DRAW")
		return nil
	}

	fmt.Println("WINNER: Player " + strconv.Itoa(result.Winner))
	return nil
}

func doInteractiveGame(board *game.Board) error {
	player := 1
	reader := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Player " + strconv.Itoa(player) + " enter a column: ")

		// In interactive mode, any scanning errors are unexpected and will
		// cause the game to end with the error encountered.  We don't print
		// the usage because presumably if they managed to get an interactive session going
		// they have figured out how to attach a TTY
		if !reader.Scan() {
			if reader.Err() != nil {
				return errors.New("error while scanning: " + reader.Err().Error())
			}
			return errors.New("EOF while scanning")
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
				return nil
			}

			fmt.Println("WINNER: Player " + strconv.Itoa(player))
			return nil
		}

		if player == 1 {
			player = 2
		} else {
			player = 1
		}
	}
}

func usage() {
	str := `For interactive mode, do: 
> docker run -it travisneilturner/connect4

For non-interactive mode (no TTY attached), do:
> echo [list of space-delimited columns] | docker run -i travisneilturner/connect4

Example: 
> echo 0 1 0 1 0 1 0 | docker run -i travisneilturner/connect4

You can also run the image then enter a list of moves, like so:
> docker run -i travisneilturner/connect4
> 0 1 0 1 0 1 0
`
	fmt.Println(str)
}
