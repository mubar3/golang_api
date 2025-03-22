package utils

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"
)

func StartCronJobs() {
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(1).Hour().Do(func() {
		log.Println("Cron job executed")
	})

	scheduler.StartAsync()
}
