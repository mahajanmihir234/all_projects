package system

import (
	"errors"
	"sort"
)

type LeaderBoard struct {
	players map[string]*Player
	users   map[string]*User
}

func (l *LeaderBoard) getOrCreatePlayers(playerIds []string) []*Player {
	players := []*Player{}
	for _, pId := range playerIds {
		player, ok := l.players[pId]
		if !ok {
			newPlayer := NewPlayer(pId)
			l.players[pId] = newPlayer
			players = append(players, newPlayer)
		} else {
			players = append(players, player)
		}
	}
	return players
}

func (l *LeaderBoard) AddUser(userId string, playerIds []string) error {
	if _, ok := l.users[userId]; ok {
		return errors.New("user already exists")
	}

	players := l.getOrCreatePlayers(playerIds)
	user := NewUser(userId, players)
	l.users[userId] = &user
	return nil
}

func (l *LeaderBoard) AddScore(playerId string, score int) error {
	player, ok := l.players[playerId]
	if !ok {
		return errors.New("player does not exist")
	}
	player.AddScore(score)
	return nil
}

func (l *LeaderBoard) GetTopK(k int) []*User {
	users := []*User{}
	for _, user := range l.users {
		users = append(users, user)
	}
	sort.Slice(users, func(i, j int) bool {
		if users[i].Score() == users[j].Score() {
			return users[i].id < users[j].id
		}
		return users[i].Score() > users[j].Score()
	})

	if len(users) <= k {
		return users
	}
	return users[:k]
}

func NewLeaderBoard() LeaderBoard {
	return LeaderBoard{players: map[string]*Player{}, users: map[string]*User{}}
}
