package service

import "fmt"

type TaskStatus string

const (
	SCHEDULED TaskStatus = "SCHEDULED"
	RUNNING   TaskStatus = "RUNNING"
	FINISHED  TaskStatus = "FINISHED"
	CANCELLED TaskStatus = "CANCELLED"
	FAILED    TaskStatus = "FAILED"
)

type Task interface {
	Id() string
	Execute() error
}

type PrintTask struct {
	id string
}

func (p PrintTask) Id() string {
	return p.id
}

func (p PrintTask) Execute() error {
	fmt.Println(p.id)
	return nil
}
