package config

import (
	"path/filepath"
	"testing"
	"time"
)

func TestLoadScheduleDir_NextTimes(t *testing.T) {
	dir := filepath.Join("..", "schedules")
	base := time.Now()
	loaded, err := LoadScheduleDir(dir, base)
	if err != nil {
		t.Fatal(err)
	}
	for _, lt := range loaded {
		next := lt.Strategy.NextExecutionTimestamp(nil)
		if next == nil {
			t.Errorf("task %s: nil next time", lt.Task.Id())
			continue
		}
		t.Logf("%s next=%v (in %v)", lt.Task.Id(), next.Format(time.RFC3339), time.Until(*next))
	}
}
