package app

import "fmt"

type Cell struct {
	value *Value
}

type Board [][]Cell

func (b Board) Len() int { return len(b) }

func NewBoard(n int) Board {
	board := [][]Cell{}
	for range n {
		row := []Cell{}
		for range n {
			row = append(row, Cell{value: nil})
		}
		board = append(board, row)
	}
	return board
}

func (b *Board) Play(row int, col int, value Value) error {
	if (*b)[row][col].value == nil {
		(*b)[row][col].value = &value
		return nil
	}
	return fmt.Errorf("Cannot play at row: %d, col: %d, already taken", row, col)
}

func Win(cellList []Cell, turn Value) bool {
	for i := range len(cellList) {
		if cellList[i].value == nil {
			return false
		}
		if *cellList[i].value != turn {
			return false
		}
	}

	return true
}

func (b Board) Win(turn Value) bool {
	for i := range b.Len() {
		if Win(b[i], turn) {
			return true
		}
	}

	for i := range b.Len() {
		columns := []Cell{}
		for j := range b.Len() {
			columns = append(columns, b[j][i])
		}
		if Win(columns, turn) {
			return true
		}
	}

	diagonal1 := []Cell{}
	diagonal2 := []Cell{}
	for i := range b.Len() {
		diagonal1 = append(diagonal1, b[i][i])
		diagonal2 = append(diagonal2, b[i][b.Len()-1-i])
	}
	if Win(diagonal1, turn) || Win(diagonal2, turn) {
		return true
	}

	return false
}

func (b Board) Print() {
	for i := range b.Len() {
		for j := range b.Len() {
			value := " "
			if b[i][j].value != nil {
				value = string(*b[i][j].value)
			}
			if j == 0 {
				value = " " + value
			}
			fmt.Print(value)
			if j < b.Len()-1 {
				fmt.Print("|")
			}
		}
		if i < b.Len()-1 {
			character := ""
			for range b.Len() {
				character += "--"
			}
			character = string(character[:len(character)-1])
			fmt.Println("\n", character)
		}
	}
	fmt.Println()
	fmt.Println()
}

func (b Board) Filled() bool {
	for i := range b.Len() {
		for j := range b.Len() {
			if b[i][j].value == nil {
				return false
			}
		}
	}

	return true
}
