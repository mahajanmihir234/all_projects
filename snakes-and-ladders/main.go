package main

import (
	"fmt"
	"snakes-and-ladders/game"
)

func main() {
	snake1, err := game.NewSnake(99, 50)
	if err != nil {
		panic(err)
	}
	snake2, err := game.NewSnake(64, 23)
	if err != nil {
		panic(err)
	}
	snakes := []game.Snake{*snake1, *snake2}
	ladders := []game.Ladder{}
	playerNames := []string{"Mihir", "Somrat"}
	board, err := game.NewBoard(100, snakes, ladders)
	if err != nil {
		panic(err)
	}

	dice := game.MultipleDice{NumDice: 2}
	game := game.NewGame(
		*board, playerNames, dice,
	)
	game.Start()

	for {
		err := game.Play()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}
