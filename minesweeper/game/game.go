package game

import "fmt"

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Gray = "\033[37m"

func NewBoard(mineGenerator MineGenerator, rows int, cols int, numMines int) Board {
	board := initBoard(rows, cols)
	board.AddMines(mineGenerator.GenerateMines(rows, cols, numMines))
	return board
}

func initBoard(rows int, cols int) Board {
	board := make([][]int, rows)
	for i := range rows {
		board[i] = make([]int, cols)
	}
	for i := range rows {
		for j := range cols {
			board[i][j] = -1
		}
	}
	return board
}

type MineSweeper struct {
	board    Board
	revealed map[Coordinate]bool
	marked   map[Coordinate]bool
}

func NewMineSweeper(mineGenerator MineGenerator, rows int, cols int, numMines int) MineSweeper {
	board := NewBoard(mineGenerator, rows, cols, numMines)
	return MineSweeper{
		board:    board,
		revealed: map[Coordinate]bool{},
		marked:   map[Coordinate]bool{},
	}
}

func (m MineSweeper) Mark(row int, col int) {
	coordinate := Coordinate{row: row, col: col}
	if m.revealed[coordinate] {
		return
	}
	m.marked[Coordinate{row: row, col: col}] = true
}

func (m MineSweeper) Unmark(row int, col int) {
	coordinate := Coordinate{row: row, col: col}
	if m.revealed[coordinate] {
		return
	}
	if m.marked[coordinate] {
		m.marked[coordinate] = false
	}
}

func (m MineSweeper) Print(revealAll bool) {
	fmt.Printf(" ")
	for row := range m.board.Cols() {
		fmt.Printf(" | %d", row)
	}
	fmt.Println()
	for range m.board.Cols() {
		fmt.Printf("----")
	}
	fmt.Printf("--")
	fmt.Println()

	if revealAll {
		for row := range m.board.Rows() {
			fmt.Printf("%d", row)
			for col := range m.board.Cols() {
				value := fmt.Sprintf(Green+"%d"+Reset, m.board[row][col])
				if m.board[row][col] == 0 {
					value = fmt.Sprint(Red + "B" + Reset)
				}
				if m.board[row][col] == -1 {
					value = fmt.Sprint(Red + " " + Reset)
				}
				fmt.Printf(" | %s", value)
			}
			fmt.Println()
			for range m.board.Cols() {
				fmt.Printf("----")
			}
			fmt.Printf("--")
			fmt.Println()
		}
		return
	}

	for row := range m.board.Rows() {
		fmt.Printf("%d", row)
		for col := range m.board.Cols() {
			coordinate := Coordinate{row: row, col: col}
			value := fmt.Sprint(Gray + "?" + Reset)
			if found, ok := m.marked[coordinate]; ok && found {
				value = fmt.Sprint(Red + "M" + Reset)
			} else if _, ok := m.revealed[coordinate]; ok {
				if m.board[row][col] != -1 {
					value = fmt.Sprintf(Green+"%d"+Reset, m.board[row][col])
				} else {
					value = fmt.Sprint(Gray + " " + Reset)
				}
			}
			fmt.Printf(" | %s", value)
		}
		fmt.Println()
		for range m.board.Cols() {
			fmt.Printf("----")
		}
		fmt.Printf("--")
		fmt.Println()
	}
}

func (m MineSweeper) Reveal(row int, col int) error {
	if m.board[row][col] == 0 {
		m.Print(true)
		return fmt.Errorf("GAME OVER")
	}

	coordinate := Coordinate{row: row, col: col}
	if m.board[row][col] != -1 {
		m.revealed[coordinate] = true
		return nil
	}

	queue := Queue{}
	queue.Push(coordinate, m.board.Rows(), m.board.Cols())

	visited := map[Coordinate]bool{}
	for len(queue) > 0 {
		element := queue.Pop()
		if visited[element] {
			continue
		}
		m.revealed[element] = true
		visited[element] = true
		if m.board[element.row][element.col] != -1 {
			continue
		}
		queue.Push(Coordinate{row: element.row, col: element.col + 1}, m.board.Rows(), m.board.Cols())
		queue.Push(Coordinate{row: element.row, col: element.col - 1}, m.board.Rows(), m.board.Cols())
		queue.Push(Coordinate{row: element.row + 1, col: element.col}, m.board.Rows(), m.board.Cols())
		queue.Push(Coordinate{row: element.row + 1, col: element.col - 1}, m.board.Rows(), m.board.Cols())
		queue.Push(Coordinate{row: element.row + 1, col: element.col + 1}, m.board.Rows(), m.board.Cols())
		queue.Push(Coordinate{row: element.row - 1, col: element.col}, m.board.Rows(), m.board.Cols())
		queue.Push(Coordinate{row: element.row - 1, col: element.col - 1}, m.board.Rows(), m.board.Cols())
		queue.Push(Coordinate{row: element.row - 1, col: element.col + 1}, m.board.Rows(), m.board.Cols())
	}
	return nil
}

func (m MineSweeper) Play(moveType string, row int, col int) error {
	if moveType == "MARK" {
		m.Mark(row, col)
		return nil
	}
	if moveType == "UNMARK" {
		m.Unmark(row, col)
		return nil
	}
	return m.Reveal(row, col)
}

func (m MineSweeper) Win() bool {
	for row := range m.board.Rows() {
		for col := range m.board.Cols() {
			coordinate := Coordinate{row: row, col: col}
			if m.board[row][col] == 0 {
				if !m.marked[coordinate] {
					return false
				}
			} else if !m.revealed[coordinate] {
				return false
			}
		}
	}

	return true
}
