package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"task_scheduler/config"
	"task_scheduler/service"
)

func main() {
	scheduleDir := "schedules"
	if len(os.Args) > 1 {
		scheduleDir = os.Args[1]
	}
	absDir, err := filepath.Abs(scheduleDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "resolve schedule dir: %v\n", err)
		os.Exit(1)
	}

	baseTime := time.Now()
	tasks, err := config.LoadScheduleDir(absDir, baseTime)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load schedules: %v\n", err)
		os.Exit(1)
	}
	if len(tasks) == 0 {
		fmt.Fprintf(os.Stderr, "no tasks found in %s\n", absDir)
		os.Exit(1)
	}

	fmt.Println("=== Task scheduler demo (JSON schedules) ===")
	fmt.Printf("Loaded %d task(s) from %s\n\n", len(tasks), absDir)

	observer := service.LoggingObserver{}
	scheduler := service.NewTaskSchedulerService(2, []service.TaskObserver{observer})
	scheduler.Start()
	defer scheduler.Stop()

	go func() {
		for err := range scheduler.Errors() {
			fmt.Fprintf(os.Stderr, "task error: %v\n", err)
		}
	}()

	for _, loaded := range tasks {
		scheduler.Schedule(loaded.Task, loaded.Strategy)
		fmt.Printf("  scheduled %-20s  (%s)\n", loaded.Task.Id(), filepath.Base(loaded.Source))
	}

	fmt.Println("\nRunning for 8s...")
	time.Sleep(8 * time.Second)
	fmt.Println("=== Demo complete ===")
}
