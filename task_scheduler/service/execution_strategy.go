package service

import "time"

type ExecutionStrategy interface {
	NextExecutionTimestamp(lastExecutionTimestamp *time.Time) *time.Time
}

type SingleExecutionStrategy struct {
	scheduledTime time.Time
}

func NewSingleExecutionStrategy(scheduledTime time.Time) SingleExecutionStrategy {
	return SingleExecutionStrategy{scheduledTime: scheduledTime}
}

func (s SingleExecutionStrategy) NextExecutionTimestamp(lastExecutionTimestamp *time.Time) *time.Time {
	if lastExecutionTimestamp != nil {
		return nil
	}
	return &s.scheduledTime
}

type RecurringExecutionStrategy struct {
	firstExecutionTimestamp time.Time
	interval                time.Duration
}

func NewRecurringExecutionStrategy(firstExecution time.Time, interval time.Duration) RecurringExecutionStrategy {
	return RecurringExecutionStrategy{
		firstExecutionTimestamp: firstExecution,
		interval:                interval,
	}
}

func (s RecurringExecutionStrategy) NextExecutionTimestamp(lastExecutionTimestamp *time.Time) *time.Time {
	if lastExecutionTimestamp == nil {
		return &s.firstExecutionTimestamp
	}
	nextTimestamp := lastExecutionTimestamp.Add(s.interval)
	return &nextTimestamp
}
