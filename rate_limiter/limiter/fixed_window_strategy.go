package limiter

import (
	"time"
)

var FIXED_WINDOW_SIZE = 1 * time.Second

type FixedWindowStrategy struct {
	threshold  int64
	windowSize time.Duration
	userMap    map[Key]*Window
}

func (s FixedWindowStrategy) Allow(key Key) bool {
	window, ok := s.userMap[key]
	if !ok || time.Since(window.startTime) > s.windowSize {
		s.userMap[key] = NewWindow()
	}
	if window.counter.Load() < s.threshold {
		window.Increment()
		return true
	}
	return false
}

func NewFixedWindowStrategy(threshold int64, windowSize time.Duration) FixedWindowStrategy {
	return FixedWindowStrategy{
		threshold:  threshold,
		windowSize: windowSize,
		userMap:    map[Key]*Window{},
	}
}
