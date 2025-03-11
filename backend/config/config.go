// config/config.go
package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port            int
	StoreMasterFile string
}

func Load() Config {
	port, err := strconv.Atoi(getEnvOrDefault("PORT", "8080"))
	if err != nil {
		port = 8080
	}

	return Config{
		Port:            port,
		StoreMasterFile: getEnvOrDefault("STORE_MASTER_FILE", "./backend/store-master.json"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
