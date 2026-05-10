package main

import (
	"fmt"
	"leaderboard/system"
)

func main() {
	leaderBoard := system.NewLeaderBoard()
	leaderBoard.AddUser("u1", []string{"p1", "p2", "p3"})
	leaderBoard.AddUser("u2", []string{"p2", "p3"})

	leaderBoard.AddScore("p2", -2)
	leaderBoard.AddScore("p2", 3)

	answer := leaderBoard.GetTopK(2)
	for _, user := range answer {
		fmt.Println(user.Id(), user.Score())
	}

	leaderBoard.AddScore("p2", 7)

	answer = leaderBoard.GetTopK(2)
	for _, user := range answer {
		fmt.Println(user.Id(), user.Score())
	}
}
