package utils

import (
	"bufio"
	"os"
	"strings"
)

func LoadEnv(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Mengabaikan komentar dan baris kosong
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Memisahkan key dan value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue // Jika tidak ada key atau value, lewati
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Menghapus tanda kutip jika ada
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = value[1 : len(value)-1] // Menghapus tanda kutip
		}

		// Mengatur variabel lingkungan
		os.Setenv(key, value)
	}

	return scanner.Err()
}
