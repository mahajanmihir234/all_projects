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
}

func (t ParkingTicket) Id() int {
	return t.id
}

func (p ParkingTicket) AddExitTime() {
	once := sync.Once{}
	once.Do(func() {
		exitTime := time.Now()
		p.exitTime = &exitTime
	})
}

func (p ParkingTicket) GetDuration() time.Duration {
	exitTime := time.Now()
	if p.exitTime != nil {
		exitTime = *p.exitTime
	}

	return exitTime.Sub(p.entryTime)
}

func NewParkingTicket(id int, v Vehicle, spot *ParkingSpot) ParkingTicket {
	return ParkingTicket{
		id:        id,
		vehicle:   v,
		spot:      spot,
		entryTime: time.Now(),
		exitTime:  nil,
	}
}
