package db

import (
	"context"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func OpenConnection(dbUrl string, maxOpenConns int, maxIdleConns int, maxIdleTime string) (*gorm.DB, error) {
	dialect := postgres.Open(dbUrl)
	db, err := gorm.Open(dialect, &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Info),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		// IMPLEMENT SELECT FIELD
		// IMPLEMENT LAZY RESULTS USING ROWS() INSTEAD OF FIND() - JIKA DATA BESAR
		// IMPLEMENT TABLE NORMALIZATION (TABLE SPLITTING)
	})
	if err != nil {
		return nil, err
	}
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	duration, _ := time.ParseDuration(maxIdleTime)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxIdleTime(duration)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	return db, nil
}
