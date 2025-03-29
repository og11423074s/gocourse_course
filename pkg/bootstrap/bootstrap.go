package bootstrap

import (
	"fmt"
	"github.com/og11423074s/gocourse_domain/domain"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

func DBConnection() (*gorm.DB, error) {
	// Connect to database
	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"))

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}
	if os.Getenv("DATABASE_DEBUG") == "true" {
		db = db.Debug()
	}

	if os.Getenv("DATABASE_MIGRATE") == "true" {

		if err := db.AutoMigrate(&domain.Course{}); err != nil {
			return nil, err
		}

	}

	return db, nil
}

func InitLogger() *log.Logger {
	return log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
}
