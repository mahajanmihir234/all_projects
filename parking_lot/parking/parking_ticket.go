package parking

import (
	"sync"
	"time"
)

type ParkingTicket struct {
	id        int
	vehicle   Vehicle
	spot      *ParkingSpot
	entryTime time.Time
	exitTime  *time.Time
	once      sync.Once
	mutex     sync.RWMutex
}

func (t *ParkingTicket) Id() int {
	return t.id
}

func (p *ParkingTicket) AddExitTime() {
	p.once.Do(func() {
		exitTime := time.Now()
		p.mutex.Lock()
		p.exitTime = &exitTime
		p.mutex.Unlock()
	})
}

func (p *ParkingTicket) GetDuration() time.Duration {
	p.mutex.RLock()
	exitTime := time.Now()
	if p.exitTime != nil {
		exitTime = *p.exitTime
	}
	p.mutex.RUnlock()

	return exitTime.Sub(p.entryTime)
}

func NewParkingTicket(id int, v Vehicle, spot *ParkingSpot) ParkingTicket {
	return ParkingTicket{
		id:        id,
		vehicle:   v,
		spot:      spot,
		entryTime: time.Now(),
		exitTime:  nil,
		once:      sync.Once{},
		mutex:     sync.RWMutex{},
	}
}
