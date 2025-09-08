package config

import (
	"fmt"
	"github.com/dblaq/buzzycash/internal/migrations"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	var err error




	dsn := AppConfig.DbUrl
	if dsn == "" {
		panic("❌ DATABASE_URL is not set")
	}

	// Configure custom logger
	newLogger := gormLogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		gormLogger.Config{
			SlowThreshold: time.Second,     
			LogLevel:      gormLogger.Silent, 
			Colorful:      true,
		},
	)

	// Open DB connection with new logger
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("❌ Failed to connect database")
	}
	

	// Run migrations
	migrations.AutoMigrate(DB)
	fmt.Println("✅ Database connected & migrated")
}