package main

import (
	"fmt"
	"hit_counter/counter"
)

func main() {
	capacity := 2
	tracker := counter.NewClickCounter(capacity)

	tracker.RecordClick(1)
	tracker.RecordClick(1)
	tracker.RecordClick(2)

	answer := tracker.GetRecentClicks(3)
	fmt.Println("Clicks at 3:", answer)
}
