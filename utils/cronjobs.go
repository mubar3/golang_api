package utils

import (
	"golang_api/config"
	"log"

	"github.com/go-co-op/gocron"
)

func StartCronJobs() {
	// currentTime := time.Now().In(config.Timezone)
	scheduler := gocron.NewScheduler(config.Timezone)
	scheduler.Every(1).Hour().Do(func() {
		log.Println("Cron job executed")
	})

	scheduler.StartAsync()
}
