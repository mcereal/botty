package cron

import (
	"fmt"
	"time"
)

// Set to 4 hours
const _duration time.Duration = 14400000000000

// ElapsedTime checks the amount of time it has been since a PR was opned and returns true if the
// duration has exceeded limits
func ElapsedTime(s time.Time) (bool, string) {
	elapsed := time.Since(s)

	fmt.Println("ELAPSED", elapsed)

	if elapsed > _duration {
		elapsedTime := elapsed.Round(time.Second)

		fmt.Println(elapsedTime)

		exceedLimit := fmt.Sprintf("Elapsed time of %v has exceeded the time limit of %v. Please convert to draft, or approve and merge", elapsedTime, _duration)
		fmt.Println(exceedLimit)
		return true, exceedLimit
	}
	return false, "Warning not detected"
}
