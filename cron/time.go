package cron

import (
	"fmt"
	"time"
)

// ElapsedTime checks the amount of time it has been since a PR was opned and returns true if the
// duration has exceeded limits
func ElapsedTime(s time.Time, d int) (bool, string) {
	elapsed := time.Since(s)

	duration := time.Duration(d)
	if elapsed > duration {
		elapsedTime := elapsed.Round(time.Second)
		exceedLimit := fmt.Sprintf("Elapsed time of %v has exceeded the time limit of %v. Please convert to draft, or approve and merge", elapsedTime, duration)
		return true, exceedLimit
	}
	return false, "Warning not detected"
}
