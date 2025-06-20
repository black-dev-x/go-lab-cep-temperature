package config

import (
	"bufio"
	"os"
	"strings"
)

var WEATHER_API_KEY string

func Load() {
	if file, err := os.Open(".env"); err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			os.Setenv(key, value)
		}
	}
	WEATHER_API_KEY = os.Getenv("WEATHER_API_KEY")
}
