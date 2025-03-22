package utils

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger() {
	appPath, err := os.Executable() // Mendapatkan path aplikasi
	if err != nil {
		logrus.Fatal("Gagal mendapatkan path aplikasi:", err)
	}

	logPath := filepath.Join(filepath.Dir(appPath), "application.log")

	logFile := &lumberjack.Logger{
		Filename:   logPath, // Path dinamis untuk file log
		MaxSize:    10,      // Maksimal ukuran file dalam MB
		MaxBackups: 3,       // Jumlah backup
		MaxAge:     30,      // Maksimal umur file dalam hari
		Compress:   true,    // Kompres file log lama
	}

	logrus.SetOutput(logFile)
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.Info("Logger berhasil diinisialisasi")
}
