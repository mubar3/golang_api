package config

import "time"

// Variabel global untuk menyimpan timezone
var Timezone *time.Location

// Fungsi untuk menginisialisasi timezone
func InitTimezone(location string) error {
	var err error
	Timezone, err = time.LoadLocation(location)
	return err
}
