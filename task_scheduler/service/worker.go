package service

import "sync"

type Worker struct {
	id    string
	tasks <-chan Task
	err   chan<- error
	quit  chan struct{}
	wg    *sync.WaitGroup
}

func NewWorker(id string, tasks <-chan Task, wg *sync.WaitGroup, err chan<- error) Worker {
	return Worker{
		id:    id,
		tasks: tasks,
		err:   err,
		wg:    wg,
		quit:  make(chan struct{}),
	}
}

func (worker *Worker) Start() {
	worker.wg.Add(1)
	go func(w *Worker) {
		defer w.wg.Done()
		for {
			select {
			case task, ok := <-w.tasks:
				if !ok {
					return
				}
				if err := task.Execute(); err != nil {
					w.err <- err
				}
			case <-w.quit:
				return
			}
		}
	}(worker)
}
