package counter

type ClickCounter struct {
	capacity     int
	timestampMap map[int]int
}

func (c ClickCounter) RecordClick(timestamp int) {
	c.timestampMap[timestamp]++
}

func (c ClickCounter) GetRecentClicks(timestamp int) int {
	counter := 0
	for t, count := range c.timestampMap {
		if t <= timestamp && t > max(timestamp-c.capacity, 0) {
			counter += count
		}
	}
	return counter
}

func NewClickCounter(capacity int) ClickCounter {
	return ClickCounter{
		capacity:     capacity,
		timestampMap: map[int]int{},
	}
}
