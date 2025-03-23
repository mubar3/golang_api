package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// CustomLogger struct untuk logger sederhana
type CustomLogger struct {
	file   *os.File
	fields map[string]interface{}
}

// Logger global yang dapat diakses di package ini
var Logger *CustomLogger

// NewLogger membuat instance baru dari CustomLogger
func NewLogger(logFilePath string) (*CustomLogger, error) {
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &CustomLogger{file: file}, nil
}

// WithField menambahkan field tambahan ke log
func (cl *CustomLogger) WithField(key string, value interface{}) *CustomLogger {
	newLogger := *cl // Salin instance logger
	if newLogger.fields == nil {
		newLogger.fields = make(map[string]interface{})
	}
	newLogger.fields[key] = value
	return &newLogger
}

// LogMessage menulis pesan log ke file
func (cl *CustomLogger) LogMessage(level string, message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	log.SetOutput(cl.file)

	// Tambahkan field ke log
	logFields := ""
	for k, v := range cl.fields {
		logFields += fmt.Sprintf("[%s: %v] ", k, v)
	}

	log.Printf("[%s] %s %s%s\n", level, timestamp, logFields, message)
}

// Close untuk menutup file log
func (cl *CustomLogger) Close() error {
	return cl.file.Close()
}

// InitLogger menginisialisasi logger global
func InitLogger(basePath string) error {
	// Dapatkan waktu sekarang
	now := time.Now()

	// Format nama folder dan file
	year := now.Format("2006")
	month := now.Format("01")
	date := now.Format("2006-01-02")
	yearPath := filepath.Join(basePath, year)
	monthPath := filepath.Join(yearPath, month)
	logFilePath := filepath.Join(monthPath, date+".log")

	// Buat folder jika belum ada
	if _, err := os.Stat(monthPath); os.IsNotExist(err) {
		err := os.MkdirAll(monthPath, 0755)
		if err != nil {
			return err
		}
	}

	// Buat logger baru
	var err error
	Logger, err = NewLogger(logFilePath)
	if err != nil {
		return err
	}

	return nil
}

// Logging dari fungsi main
// utils.Logger.LogMessage("INFO", "Aplikasi dimulai")
// utils.Logger.LogMessage("WARNING", "Peringatan contoh")
// utils.Logger.LogMessage("ERROR", "Pesan kesalahan contoh")
// utils.Logger.LogMessage("DEBUG", "Debugging aplikasi")
// utils.Logger.WithField("user_mobile", "1").LogMessage("LOG", "Debugging aplikasi")
