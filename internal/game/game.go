package game

import (
	"errors"
	"strconv"
)

// Board represents the state of the game
type Board struct {
	// State is the 2-D array containing the placed pieces
	State [][]int
	// Top is an array of numbers representing the "lowest" (in physical terms)
	// available row per column.  This corresponds to the maximum unoccupied index in
	// the State array
	Top []int
	// Finished allows us to quickly check to make sure that no more pieces are added
	// to a board that has a winner
	Finished bool
}

// Outcome represents the outcome of placing pieces on the Board
// A nil Outcome is used to represent the successful completion of a
// move that does not produce a winner
type Outcome struct {
	// Winner is the player's number (or -1 for a draw)
	Winner int
}

// NewBoard returns a new board with the provided dimensions
func NewBoard(width, height int) (*Board, error) {
	if width < 4 {
		return nil, errors.New("width must be positive and at least 4 (this is connect 4 after all)")
	}

	if height < 4 {
		return nil, errors.New("height must be positive and at least 4 (this is connect 4 after all)")
	}

	rows := make([][]int, height)
	for i := range rows {
		rows[i] = make([]int, width)
		for j := 0; j < width; j++ {
			rows[i][j] = -1
		}
	}

	// we place pieces at the maximum available index per column
	// and keep track of that value in "top"
	top := make([]int, width)
	for i := range top {
		top[i] = height - 1
	}

	return &Board{
		State: rows,
		Top:   top,
	}, nil
}

// AddPiece adds the piece corresponding to "player" to the column specified by "col"
/* Returns:
	nil, err if there was an invalid move
	*Outcome, nil if the move results in a win or a draw
	nil, nil otherwise

The *Outcome contains the Winner's id (or -1 for a draw)
*/
func (b *Board) AddPiece(col int, player int) (*Outcome, error) {
	if b.Finished {
		return nil, errors.New("can't place more pieces on a completed board")
	}

	if col < 0 || col > len(b.State[0])-1 {
		return nil, errors.New("invalid move: column must be between 0 and " + strconv.Itoa(len(b.State[0])-1))
	}

	if b.Top[col] == -1 {
		return nil, errors.New("invalid move: column " + strconv.Itoa(col) + " is full")
	}

	b.State[b.Top[col]][col] = player
	b.Top[col]--

	return b.checkWinner(col), nil
}

// AddAllPieces adds all the pieces to the board one by one and stops if there is a
// winner, a draw, or an error
func (b *Board) AddAllPieces(cols []int) (*Outcome, error) {
	player := 1
	for i, col := range cols {
		outcome, err := b.AddPiece(col, player)
		if err != nil {
			return nil, errors.New("move at pos " + strconv.Itoa(i+1) + " caused an error: " + err.Error())
		}

		if outcome != nil {
			return outcome, nil
		}

		if player == 1 {
			player = 2
		} else {
			player = 1
		}
	}

	return nil, nil
}

// checkWinner performs a 7-directional search to determine whether there
// are any winning combinations originating from the most recent move
func (b *Board) checkWinner(col int) *Outcome {
	// start of the search is the row just below the new top of the column
	// so we add 1 to the row index
	row := b.Top[col] + 1
	player := b.State[row][col]

	// first check right and left
	rightMatches := b.checkDirection(player, row, col, 0, 1, 3)
	if rightMatches == 3 {
		b.Finished = true
		return &Outcome{
			Winner: player,
		}
	}

	leftMatches := b.checkDirection(player, row, col, 0, -1, 3-rightMatches)
	if rightMatches+leftMatches >= 3 {
		b.Finished = true
		return &Outcome{
			Winner: player,
		}
	}

	// next check down (since we won't have any pieces above us)
	// (remember that "down" corresponds to a HIGHER index so the rowDir is +1)
	downMatches := b.checkDirection(player, row, col, 1, 0, 3)
	if downMatches == 3 {
		b.Finished = true
		return &Outcome{
			Winner: player,
		}
	}

	// then check up right and down left
	upRightMatches := b.checkDirection(player, row, col, -1, 1, 3)
	if upRightMatches == 3 {
		b.Finished = true
		return &Outcome{
			Winner: player,
		}
	}

	downLeftMatches := b.checkDirection(player, row, col, 1, -1, 3-upRightMatches)
	if upRightMatches+downLeftMatches >= 3 {
		b.Finished = true
		return &Outcome{
			Winner: player,
		}
	}

	// finally check up left and down right
	upLeftMatches := b.checkDirection(player, row, col, -1, -1, 3)
	if upLeftMatches == 3 {
		b.Finished = true
		return &Outcome{
			Winner: player,
		}
	}

	downRightMatches := b.checkDirection(player, row, col, 1, 1, 3-upLeftMatches)
	if upLeftMatches+downRightMatches >= 3 {
		b.Finished = true
		return &Outcome{
			Winner: player,
		}
	}

	// last but not least check if all the columns are full, if so we have a draw
	full := true
	for i := range b.Top {
		if b.Top[i] != -1 {
			full = false
			break
		}
	}

	if full {
		b.Finished = true
		return &Outcome{
			Winner: -1,
		}
	}

	return nil
}

/**
checkDirection scans in one of the 8 cardinal directions specified by [rowDir, colDir] at most "stop" units
starting from [rowStart, colStart) and returns the number of pieces that match "player"
*/
func (b *Board) checkDirection(player int, rowStart int, colStart int, rowDir int, colDir int, stop int) int {
	matches := 0

	row := rowStart
	col := colStart
	for i := 0; i < stop; i++ {
		row = row + rowDir
		col = col + colDir

		if row >= 0 && row < len(b.State) && col >= 0 && col < len(b.State[0]) && b.State[row][col] == player {
			matches++
		} else {
			return matches
		}
	}

	return matches
}
