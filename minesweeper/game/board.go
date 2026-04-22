package game

type Coordinate struct {
	row int
	col int
}

type Board [][]int

func (b Board) Rows() int {
	return len(b)
}

func (b Board) Cols() int {
	return len(b[0])
}

func (b *Board) UpdateValue(row int, col int) {
	if row < 0 || col < 0 || row == b.Rows() || col == b.Cols() {
		return
	}

	if (*b)[row][col] == 0 {
		return
	}

	if (*b)[row][col] == -1 {
		(*b)[row][col] = 1
	} else {
		(*b)[row][col] += 1
	}
}

func (b *Board) AddMines(mineCoordinates map[Coordinate]bool) {
	for coordinate := range mineCoordinates {
		(*b)[coordinate.row][coordinate.col] = 0
	}

	for coordinate := range mineCoordinates {
		b.UpdateValue(coordinate.row, coordinate.col-1)
		b.UpdateValue(coordinate.row, coordinate.col+1)

		b.UpdateValue(coordinate.row-1, coordinate.col)
		b.UpdateValue(coordinate.row-1, coordinate.col-1)
		b.UpdateValue(coordinate.row-1, coordinate.col+1)

		b.UpdateValue(coordinate.row+1, coordinate.col)
		b.UpdateValue(coordinate.row+1, coordinate.col-1)
		b.UpdateValue(coordinate.row+1, coordinate.col+1)
	}
}
