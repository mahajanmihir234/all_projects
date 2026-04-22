package main

import (
	"bufio"
	"fmt"
	"minesweeper/game"
	"os"
)

func main() {
	mineGenerator := game.RandomMineGenerator{}
	mineSweeper := game.NewMineSweeper(mineGenerator, 6, 6, 5)
	mineSweeper.Print(false)
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)

	for !mineSweeper.Win() {
		line, _, err := reader.ReadLine()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		var row, col int
		var moveType string
		_, err = fmt.Sscanf(string(line), "%s %d %d", &moveType, &row, &col)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		err = mineSweeper.Play(moveType, row, col)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println()
		mineSweeper.Print(false)
	}
	fmt.Println("GAME WON!")
}
