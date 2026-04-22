package main

import (
	"bufio"
	"fmt"
	"os"
	"tictactoe/app"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	ticTacToe := app.NewTicTacToe(4)

	turn := app.X
	for !ticTacToe.Win() && !ticTacToe.Board.Filled() {
		ticTacToe.Board.Print()
		fmt.Printf("Turn: %s\n", string(ticTacToe.Turn))
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

		err = ticTacToe.Play(row, col, app.Value(turn))
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		turn = turn.Opposite()
		fmt.Println()
	}
	ticTacToe.Board.Print()

	if !ticTacToe.Win() && ticTacToe.Board.Filled() {
		fmt.Println("Draw")
	} else {
		fmt.Printf("Game won by %s\n", string(ticTacToe.Turn.Opposite()))
	}
}
