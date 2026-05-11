package limiter

import (
	"sync/atomic"
	"time"
)

type Window struct {
	startTime time.Time
	counter   atomic.Int64
}

func (w *Window) Increment() {
	w.counter.Add(1)
}

func NewWindow() *Window {
	return &Window{startTime: time.Now(), counter: atomic.Int64{}}
}

func (w *Window) SetCounter(value int) {
	w.counter.Store(int64(value))
}
