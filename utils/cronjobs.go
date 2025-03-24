package utils

import (
	"golang_api/config"
	"time"
)

func StartCronJobs() {
	scheduler := NewScheduler()

	scheduler.AddJob(func() { Isnotnull(nil) }, "0 0 * * 1") // Format cron (menit jam hari-bulan bulan hari-minggu)

	// Mulai scheduler
	go scheduler.Start()

	// Menjaga agar aplikasi tetap berjalan
	select {}
}

// CronJob adalah struktur untuk mendefinisikan tugas cron
type CronJob struct {
	Job      func() // Fungsi yang akan dijalankan
	Schedule string // Jadwal dalam format string
}

// Scheduler adalah struktur untuk mengelola cron job
type Scheduler struct {
	jobs []CronJob
}

// NewScheduler membuat instance baru dari Scheduler
func NewScheduler() *Scheduler {
	return &Scheduler{}
}

// AddJob menambahkan tugas baru ke scheduler
func (s *Scheduler) AddJob(job func(), schedule string) {
	s.jobs = append(s.jobs, CronJob{Job: job, Schedule: schedule})
}

// Start menjalankan scheduler
func (s *Scheduler) Start() {
	for {
		now := time.Now().In(config.Timezone)
		// Cek apakah sekarang adalah hari Senin dan pukul 00:00
		if now.Weekday() == time.Monday && now.Hour() == 0 && now.Minute() == 0 {
			for _, job := range s.jobs {
				go job.Job() // Jalankan tugas dalam goroutine
			}
			// Tunggu satu menit sebelum memeriksa lagi
			time.Sleep(1 * time.Minute)
		} else {
			// Tunggu satu detik sebelum memeriksa lagi
			time.Sleep(1 * time.Second)
		}
	}
}
