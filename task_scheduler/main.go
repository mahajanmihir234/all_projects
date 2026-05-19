package main

import (
	"fmt"
	"os"
	"time"

	"task_scheduler/service"
)

// failTask demonstrates observer and error-channel behavior on failure.
type failTask struct {
	id string
}

func (f failTask) Id() string { return f.id }

func (f failTask) Execute() error {
	return fmt.Errorf("simulated failure for %s", f.id)
}

func main() {
	observer := service.LoggingObserver{}
	scheduler := service.NewTaskSchedulerService(2, []service.TaskObserver{observer})
	scheduler.Start()
	defer scheduler.Stop()

	go func() {
		for err := range scheduler.Errors() {
			fmt.Fprintf(os.Stderr, "task error: %v\n", err)
		}
	}()

	now := time.Now()
	fmt.Println("=== Task scheduler demo ===")
	fmt.Println("Scheduling one-shot, recurring, and failing tasks...")

	// One-shot tasks at different times (heap orders by soonest first).
	scheduler.Schedule(
		service.NewPrintTask("welcome-email"),
		service.NewSingleExecutionStrategy(now.Add(1*time.Second)),
	)
	scheduler.Schedule(
		service.NewPrintTask("daily-report"),
		service.NewSingleExecutionStrategy(now.Add(2*time.Second)),
	)
	scheduler.Schedule(
		service.NewPrintTask("nightly-backup"),
		service.NewSingleExecutionStrategy(now.Add(3*time.Second)),
	)

	// Recurring task: first run in 500ms, then every 500ms until shutdown.
	scheduler.Schedule(
		service.NewPrintTask("health-check"),
		service.NewRecurringExecutionStrategy(now.Add(500*time.Millisecond), 500*time.Millisecond),
	)

	scheduler.Schedule(
		failTask{id: "bad-import"},
		service.NewSingleExecutionStrategy(now.Add(1500*time.Millisecond)),
	)

	// Let all scheduled work finish before shutdown.
	time.Sleep(5 * time.Second)
	fmt.Println("=== Demo complete ===")
}
