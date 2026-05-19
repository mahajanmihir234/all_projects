package service

import (
	"fmt"
	"log"
)

type TaskObserver interface {
	Id() string
	OnTaskStarted(task Task)
	OnTaskCompleted(task Task)
	OnTaskFailed(task Task, err error)
}

type LoggingObserver struct {
	id     string
	logger *log.Logger
}

func (l LoggingObserver) Id() string {
	return l.id
}

func (l LoggingObserver) OnTaskStarted(task Task) {
	fmt.Printf("Task: %s started\n", task.Id())
}

func (l LoggingObserver) OnTaskCompleted(task Task) {
	fmt.Printf("Task: %s completed\n", task.Id())
}

func (l LoggingObserver) OnTaskFailed(task Task, err error) {
	fmt.Printf("Task: %s failed with error: %s\n", task.Id(), err.Error())
}
