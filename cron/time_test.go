package cron

import (
	"testing"
	"time"
)

// TestElapsedTime checks the amount of time it has been since a PR was opned and returns true if the
// duration has exceeded limits
func TestElapsedTime(t *testing.T) {
	currentTime := time.Now()

	duration := 14400000000000
	count := (duration / 60000000000) - 1
	notElapsed := currentTime.Add(time.Duration(-count) * time.Minute)

	elapsed, message := ElapsedTime(notElapsed, duration)
	if elapsed {
		t.Errorf("got %v and duration is %v", elapsed, count)
	}
	if message != "Warning not detected" {
		t.Errorf("got %v", message)
	}

}
