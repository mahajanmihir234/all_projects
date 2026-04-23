package main

import (
	"fmt"
	"parking_lot/parking"
	"sync"
)

func main() {
	parkingSpots := parking.NewParkingSpots(map[parking.VehicleSize]int{
		parking.SMALL:  0,
		parking.MEDIUM: 0,
		parking.LARGE:  1,
	})

	parkingFloor := parking.NewParkingFloor(parkingSpots)
	feeStrategy := parking.HourlyParkingStrategy{HourlyRate: 20}
	spotAllocationStrategy := parking.FirstSpotStrategy{}

	parkingLot := parking.GetParkingLot([]parking.ParkingFloor{parkingFloor}, feeStrategy, spotAllocationStrategy)
	fmt.Println("Fresh lot")
	for size, availability := range parkingLot.Availability() {
		fmt.Println(parking.VehicleSize(size), availability)
	}

	// car := parking.NewCar("Baleno")
	truck := parking.NewTruck("Truck")
	truck2 := parking.NewTruck("Truck2")

	// bike := parking.NewBike("Bike")

	wg := sync.WaitGroup{}
	var truckTicket *parking.ParkingTicket
	var truck2Ticket *parking.ParkingTicket
	var err error

	wg.Add(1)
	go func() {
		defer wg.Done()
		truck2Ticket, err = parkingLot.Park(truck2)
		if err != nil {
			fmt.Println(err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		truckTicket, err = parkingLot.Park(truck)
		if err != nil {
			fmt.Println(err)
		}
	}()

	// carTicket, err := parkingLot.Park(car)
	// if err != nil {
	// 	panic(err)
	// }

	// bikeTicket, err := parkingLot.Park(bike)
	// if err != nil {
	// 	panic(err)
	// }

	wg.Wait()

	for size, availability := range parkingLot.Availability() {
		fmt.Println(parking.VehicleSize(size), availability)
	}

	if truckTicket != nil {
		truckFee, err := parkingLot.UnPark(truckTicket.Id())
		if err != nil {
			panic(err)
		}
		fmt.Println("TruckFee: ", *truckFee)
	}

	if truck2Ticket != nil {
		truckFee, err := parkingLot.UnPark(truck2Ticket.Id())
		if err != nil {
			panic(err)
		}
		fmt.Println("Truck2Fee: ", *truckFee)
	}

	// carFee, err := parkingLot.UnPark(carTicket.Id())
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("CarFee: ", *carFee)

	// bikeFee, err := parkingLot.UnPark(bikeTicket.Id())
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("bikeFee: ", *bikeFee)
}
