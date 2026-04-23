package parking

type ParkingFloor struct {
	spots []*ParkingSpot
}

func (p ParkingFloor) Available(size VehicleSize) bool {
	for _, spot := range p.spots {
		if spot.Available() && spot.CanFitVehicle(size) {
			return true
		}
	}
	return false
}

func (p ParkingFloor) GetAvailableSpot(size VehicleSize) *ParkingSpot {
	for _, spot := range p.spots {
		if spot.Available() && spot.CanFitVehicle(size) {
			return spot
		}
	}
	return nil
}

func (p ParkingFloor) AvailableSpots(size VehicleSize) int {
	counter := 0
	for _, spot := range p.spots {
		if spot.Available() && spot.CanFitVehicle(size) {
			counter++
		}
	}
	return counter
}
func (p ParkingFloor) GetAvailableSpotsForSize(size VehicleSize) int {
	counter := 0
	for _, spot := range p.spots {
		if spot.Available() && spot.size == size {
			counter++
		}
	}
	return counter
}

func NewParkingFloor(spots []*ParkingSpot) ParkingFloor {
	return ParkingFloor{spots: spots}
}
