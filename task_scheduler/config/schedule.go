package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"task_scheduler/service"
)

// ScheduleFile is the top-level document in each JSON file under schedules/.
type ScheduleFile struct {
	Version int          `json:"version"`
	Tasks   []TaskConfig `json:"tasks"`
}

// TaskConfig describes one schedulable task and how it should run.
type TaskConfig struct {
	ID       string         `json:"id"`
	TaskType string         `json:"task_type,omitempty"` // print (default) or fail
	Schedule ScheduleConfig `json:"schedule"`
}

// ScheduleConfig supports once, fixed interval, and cron expression strategies.
type ScheduleConfig struct {
	Type       string `json:"type"` // once | interval | cron
	In         string `json:"in,omitempty"`
	Every      string `json:"every,omitempty"`
	StartIn    string `json:"start_in,omitempty"`
	At         string `json:"at,omitempty"`         // RFC3339
	Expression string `json:"expression,omitempty"` // cron, e.g. "*/2 * * * *"
}

// LoadedTask pairs a service task with its execution strategy.
type LoadedTask struct {
	Task     service.Task
	Strategy service.ExecutionStrategy
	Source   string
}

// LoadScheduleDir reads every *.json file in dir and returns tasks to schedule.
func LoadScheduleDir(dir string, baseTime time.Time) ([]LoadedTask, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read schedule dir %q: %w", dir, err)
	}

	var loaded []LoadedTask
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		tasks, err := LoadScheduleFile(path, baseTime)
		if err != nil {
			return nil, err
		}
		loaded = append(loaded, tasks...)
	}
	return loaded, nil
}

// LoadScheduleFile parses a single schedule JSON file.
func LoadScheduleFile(path string, baseTime time.Time) ([]LoadedTask, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %q: %w", path, err)
	}

	var file ScheduleFile
	if err := json.Unmarshal(data, &file); err != nil {
		return nil, fmt.Errorf("parse %q: %w", path, err)
	}
	if len(file.Tasks) == 0 {
		return nil, nil
	}

	loaded := make([]LoadedTask, 0, len(file.Tasks))
	for _, tc := range file.Tasks {
		task, strategy, err := tc.build(baseTime)
		if err != nil {
			return nil, fmt.Errorf("%s: task %q: %w", path, tc.ID, err)
		}
		loaded = append(loaded, LoadedTask{
			Task:     task,
			Strategy: strategy,
			Source:   path,
		})
	}
	return loaded, nil
}

func (tc TaskConfig) build(baseTime time.Time) (service.Task, service.ExecutionStrategy, error) {
	if tc.ID == "" {
		return nil, nil, fmt.Errorf("id is required")
	}

	task, err := tc.buildTask()
	if err != nil {
		return nil, nil, err
	}

	strategy, err := tc.Schedule.buildStrategy(baseTime)
	if err != nil {
		return nil, nil, err
	}
	return task, strategy, nil
}

func (tc TaskConfig) buildTask() (service.Task, error) {
	switch strings.ToLower(tc.TaskType) {
	case "", "print":
		return service.NewPrintTask(tc.ID), nil
	case "fail":
		return failTask{id: tc.ID}, nil
	default:
		return nil, fmt.Errorf("unknown task_type %q", tc.TaskType)
	}
}

func (sc ScheduleConfig) buildStrategy(baseTime time.Time) (service.ExecutionStrategy, error) {
	switch strings.ToLower(sc.Type) {
	case "once":
		return sc.buildOnce(baseTime)
	case "interval":
		return sc.buildInterval(baseTime)
	case "cron":
		return sc.buildCron(baseTime)
	default:
		return nil, fmt.Errorf("unknown schedule type %q", sc.Type)
	}
}

func (sc ScheduleConfig) buildOnce(baseTime time.Time) (service.ExecutionStrategy, error) {
	switch {
	case sc.In != "":
		d, err := time.ParseDuration(sc.In)
		if err != nil {
			return nil, fmt.Errorf("in: %w", err)
		}
		return service.NewSingleExecutionStrategy(baseTime.Add(d)), nil
	case sc.At != "":
		t, err := time.Parse(time.RFC3339, sc.At)
		if err != nil {
			return nil, fmt.Errorf("at: %w", err)
		}
		return service.NewSingleExecutionStrategy(t), nil
	default:
		return nil, fmt.Errorf("once schedule requires \"in\" or \"at\"")
	}
}

func (sc ScheduleConfig) buildInterval(baseTime time.Time) (service.ExecutionStrategy, error) {
	if sc.Every == "" {
		return nil, fmt.Errorf("interval schedule requires \"every\"")
	}
	every, err := time.ParseDuration(sc.Every)
	if err != nil {
		return nil, fmt.Errorf("every: %w", err)
	}
	if every <= 0 {
		return nil, fmt.Errorf("every must be positive")
	}

	first := baseTime
	if sc.StartIn != "" {
		startIn, err := time.ParseDuration(sc.StartIn)
		if err != nil {
			return nil, fmt.Errorf("start_in: %w", err)
		}
		first = baseTime.Add(startIn)
	}
	return service.NewRecurringExecutionStrategy(first, every), nil
}

func (sc ScheduleConfig) buildCron(anchor time.Time) (service.ExecutionStrategy, error) {
	if sc.Expression == "" {
		return nil, fmt.Errorf("cron schedule requires \"expression\"")
	}
	return service.NewCronExecutionStrategy(sc.Expression, anchor)
}

// failTask is used when task_type is "fail" in JSON schedules.
type failTask struct {
	id string
}

func (f failTask) Id() string { return f.id }

func (f failTask) Execute() error {
	return fmt.Errorf("simulated failure for %s", f.id)
}
