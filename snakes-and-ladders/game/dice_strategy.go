package game

import "math/rand/v2"

type Dice interface {
	Roll() int
	MaxValue() int
}

type FairDice struct{}

func (f FairDice) Roll() int {
	return rand.IntN(6) + 1
}
func (f FairDice) MaxValue() int {
	return 6
}

type MultipleDice struct {
	NumDice int
}

func (m MultipleDice) Roll() int {
	position := 0
	for range m.NumDice {
		position += rand.IntN(6) + 1
	}
	return position
}

func (m MultipleDice) MaxValue() int {
	return 6 * m.NumDice
}
