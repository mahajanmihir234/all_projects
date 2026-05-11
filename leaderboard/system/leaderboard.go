package system

import (
	"errors"
	"sort"
)

type LeaderBoard struct {
	players       map[string]*Player
	users         map[string]*User
	playerToUsers map[string]map[string]struct{}
	userScores    map[string]int
	ranking       []string
	userIndex     map[string]int
}

func (l *LeaderBoard) getOrCreatePlayer(playerId string) *Player {
	player, ok := l.players[playerId]
	if ok {
		return player
	}

	newPlayer := NewPlayer(playerId)
	l.players[playerId] = newPlayer
	return newPlayer
}

func (l *LeaderBoard) getOrCreatePlayers(playerIds []string) []*Player {
	players := []*Player{}
	for _, pId := range playerIds {
		players = append(players, l.getOrCreatePlayer(pId))
	}
	return players
}

func (l *LeaderBoard) betterUser(leftUserId string, leftScore int, rightUserId string, rightScore int) bool {
	if leftScore == rightScore {
		return leftUserId < rightUserId
	}
	return leftScore > rightScore
}

func (l *LeaderBoard) lowerBound(userId string, score int) int {
	return sort.Search(len(l.ranking), func(i int) bool {
		rankedUserId := l.ranking[i]
		rankedScore := l.userScores[rankedUserId]
		return !l.betterUser(rankedUserId, rankedScore, userId, score)
	})
}

func (l *LeaderBoard) insertIntoRanking(userId string) {
	score := l.userScores[userId]
	idx := l.lowerBound(userId, score)
	l.ranking = append(l.ranking, "")
	copy(l.ranking[idx+1:], l.ranking[idx:])
	l.ranking[idx] = userId
	for i := idx; i < len(l.ranking); i++ {
		l.userIndex[l.ranking[i]] = i
	}
}

func (l *LeaderBoard) removeFromRanking(userId string) {
	idx := l.userIndex[userId]
	l.ranking = append(l.ranking[:idx], l.ranking[idx+1:]...)
	delete(l.userIndex, userId)
	for i := idx; i < len(l.ranking); i++ {
		l.userIndex[l.ranking[i]] = i
	}
}

func (l *LeaderBoard) updateUserScore(userId string, delta int) {
	l.removeFromRanking(userId)
	l.userScores[userId] += delta
	l.insertIntoRanking(userId)
}

func (l *LeaderBoard) AddUser(userId string, playerIds []string) error {
	if _, ok := l.users[userId]; ok {
		return errors.New("user already exists")
	}

	seenPlayers := map[string]struct{}{}
	dedupedPlayerIds := make([]string, 0, len(playerIds))
	for _, playerId := range playerIds {
		if _, ok := seenPlayers[playerId]; ok {
			continue
		}
		seenPlayers[playerId] = struct{}{}
		dedupedPlayerIds = append(dedupedPlayerIds, playerId)
	}

	players := l.getOrCreatePlayers(dedupedPlayerIds)
	user := NewUser(userId, players)
	l.users[userId] = &user

	userScore := 0
	for _, player := range players {
		userScore += player.GetScore()
		if _, ok := l.playerToUsers[player.id]; !ok {
			l.playerToUsers[player.id] = map[string]struct{}{}
		}
		l.playerToUsers[player.id][userId] = struct{}{}
	}

	l.userScores[userId] = userScore
	l.insertIntoRanking(userId)
	return nil
}

func (l *LeaderBoard) AddScore(playerId string, score int) error {
	player := l.getOrCreatePlayer(playerId)
	player.AddScore(score)

	affectedUsers, ok := l.playerToUsers[playerId]
	if !ok {
		return nil
	}

	for userId := range affectedUsers {
		l.updateUserScore(userId, score)
	}
	return nil
}

func (l *LeaderBoard) GetTopK(k int) []*User {
	if k <= 0 {
		return []*User{}
	}

	limit := k
	if len(l.ranking) < limit {
		limit = len(l.ranking)
	}

	users := make([]*User, 0, limit)
	for i := 0; i < limit; i++ {
		userId := l.ranking[i]
		users = append(users, l.users[userId])
	}
	return users
}

func NewLeaderBoard() LeaderBoard {
	return LeaderBoard{
		players:       map[string]*Player{},
		users:         map[string]*User{},
		playerToUsers: map[string]map[string]struct{}{},
		userScores:    map[string]int{},
		ranking:       []string{},
		userIndex:     map[string]int{},
	}
}
