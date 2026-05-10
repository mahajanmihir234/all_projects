package system

type User struct {
	id      string
	players []*Player
}

func (u *User) Score() int {
	score := 0
	for _, p := range u.players {
		score += p.GetScore()
	}
	return score
}

func NewUser(userId string, players []*Player) User {
	return User{id: userId, players: players}
}

func (u *User) Id() string {
	return u.id
}
