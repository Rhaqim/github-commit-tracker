package database

import (
	"savannahtech/internal/config"
	"savannahtech/internal/log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB // Global database connection

func Init() {
	var err error

	log.InfoLogger.Println("Connecting to PostgreSQL database...")

	var dsn string = "host=" + config.PgHost + " port=" + config.PgPort + " user=" + config.PgUser + " dbname=" + config.Database + " sslmode=" + config.SSLMode + " password=" + config.PgPassword

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.ErrorLogger.Fatal("Failed to connect database:", err)
	}
}

func Close() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.ErrorLogger.Fatal("Failed to close database:", err)
	}
	sqlDB.Close()
}
