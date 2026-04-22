package game

type Queue []Coordinate

func (q *Queue) Push(val Coordinate, rows int, cols int) {
	if val.row < 0 || val.col < 0 || val.row == rows || val.col == cols {
		return
	}
	*q = append(*q, val)
}

func (q *Queue) Pop() Coordinate {
	element := (*q)[0]
	(*q) = (*q)[1:]
	return element
}
