package cron

import (
	"testing"
	"time"

	"github.com/mcereal/botty/config"
)

// TestElapsedTime checks the amount of time it has been since a PR was opned and returns true if the
// duration has exceeded limits
func TestElapsedTime(t *testing.T) {
	currentTime := time.Now()

	for _, v := range config.AppConfig.Team {
		count := (v.CronElapsedDuration / 60000000000) - 1
		notElapsed := currentTime.Add(time.Duration(-count) * time.Minute)

		elapsed, message := ElapsedTime(notElapsed, v.CronElapsedDuration)
		if elapsed != true {
			t.Errorf("got %v", elapsed)
		}
		if message != "Warning not detected" {
			t.Errorf("got %v", message)
		}
	}
}
