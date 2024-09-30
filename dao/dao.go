package dao

import (
	"server/config"
	"server/pkg/logger"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// get DB reference
var DB *gorm.DB

func Init() {
	// config the neon
	dsn := config.ProgresSQL
	var err error
	//connect neon
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// print error if occured
	if err != nil {
		logger.Error(map[string]interface{}{"sql is failed to connect": err.Error()})
	}

	sqlDB, err := DB.DB()
	if err != nil {
		logger.Error(map[string]interface{}{"sql is failed to connect": err.Error()})
	}
	// set num of idle connections
	sqlDB.SetMaxIdleConns(10)
	// set num of max connections
	sqlDB.SetMaxOpenConns(100)
	//set the duration of connection life
	sqlDB.SetConnMaxLifetime(time.Hour)
}
