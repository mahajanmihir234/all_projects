package main

import (
	"bufio"
	"fmt"
	"os"
	"tictactoe_v2/game"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	ticTacToe := game.NewGame(3)

	turn := game.X
	for {
		ticTacToe.Print()
		fmt.Printf("Turn: %s\n", string(turn))
		line, _, err := reader.ReadLine()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		var row, col int
		_, err = fmt.Sscanf(string(line), "%d %d", &row, &col)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		stopped, err := ticTacToe.Play(turn, row, col)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if stopped {
			fmt.Println("GAME WON!")
			return
		}
		turn = turn.Opposite()
	}
}
