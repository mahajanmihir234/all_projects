package game

import (
	"errors"
	"fmt"
)

type Player struct {
	name     string
	position int
}

func (p Player) Name() string {
	return p.name
}

func (p Player) Position() int {
	return p.position
}

type Game struct {
	board              Board
	players            []Player
	currentPlayerIndex int
	dice               Dice
	status             GameStatus
}

func (g Game) CurrentPlayer() Player {
	return g.players[g.currentPlayerIndex]
}

type Roll struct {
	player        Player
	diceRoll      int
	finalPosition int
}

func (g *Game) Play() error {
	if g.status == NOT_STARTED {
		return errors.New("Game has not started")
	}
	if g.status == FINISHED {
		return errors.New("Game has finished")
	}

	for {
		diceRoll := g.dice.Roll()
		fmt.Println("Dice rolls a:", diceRoll)
		player := g.CurrentPlayer()
		finalPosition := g.board.GetFinalPosition(player.position + diceRoll)

		if finalPosition <= g.board.size {
			player.position = finalPosition
			g.players[g.currentPlayerIndex] = player
			fmt.Printf("Player %s is at position %d\n", player.name, player.position)

			if finalPosition == g.board.size {
				g.status = FINISHED
				return fmt.Errorf("Game won by player: %d", g.currentPlayerIndex)
			}
		}
		if diceRoll != g.dice.MaxValue() {
			g.MoveToNextPlayer()
			break
		}
	}

	return nil
}

func NewGame(board Board, playerNames []string, dice Dice) Game {
	players := []Player{}
	for _, name := range playerNames {
		players = append(players, Player{
			name:     name,
			position: 0,
		})
	}
	return Game{
		board:              board,
		players:            players,
		dice:               dice,
		currentPlayerIndex: 0,
		status:             NOT_STARTED,
	}
}

func (g *Game) Start() {
	g.status = IN_PROGRESS
}

func (g *Game) MoveToNextPlayer() {
	g.currentPlayerIndex = (g.currentPlayerIndex + 1) % len(g.players)
	fmt.Printf("Player changed to %s \n", g.CurrentPlayer().name)
}
