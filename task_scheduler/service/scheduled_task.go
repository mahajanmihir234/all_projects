package service

import (
	"fmt"
	"time"
)

type ScheduledTask struct {
	id                string
	executionStrategy ExecutionStrategy
	task              Task
	lastExecutionTime *time.Time
	sequenceNumber    int64
	status            TaskStatus
}

func (s *ScheduledTask) SetStatus(status TaskStatus) {
	s.status = status
}

func (s ScheduledTask) NextExecutionTime() *time.Time {
	return s.executionStrategy.NextExecutionTimestamp(s.lastExecutionTime)
}

func (this ScheduledTask) LessThan(other ScheduledTask) bool {
	thisExecutionTime := this.NextExecutionTime()
	otherExecutionTime := other.NextExecutionTime()
	if thisExecutionTime == nil && otherExecutionTime == nil {
		return false
	}
	if thisExecutionTime == nil {
		return false
	}
	if otherExecutionTime == nil {
		return true
	}
	if thisExecutionTime.Equal(*otherExecutionTime) {
		return this.sequenceNumber < other.sequenceNumber
	}
	return otherExecutionTime.After(*thisExecutionTime)
}

func (s ScheduledTask) HasMoreExecutions() bool {
	return s.NextExecutionTime() != nil
}

func (s ScheduledTask) Print() {
	fmt.Printf("Task: %s, nextExecutionTime: %v, status: %s\n", s.id, s.NextExecutionTime(), s.status)
}

func (s *ScheduledTask) UpdateForNextExecution() {
	now := time.Now()
	s.lastExecutionTime = &now
}
