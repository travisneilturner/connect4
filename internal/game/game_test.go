package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBoard(t *testing.T) {
	board, err := NewBoard(7, 6)
	assert.Nil(t, err, "should be no error")
	assert.NotEmpty(t, board, "board should exist")
	assert.Equal(t, 7, len(board.State[0]), "width should be 7")
	assert.Equal(t, 6, len(board.State), "height should be 6")
}

func TestNewInvalidBoard(t *testing.T) {
	_, err := NewBoard(-7, 6)
	assert.NotNil(t, err, "should be an error")
}

func TestAddPiece(t *testing.T) {
	board, _ := NewBoard(10, 10)

	board.AddPiece(0, 1)
	assert.Equal(t, 8, board.Top[0], "highest available row in col 0 should be 8")

	board.AddPiece(0, 1)
	board.AddPiece(0, 1)
	assert.Equal(t, 6, board.Top[0], "highest available row in col 0 should be 6")

	expectedState := [][]int{
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	}
	assert.Equal(t, expectedState, board.State, "states should match")
}

func TestAddPieceError(t *testing.T) {
	board, _ := NewBoard(10, 10)
	for i := 0; i < 10; i++ {
		board.AddPiece(0, 1)
	}

	_, err := board.AddPiece(0, 1)
	assert.NotNil(t, err, "shouldn't be able to place 11 pieces in a column that is only 10 spaces tall")
}

func TestLeftRightWin(t *testing.T) {
	board, _ := NewBoard(7, 6)
	var outcome *Outcome
	var err error
	for i := 0; i < 3; i++ {
		outcome, err = board.AddPiece(i, 0)
		assert.Nil(t, outcome, "AddPiece should return nil outcome")
		assert.Nil(t, err, "AddPiece should return no error")
	}

	outcome, _ = board.AddPiece(3, 0)
	assert.Equal(t, &Outcome{
		Winner: 0,
	}, outcome, "player 0 should be the winner")
}

func TestRightLeftWin(t *testing.T) {
	board, _ := NewBoard(7, 6)
	var outcome *Outcome
	for i := 3; i >= 1; i-- {
		outcome, _ = board.AddPiece(i, 0)
	}

	outcome, _ = board.AddPiece(0, 0)
	assert.Equal(t, &Outcome{
		Winner: 0,
	}, outcome, "player 0 should be the winner")
}

func TestLeftRightMiddleWin(t *testing.T) {
	board, _ := NewBoard(7, 6)
	board.AddPiece(0, 0)
	board.AddPiece(2, 0)
	board.AddPiece(3, 0)

	outcome, _ := board.AddPiece(1, 0)
	assert.Equal(t, &Outcome{
		Winner: 0,
	}, outcome, "player 0 should be the winner")
}

func TestUpDownWin(t *testing.T) {
	board, _ := NewBoard(7, 6)
	board.AddPiece(0, 0)
	board.AddPiece(0, 0)
	board.AddPiece(0, 0)

	outcome, _ := board.AddPiece(0, 0)
	assert.Equal(t, &Outcome{
		Winner: 0,
	}, outcome, "player 0 should be the winner")

}

/*

1
0 1 0
1 0 1   0
0 0 1 1 1
*/
func TestUpLeftWin(t *testing.T) {
	board, _ := NewBoard(7, 6)
	board.AddPiece(0, 0)
	board.AddPiece(0, 1)
	board.AddPiece(0, 0)
	board.AddPiece(0, 1)
	board.AddPiece(1, 0)
	board.AddPiece(4, 1)
	board.AddPiece(1, 0)
	board.AddPiece(1, 1)
	board.AddPiece(2, 0)
	board.AddPiece(2, 1)
	board.AddPiece(4, 0)

	outcome, _ := board.AddPiece(3, 1)
	assert.Equal(t, &Outcome{
		Winner: 1,
	}, outcome, "player 1 should be the winner")
}

func TestDraw(t *testing.T) {
	board, _ := NewBoard(4, 4)
	board.AddPiece(0, 0)
	board.AddPiece(0, 1)
	board.AddPiece(0, 0)
	board.AddPiece(0, 1)
	board.AddPiece(1, 0)
	board.AddPiece(1, 1)
	board.AddPiece(1, 0)
	board.AddPiece(1, 1)
	board.AddPiece(2, 0)
	board.AddPiece(3, 1)
	board.AddPiece(2, 0)
	board.AddPiece(2, 1)
	board.AddPiece(2, 0)
	board.AddPiece(3, 1)
	board.AddPiece(3, 0)

	outcome, _ := board.AddPiece(3, 1)
	assert.Equal(t, &Outcome{
		Winner: -1,
	}, outcome, "outcome should be a draw")
}

func TestAddAllPieces(t *testing.T) {
	board, _ := NewBoard(4, 4)
	outcome, _ := board.AddAllPieces([]int{0, 1, 0, 1, 0, 1, 0})

	assert.Equal(t, &Outcome{
		Winner: 1,
	}, outcome, "player 1 should be the winner")
}
