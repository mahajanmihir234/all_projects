package game

import "math/rand"

type MineGenerator interface {
	GenerateMines(rows int, cols int, mines int) map[Coordinate]bool
}

type RandomMineGenerator struct{}

func (r RandomMineGenerator) GenerateMines(rows int, cols int, numMines int) map[Coordinate]bool {
	mines := map[Coordinate]bool{}
	for numMines > 0 {
		coordinate := Coordinate{
			row: rand.Intn(rows),
			col: rand.Intn(cols),
		}
		if _, ok := mines[coordinate]; !ok {
			mines[coordinate] = true
			numMines -= 1
		}
	}
	return mines
}
