package service

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

// CronExecutionStrategy runs a task on a standard cron schedule (supports optional seconds field).
// anchor fixes the first fire time so heap ordering stays stable (Next must not call time.Now() on every peek).
type CronExecutionStrategy struct {
	schedule cron.Schedule
	anchor   time.Time
}

func NewCronExecutionStrategy(expression string, anchor time.Time) (CronExecutionStrategy, error) {
	parser := cron.NewParser(
		cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)
	schedule, err := parser.Parse(expression)
	if err != nil {
		return CronExecutionStrategy{}, fmt.Errorf("parse cron %q: %w", expression, err)
	}
	if anchor.IsZero() {
		anchor = time.Now()
	}
	return CronExecutionStrategy{schedule: schedule, anchor: anchor}, nil
}

func (s CronExecutionStrategy) NextExecutionTimestamp(lastExecutionTimestamp *time.Time) *time.Time {
	from := s.anchor
	if lastExecutionTimestamp != nil {
		from = *lastExecutionTimestamp
	}
	next := s.schedule.Next(from)
	return &next
}
