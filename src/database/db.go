package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"github.com/fiber/src/models"
	"log"
	"os"
	"time"
)

// Database is the struct to hold DB connection
type Database struct {
	DB *gorm.DB
}

// NewDatabase creates a new database connection
func NewDatabase() (*Database, error) {
	psqlInfo := fmt.Sprintf("host=localhost port=5432 user=root "+
		"password=apl@992132 dbname=myDB sslmode=disable")

	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
		Logger: loggerConfig(true),
	})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		return nil, err
	}

	return &Database{DB: db}, nil
}

// loggerConfig for database
func loggerConfig(enable bool) logger.Interface {
	if enable {
		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second,   // Slow SQL threshold
				LogLevel:                  logger.Info,   // Set log level
				IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound errors
				ParameterizedQueries:      true,          // Don't include raw SQL queries in logs
				Colorful:                  true,          // Colorize logs
			},
		)
		return newLogger
	}

	// Default silent logger if not enabled
	return logger.Default.LogMode(logger.Silent)
}