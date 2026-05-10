package system

import "sync"

type Player struct {
	id    string
	mutex sync.Mutex
	score int
}

func (p *Player) AddScore(score int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.score += score
}

func (p *Player) GetScore() int {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.score
}

func NewPlayer(id string) *Player {
	return &Player{
		id:    id,
		mutex: sync.Mutex{},
		score: 0,
	}
}
