package database

import (
	"github.com/Rhaqim/savannahtech/internal/config"
	"github.com/Rhaqim/savannahtech/internal/log"

	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client // Global Redis connection

// Init initializes the Redis connection using the configuration values from the config package.
// It establishes a connection to the Redis database and assigns the connection to the global Redis variable.
// If an error occurs during the connection process, it logs the error and shuts down the logger.
func CacheInit() {
	log.InfoLogger.Println("Connecting to Redis database...")

	Redis = redis.NewClient(&redis.Options{
		Addr:     config.RedisAddress,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})
}

func CacheClose() {
	err := Redis.Close()
	if err != nil {
		log.ErrorLogger.Fatal("Failed to close Redis connection:", err)
	}
}
