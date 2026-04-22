package game

import (
	"errors"
	"fmt"
)

type Game struct {
	board      Board
	turn       Symbol
	strategies []WinningStrategy
}

func NewGame(size int) Game {
	return Game{
		board: NewBoard(size),
		turn:  X,
		strategies: []WinningStrategy{
			RowWinningStrategy{}, ColumnWinningStrategy{}, Diagonal1WinningStrategy{}, Diagonal2WinningStrategy{},
		},
	}
}

func (g Game) Play(turn Symbol, row int, col int) (bool, error) {
	if g.turn != turn {
		return false, errors.New("Invalid turn")
	}

	err := g.board.Play(turn, row, col)
	if err != nil {
		return false, err
	}

	for _, ws := range g.strategies {
		if ws.Win(g.board, turn, row, col) {
			return true, nil
		}
	}

	if g.board.Filled() {
		return true, errors.New("DRAW")
	}
	g.turn = g.turn.Opposite()
	fmt.Println("NewTurn", g.turn)

	return false, nil
}

func (g Game) Turn() Symbol {
	return g.turn
}

func (g Game) Print() {
	g.board.Print()
}
