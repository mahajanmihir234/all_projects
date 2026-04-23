package parking

type FeeStrategy interface {
	GetFee(ticket ParkingTicket) float64
}

type FlatFeeStrategy struct{}

func (f FlatFeeStrategy) GetFee(ticket ParkingTicket) float64 {
	return 10.0
}

type HourlyParkingStrategy struct {
	HourlyRate int
}

func (h HourlyParkingStrategy) GetFee(ticket ParkingTicket) float64 {
	duration := ticket.GetDuration()
	return float64(h.HourlyRate) * duration.Hours()
}
