package parking

import (
	"errors"
	"fmt"
	"sync"
)

var (
	parkingLotInstance *ParkingLot
	parkingLotOnce     sync.Once
)

type ParkingLot struct {
	counter                int
	floors                 []ParkingFloor
	activeTickets          map[int]*ParkingTicket
	feeStrategy            FeeStrategy
	spotAllocationStrategy SpotAllocationStrategy
	mutex                  sync.Mutex
}

func (p *ParkingLot) Park(v Vehicle) (*ParkingTicket, error) {
	for {
		spot, err := p.spotAllocationStrategy.Allot(p.floors, v)
		if err != nil {
			return nil, err
		}

		if err := spot.ParkVehicle(v); err != nil {
			if errors.Is(err, ErrSpotOccupied) {
				// Another goroutine claimed this spot first; retry with a fresh allocation.
				continue
			}
			return nil, err
		}

		p.mutex.Lock()
		parkingTicket := NewParkingTicket(p.counter, v, spot)
		p.activeTickets[parkingTicket.id] = &parkingTicket
		p.counter++
		p.mutex.Unlock()

		return &parkingTicket, nil
	}
}

func (p *ParkingLot) UnPark(ticketId int) (*float64, error) {
	p.mutex.Lock()
	ticket, ok := p.activeTickets[ticketId]
	if !ok {
		p.mutex.Unlock()
		return nil, fmt.Errorf("Invalid ticket id")
	}
	delete(p.activeTickets, ticketId)
	p.mutex.Unlock()

	ticket.AddExitTime()
	parkingSpot := ticket.spot
	_, err := parkingSpot.UnParkVehicle()
	if err != nil {
		return nil, err
	}

	fee := p.feeStrategy.GetFee(ticket)
	return &fee, nil
}

func (p *ParkingLot) Availability() map[VehicleSize]int {
	available := map[VehicleSize]int{
		SMALL:  0,
		MEDIUM: 0,
		LARGE:  0,
	}
	for _, floor := range p.floors {
		available[SMALL] += floor.GetAvailableSpotsForSize(SMALL)
		available[MEDIUM] += floor.GetAvailableSpotsForSize(MEDIUM)
		available[LARGE] += floor.GetAvailableSpotsForSize(LARGE)
	}

	return available
}

func NewParkingLot(floors []ParkingFloor, feeStrategy FeeStrategy, spotStrategy SpotAllocationStrategy) ParkingLot {
	return ParkingLot{
		counter:                1,
		floors:                 floors,
		activeTickets:          map[int]*ParkingTicket{},
		feeStrategy:            feeStrategy,
		spotAllocationStrategy: spotStrategy,
		mutex:                  sync.Mutex{},
	}
}

func GetParkingLot(floors []ParkingFloor, feeStrategy FeeStrategy, spotStrategy SpotAllocationStrategy) *ParkingLot {
	parkingLotOnce.Do(func() {
		lot := NewParkingLot(floors, feeStrategy, spotStrategy)
		parkingLotInstance = &lot
	})
	return parkingLotInstance
}
