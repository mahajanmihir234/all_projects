package service

import (
	"container/heap"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type TaskSchedulerService struct {
	mutex           sync.Mutex
	taskQueue       ScheduledTaskQueue
	observers       []TaskObserver
	sequenceCounter atomic.Int64
	workers         []*Worker
	workerCount     int
	tasks           chan Task
	errCh           chan error
	wake            chan struct{}
	done            chan struct{}
	workerWg        sync.WaitGroup
	started         bool
}

func NewTaskSchedulerService(workerCount int, observers []TaskObserver) *TaskSchedulerService {
	if workerCount < 1 {
		workerCount = 1
	}

	taskQueue := ScheduledTaskQueue{}
	heap.Init(&taskQueue)

	return &TaskSchedulerService{
		taskQueue:       taskQueue,
		observers:       observers,
		sequenceCounter: atomic.Int64{},
		workerCount:     workerCount,
		tasks:           make(chan Task, workerCount),
		errCh:           make(chan error, workerCount),
		wake:            make(chan struct{}, 1),
		done:            make(chan struct{}),
	}
}

// Start launches worker goroutines and the scheduler loop that dispatches due tasks.
func (s *TaskSchedulerService) Start() {
	s.mutex.Lock()
	if s.started {
		s.mutex.Unlock()
		return
	}
	s.started = true
	s.mutex.Unlock()

	for i := 0; i < s.workerCount; i++ {
		worker := NewWorker(fmt.Sprint(i), s.tasks, &s.workerWg, s.errCh)
		s.workers = append(s.workers, &worker)
		worker.Start()
	}

	go s.runScheduler()
}

// Stop shuts down the scheduler and workers. It is safe to call only after Start.
func (s *TaskSchedulerService) Stop() {
	select {
	case <-s.done:
		return
	default:
		close(s.done)
	}

	close(s.tasks)
	s.workerWg.Wait()
}

// Errors returns a channel that receives task execution errors reported by workers.
func (s *TaskSchedulerService) Errors() <-chan error {
	return s.errCh
}

func (s *TaskSchedulerService) Schedule(task Task, executionStrategy ExecutionStrategy) string {
	sequenceNumber := s.sequenceCounter.Add(1)
	scheduledTask := ScheduledTask{
		id:                task.Id(),
		executionStrategy: executionStrategy,
		task:              task,
		lastExecutionTime: nil,
		sequenceNumber:    sequenceNumber,
		status:            SCHEDULED,
	}

	s.mutex.Lock()
	heap.Push(&s.taskQueue, scheduledTask)
	s.mutex.Unlock()

	s.signal()
	return task.Id()
}

func (s *TaskSchedulerService) signal() {
	select {
	case s.wake <- struct{}{}:
	default:
	}
}

func (s *TaskSchedulerService) runScheduler() {
	for {
		if s.dispatchDueTasks() {
			continue
		}

		wait := s.timeUntilNextTask()

		select {
		case <-s.done:
			return
		case <-s.wake:
		case <-time.After(wait):
		}
	}
}

// dispatchDueTasks pops and dispatches all tasks that are due now. Returns true if any were dispatched.
func (s *TaskSchedulerService) dispatchDueTasks() bool {
	dispatched := false
	now := time.Now()

	for {
		s.mutex.Lock()
		if s.taskQueue.Len() == 0 {
			s.mutex.Unlock()
			return dispatched
		}

		next := s.taskQueue[0]
		nextTime := next.NextExecutionTime()
		if nextTime == nil {
			heap.Pop(&s.taskQueue)
			s.mutex.Unlock()
			continue
		}
		if now.Before(*nextTime) {
			s.mutex.Unlock()
			return dispatched
		}

		scheduled := heap.Pop(&s.taskQueue).(ScheduledTask)
		s.mutex.Unlock()

		scheduled.SetStatus(RUNNING)
		s.tasks <- workerTask{service: s, scheduled: scheduled}
		dispatched = true
	}
}

// timeUntilNextTask returns how long to wait before the next task may be due, or a long duration if the queue is empty.
func (s *TaskSchedulerService) timeUntilNextTask() time.Duration {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.taskQueue.Len() == 0 {
		return 24 * time.Hour
	}

	nextTime := s.taskQueue[0].NextExecutionTime()
	if nextTime == nil {
		return 0
	}

	if d := time.Until(*nextTime); d > 0 {
		return d
	}
	return 0
}

func (s *TaskSchedulerService) requeue(scheduled ScheduledTask) {
	s.mutex.Lock()
	heap.Push(&s.taskQueue, scheduled)
	s.mutex.Unlock()
	s.signal()
}

func (s *TaskSchedulerService) notifyStarted(task Task) {
	for _, o := range s.observers {
		o.OnTaskStarted(task)
	}
}

func (s *TaskSchedulerService) notifyCompleted(task Task) {
	for _, o := range s.observers {
		o.OnTaskCompleted(task)
	}
}

func (s *TaskSchedulerService) notifyFailed(task Task, err error) {
	for _, o := range s.observers {
		o.OnTaskFailed(task, err)
	}
}

// workerTask wraps a scheduled task so workers can execute it with observer and re-queue logic.
type workerTask struct {
	service   *TaskSchedulerService
	scheduled ScheduledTask
}

func (w workerTask) Id() string {
	return w.scheduled.task.Id()
}

func (w workerTask) Execute() error {
	task := w.scheduled.task
	w.service.notifyStarted(task)

	err := task.Execute()
	if err != nil {
		w.scheduled.SetStatus(FAILED)
		w.service.notifyFailed(task, err)
		return err
	}

	w.scheduled.UpdateForNextExecution()
	if w.scheduled.HasMoreExecutions() {
		w.scheduled.SetStatus(SCHEDULED)
		w.service.requeue(w.scheduled)
	} else {
		w.scheduled.SetStatus(FINISHED)
	}
	w.service.notifyCompleted(task)
	return nil
}
