package cron

import (
	"fmt"
	"time"
)

// Set to 4 hours
// const _duration time.Duration = 14400000000000

// ElapsedTime checks the amount of time it has been since a PR was opned and returns true if the
// duration has exceeded limits
func ElapsedTime(s time.Time, d int) (bool, string) {
	elapsed := time.Since(s)

	_duration := time.Duration(d)
	if elapsed > _duration {
		elapsedTime := elapsed.Round(time.Second)
		exceedLimit := fmt.Sprintf("Elapsed time of %v has exceeded the time limit of %v. Please convert to draft, or approve and merge", elapsedTime, _duration)
		return true, exceedLimit
	}
	return false, "Warning not detected"
}
