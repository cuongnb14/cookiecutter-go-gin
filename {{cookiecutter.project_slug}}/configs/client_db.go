package configs

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"sync"
	"time"
)

var db *gorm.DB
var onceGetDB sync.Once

func GetDB() *gorm.DB {
	onceGetDB.Do(func() {
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			Env.DBHost, Env.DBPort, Env.DBUser, Env.DBPass, Env.DBName, Env.DBSslMode)

		var err error
		logLevel := glogger.LogLevel(Env.DBLogLevel)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger:                 glogger.Default.LogMode(logLevel),
			TranslateError:         true,
			SkipDefaultTransaction: true,
		})

		if err != nil {
			panic("failed to connect database")
		}

		// Setup connection pool
		sqlDB, err := db.DB()
		if err != nil {
			panic("failed to get database")
		}

		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(500)
		sqlDB.SetConnMaxLifetime(time.Hour)

	})

	return db
}
