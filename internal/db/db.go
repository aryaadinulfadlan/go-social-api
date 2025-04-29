package db

import (
	"context"
	"log"
	"time"

	"github.com/aryaadinulfadlan/go-social-api/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init() {
	dialect := postgres.Open(config.DB.DATABASE_URL)
	var err error
	DB, err = gorm.Open(dialect, &gorm.Config{
		Logger:      logger.Default.LogMode(logger.Info),
		PrepareStmt: true,
	})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	sqlDB, errSql := DB.DB()
	if err != nil {
		log.Fatalf("Failed to connect database: %v", errSql)
	}
	duration, _ := time.ParseDuration(config.DB.MaxIdleTime)
	sqlDB.SetMaxOpenConns(config.DB.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.DB.MaxIdleConns)
	sqlDB.SetConnMaxIdleTime(duration)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
}
