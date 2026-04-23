package game

import "fmt"

type BoardEntity interface {
	Move(start int) int
}

type Snake struct {
	start int
	end   int
}

func NewSnake(start int, end int) (*Snake, error) {
	if start <= end {
		return nil, fmt.Errorf("snake's start has to be bigger than end")
	}
	return &Snake{
		start: start,
		end:   end,
	}, nil
}

func (s Snake) Move(start int) int {
	if start == s.start {
		return s.end
	}
	return start
}

type Ladder struct {
	start int
	end   int
}

func (l Ladder) Move(start int) int {
	if start == l.start {
		return l.end
	}
	return start
}

func NewLadder(start int, end int) (*Ladder, error) {
	if start >= end {
		return nil, fmt.Errorf("ladder's start has to be lesser than end")
	}
	return &Ladder{
		start: start,
		end:   end,
	}, nil
}
