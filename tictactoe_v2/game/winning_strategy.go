package game

type WinningStrategy interface {
	Win(board Board, turn Symbol, row int, col int) bool
}

type RowWinningStrategy struct{}
type ColumnWinningStrategy struct{}
type Diagonal1WinningStrategy struct{}
type Diagonal2WinningStrategy struct{}

func (r RowWinningStrategy) Win(board Board, turn Symbol, row int, col int) bool {
	for i := range board.size {
		if board.grid[i][col] != &turn {
			return false
		}
	}

	return true
}

func (r ColumnWinningStrategy) Win(board Board, turn Symbol, row int, col int) bool {
	for i := range board.size {
		if board.grid[row][i] != &turn {
			return false
		}
	}

	return true
}

func (r Diagonal1WinningStrategy) Win(board Board, turn Symbol, row int, col int) bool {
	if row != col {
		return false
	}
	for i := range board.size {
		if board.grid[i][i] != &turn {
			return false
		}
	}

	return true
}

func (r Diagonal2WinningStrategy) Win(board Board, turn Symbol, row int, col int) bool {
	if row+col != board.size-1 {
		return false
	}
	for i := range board.size {
		if board.grid[i][board.size-1-i] != &turn {
			return false
		}
	}

	return true
}
