package parking

import "errors"

type SpotAllocationStrategy interface {
	Allot(floors []ParkingFloor, vehicle Vehicle) (*ParkingSpot, error)
}

type FirstSpotStrategy struct{}

func (f FirstSpotStrategy) Allot(floors []ParkingFloor, v Vehicle) (*ParkingSpot, error) {
	var spot *ParkingSpot
	for _, floor := range floors {
		spot = floor.GetAvailableSpot(v.Size())
		if spot == nil {
			continue
		}
	}
	if spot == nil {
		return nil, errors.New("No parking spot found")
	}
	return spot, nil
}
