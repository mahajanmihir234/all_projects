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

func (p ParkingSpot) Available() bool {
	return p.vehicle == nil
}
func (p ParkingSpot) CanFitVehicle(size VehicleSize) bool {
	return size <= p.size
}

func (p *ParkingSpot) ParkVehicle(v Vehicle) error {
	if !p.Available() {
		return errors.New("Already occupied")
	}

	if !p.CanFitVehicle(v.Size()) {
		return fmt.Errorf("Cannot park vehicle with size %d, parking spot has size %d", v.Size(), p.size)
	}

	p.mutex.Lock()
	p.vehicle = v
	p.mutex.Unlock()
	return nil
}

func (p *ParkingSpot) UnParkVehicle() (Vehicle, error) {
	if p.Available() {
		return nil, errors.New("no vehicle to unpark")
	}

	p.mutex.Lock()
	vehicle := p.vehicle
	p.vehicle = nil
	p.mutex.Unlock()
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
