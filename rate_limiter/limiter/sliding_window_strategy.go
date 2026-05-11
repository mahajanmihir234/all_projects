package limiter

import (
	"time"
)

type SlidingWindowStrategy struct {
	threshold  int64
	timestamps []time.Time
	windowSize time.Duration
	userMap    map[Key]*Window
}

func (s *SlidingWindowStrategy) setTimestamps() {
	now := time.Now()
	newTimestamps := []time.Time{}
	for i := len(s.timestamps) - 1; i >= 0; i-- {
		if now.Sub(s.timestamps[i]) <= s.windowSize {
			newTimestamps = append(newTimestamps, s.timestamps[i])
		}
		if len(newTimestamps) == int(s.threshold) {
			break
		}
	}
	s.timestamps = newTimestamps
}

func (s *SlidingWindowStrategy) Allow(key Key) bool {
	s.setTimestamps()
	if _, ok := s.userMap[key]; !ok {
		s.userMap[key] = NewWindow()
	}
	window := s.userMap[key]
	window.SetCounter(len(s.timestamps))

	if window.counter.Load() < s.threshold {
		window.Increment()
		return true
	}
	return false
}
