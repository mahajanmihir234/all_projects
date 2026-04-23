package parking

import (
	"sync"
	"testing"
)

func TestParkingLotConcurrentParkSameSpot(t *testing.T) {
	spots := NewParkingSpots(map[VehicleSize]int{
		SMALL:  0,
		MEDIUM: 0,
		LARGE:  1,
	})
	floor := NewParkingFloor(spots)
	lot := NewParkingLot(
		[]ParkingFloor{floor},
		HourlyParkingStrategy{HourlyRate: 10},
		FirstSpotStrategy{},
	)

	vehicles := []Vehicle{
		NewTruck("truck-1"),
		NewTruck("truck-2"),
	}

	start := make(chan struct{})
	var wg sync.WaitGroup
	var successCount int
	var failureCount int
	var resultMu sync.Mutex

	for _, v := range vehicles {
		wg.Add(1)
		go func(vehicle Vehicle) {
			defer wg.Done()
			<-start
			_, err := lot.Park(vehicle)

			resultMu.Lock()
			defer resultMu.Unlock()
			if err != nil {
				failureCount++
				return
			}
			successCount++
		}(v)
	}

	close(start)
	wg.Wait()

	if successCount != 1 {
		t.Fatalf("expected exactly 1 successful park, got %d", successCount)
	}
	if failureCount != 1 {
		t.Fatalf("expected exactly 1 failed park, got %d", failureCount)
	}

	availability := lot.Availability()
	if availability[LARGE] != 0 {
		t.Fatalf("expected no large spots left, got %d", availability[LARGE])
	}
}
