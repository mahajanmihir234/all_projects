package app

import "fmt"

type TicTacToe struct {
	Board Board
	Turn  Value
}

func NewTicTacToe(n int) TicTacToe { return TicTacToe{Board: NewBoard(n), Turn: X} }

func (t *TicTacToe) Play(row int, col int, value Value) error {
	if t.Turn != value {
		return fmt.Errorf("Incorrect turn")
	}
	err := t.Board.Play(row, col, value)
	if err != nil {
		return err
	}
	t.Turn = t.Turn.Opposite()
	return nil
}

func (t TicTacToe) Win() bool {
	return t.Board.Win(t.Turn.Opposite())
}

func (t TicTacToe) Print() {
	t.Board.Print()
}
