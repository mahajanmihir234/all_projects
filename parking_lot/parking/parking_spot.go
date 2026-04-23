package parking

import (
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type ParkingSpot struct {
	id      uuid.UUID
	size    VehicleSize
	vehicle Vehicle
	mutex   sync.RWMutex
}

var ErrSpotOccupied = errors.New("already occupied")

func (p *ParkingSpot) Available() bool {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	return p.vehicle == nil
}
func (p *ParkingSpot) CanFitVehicle(size VehicleSize) bool {
	return size <= p.size
}

func (p *ParkingSpot) ParkVehicle(v Vehicle) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.vehicle != nil {
		return ErrSpotOccupied
	}

	if !p.CanFitVehicle(v.Size()) {
		return fmt.Errorf("Cannot park vehicle with size %d, parking spot has size %d", v.Size(), p.size)
	}

	p.vehicle = v
	return nil
}

func (p *ParkingSpot) UnParkVehicle() (Vehicle, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.vehicle == nil {
		return nil, errors.New("no vehicle to unpark")
	}

	vehicle := p.vehicle
	p.vehicle = nil
	return vehicle, nil
}

func NewParkingSpots(mp map[VehicleSize]int) []*ParkingSpot {
	numSmallSlots := mp[SMALL]
	numMediumSlots := mp[MEDIUM]
	numLargeSlots := mp[LARGE]

	spots := []*ParkingSpot{}
	for range numSmallSlots {
		id := uuid.New()
		spots = append(spots, &ParkingSpot{
			id:      id,
			size:    SMALL,
			vehicle: nil,
			mutex:   sync.RWMutex{},
		})
	}

	for range numMediumSlots {
		id := uuid.New()
		spots = append(spots, &ParkingSpot{
			id:      id,
			size:    MEDIUM,
			vehicle: nil,
			mutex:   sync.RWMutex{},
		})
	}

	for range numLargeSlots {
		id := uuid.New()
		spots = append(spots, &ParkingSpot{
			id:      id,
			size:    LARGE,
			vehicle: nil,
			mutex:   sync.RWMutex{},
		})
	}
	return spots
}
