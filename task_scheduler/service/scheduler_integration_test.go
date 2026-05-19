package service

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestScheduler_DispatchesDueTask(t *testing.T) {
	scheduler := NewTaskSchedulerService(2, nil)
	scheduler.Start()
	defer scheduler.Stop()

	done := make(chan struct{})
	var runs atomic.Int32
	task := runOnceTask{
		id: "test",
		run: func() error {
			runs.Add(1)
			close(done)
			return nil
		},
	}

	scheduler.Schedule(task, NewSingleExecutionStrategy(time.Now().Add(50*time.Millisecond)))

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatalf("task did not run, runs=%d", runs.Load())
	}
}

type runOnceTask struct {
	id  string
	run func() error
}

func (r runOnceTask) Id() string { return r.id }

func (r runOnceTask) Execute() error { return r.run() }
