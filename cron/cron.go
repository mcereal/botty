package cron

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
)

// ScheduleCron triggers github api calls on a cadence
func ScheduleCron() {
	s := gocron.NewScheduler(time.UTC)

	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Warn("Invalid time zone")
	}
	s.ChangeLocation(location)
	fmt.Println("Timezone:", s.Location())
	// j, _ := s.Every(4).Hours().Do(github.GetOpenPrs)
	j, _ := s.Every(1).Monday().Tuesday().Wednesday().Thursday().Friday().At("09:30;13:30;17:30").Do(GetOpenPrs)

	s.StartAsync()

	fmt.Println("Running job:", j.IsRunning())

}
