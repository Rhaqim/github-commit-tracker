package database

import (
	"fmt"
	"savannahtech/src/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB // Global database connection

// InitDB initializes the database connection using the configuration values from the config package.
// It establishes a connection to the PostgreSQL database and assigns the connection to the global DB variable.
// If an error occurs during the connection process, it logs the error and shuts down the logger.
func Init() {
	var err error

	fmt.Println("Connecting to PostgreSQL database...", config.PgHost, config.PgPort, config.PgUser, config.PgPassword, config.Database, config.SSLMode)

	var dsn string = "host=" + config.PgHost + " port=" + config.PgPort + " user=" + config.PgUser + " dbname=" + config.Database + " sslmode=" + config.SSLMode + " password=" + config.PgPassword

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed to connect database")
	}
}

func Close() {
	sqlDB, err := DB.DB()
	if err != nil {
		panic("failed to close database")
	}
	sqlDB.Close()
}
