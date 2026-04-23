package parking

import "fmt"

type ParkingLot struct {
	counter                int
	floors                 []ParkingFloor
	activeTickets          map[int]ParkingTicket
	feeStrategy            FeeStrategy
	spotAllocationStrategy SpotAllocationStrategy
}

func (p *ParkingLot) Park(v Vehicle) (*ParkingTicket, error) {
	spot, err := p.spotAllocationStrategy.Allot(p.floors, v)
	if err != nil {
		return nil, err
	}
	fmt.Println("BeforeParking: ", spot.id, spot.size, spot.vehicle)

	parkingTicket := NewParkingTicket(p.counter, v, spot)
	err = spot.ParkVehicle(v)
	if err != nil {
		return nil, err
	}
	fmt.Println("AfterParking: ", spot.id, spot.size, spot.vehicle)
	p.activeTickets[parkingTicket.id] = parkingTicket
	p.counter++

	return &parkingTicket, nil
}

func (p ParkingLot) UnPark(ticketId int) (*float64, error) {
	if _, ok := p.activeTickets[ticketId]; !ok {
		return nil, fmt.Errorf("Invalid ticket id")
	}
	ticket := p.activeTickets[ticketId]
	delete(p.activeTickets, ticketId)

	ticket.AddExitTime()
	parkingSpot := ticket.spot
	_, err := parkingSpot.UnParkVehicle()
	if err != nil {
		return nil, err
	}

	fee := p.feeStrategy.GetFee(ticket)
	return &fee, nil
}

func (p ParkingLot) Availability() map[VehicleSize]int {
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
		activeTickets:          map[int]ParkingTicket{},
		feeStrategy:            feeStrategy,
		spotAllocationStrategy: spotStrategy,
	}
}
