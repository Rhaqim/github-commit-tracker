package database

import (
	"github.com/Rhaqim/savannahtech/config"
	"github.com/Rhaqim/savannahtech/internal/core/repo"
	"github.com/Rhaqim/savannahtech/internal/core/repositories"
	"github.com/Rhaqim/savannahtech/pkg/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB // Global database connection

func InitDB() {
	dsn := config.Config.DatabaseURL

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.ErrorLogger.Fatalf("failed to connect to database: %v", err)
	}

	DB = db

	// Initialize repositories
	repositories.RepoStore = repo.NewRepositoryRepo(DB)
	repositories.CommitStore = repo.NewCommitRepo(DB)
}

// func Init() {
// 	var err error

// 	logger.InfoLogger.Println("Connecting to PostgreSQL database...")

// 	var dsn string = "host=" + config.PgHost + " port=" + config.PgPort + " user=" + config.PgUser + " dbname=" + config.Database + " sslmode=" + config.SSLMode + " password=" + config.PgPassword

// 	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
// 		Logger: logger.Default.LogMode(logger.Silent),
// 	})
// 	if err != nil {
// 		logger.ErrorLogger.Fatal("Failed to connect database:", err)
// 	}
// }

// func Close() {
// 	sqlDB, err := DB.DB()
// 	if err != nil {
// 		logger.ErrorLogger.Fatal("Failed to close database:", err)
// 	}
// 	sqlDB.Close()
// }
