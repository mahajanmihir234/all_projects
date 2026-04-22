package game

import (
	"errors"
	"fmt"
)

type Cell *Symbol

type Board struct {
	grid [][]Cell
	size int
}

func NewBoard(size int) Board {
	grid := [][]Cell{}
	for range size {
		row := []Cell{}
		for range size {
			row = append(row, nil)
		}
		grid = append(grid, row)
	}

	return Board{
		grid: grid,
		size: size,
	}
}

func (b Board) Play(turn Symbol, row int, col int) error {
	if b.grid[row][col] != nil {
		return errors.New("Cannot play in an already filled cell")
	}

	b.grid[row][col] = &turn
	return nil
}

func (b Board) Filled() bool {
	for i := range b.size {
		for j := range b.size {
			if b.grid[i][j] == nil {
				return false
			}
		}
	}
	return true
}

func (b Board) Print() {
	for i := range b.size {
		for j := range b.size {
			value := " "
			if b.grid[i][j] != nil {
				value = string(*b.grid[i][j])
			}
			if j == 0 {
				value = " " + value
			}
			fmt.Print(value)
			if j < b.size-1 {
				fmt.Print("|")
			}
		}
		if i < b.size-1 {
			character := ""
			for range b.size {
				character += "--"
			}
			character = string(character[:len(character)-1])
			fmt.Println("\n", character)
		}
	}
	fmt.Println()
	fmt.Println()
}
