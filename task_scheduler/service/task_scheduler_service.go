package service

import (
	"container/heap"
	"fmt"
	"sync"
	"sync/atomic"
)

type TaskSchedulerService struct {
	mutex           sync.Mutex
	taskQueue       ScheduledTaskQueue
	observers       []TaskObserver
	sequenceCounter atomic.Int64
	workers         []*Worker
	err             chan<- error
}

func NewTaskSchedulerService(workerCount int, observers []TaskObserver) TaskSchedulerService {
	err := make(chan<- error)
	taskChannel := make(<-chan Task)
	waitGroup := sync.WaitGroup{}

	workers := []*Worker{}
	for i := range workerCount {
		worker := NewWorker(fmt.Sprint(i), taskChannel, &waitGroup, err)
		workers = append(workers, &worker)
		worker.Start()
	}

	taskQueue := ScheduledTaskQueue{}
	heap.Init(&taskQueue)

	return TaskSchedulerService{
		mutex:           sync.Mutex{},
		taskQueue:       taskQueue,
		observers:       observers,
		sequenceCounter: atomic.Int64{},
		workers:         workers,
		err:             err,
	}
}

func (service *TaskSchedulerService) Schedule(task Task, executionStrategy ExecutionStrategy) string {
	service.sequenceCounter.Add(1)
	sequenceNumber := service.sequenceCounter.Load()
	scheduledTask := ScheduledTask{
		id:                task.Id(),
		executionStrategy: executionStrategy,
		task:              task,
		lastExecutionTime: nil,
		sequenceNumber:    sequenceNumber,
		status:            SCHEDULED,
	}

	service.mutex.Lock()
	defer service.mutex.Unlock()
	heap.Push(&service.taskQueue, scheduledTask)

	return task.Id()
}
