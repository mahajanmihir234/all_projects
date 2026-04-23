package game

import (
	"errors"
	"fmt"
)

type Board struct {
	size             int
	snakesAndLadders map[int]int
}

func (b Board) GetFinalPosition(position int) int {
	if value, ok := b.snakesAndLadders[position]; ok {
		if position < value {
			fmt.Printf("Ladder from %d to %d\n", position, value)
		} else {
			fmt.Printf("Snake from %d to %d\n", position, value)
		}
		return value
	}
	return position
}

func NewBoard(size int, snakes []Snake, ladders []Ladder) (*Board, error) {
	snakesAndLadders := map[int]int{}
	for _, snake := range snakes {
		if _, ok := snakesAndLadders[snake.start]; ok {
			return nil, errors.New("Cannot have two snakes with the same starting position")
		}
		snakesAndLadders[snake.start] = snake.end
	}

	for _, ladder := range ladders {
		if _, ok := snakesAndLadders[ladder.start]; ok {
			return nil, errors.New("Cannot have ladders with the same starting position")
		}
		snakesAndLadders[ladder.start] = ladder.end
	}

	return &Board{size: size, snakesAndLadders: snakesAndLadders}, nil
}
