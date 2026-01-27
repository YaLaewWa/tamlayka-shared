package db

import (
	"fmt"
	"log"
	"math"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConfig struct {
	Host     string `env:"HOST,required"`
	Port     int    `env:"PORT,required"`
	User     string `env:"USER,required"`
	Password string `env:"PASSWORD,required"`
	DBName   string `env:"NAME,required"`
	SSLMode  string `env:"SSLMODE,required"`
}

func New(config DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("failed to connect database")
		return nil, err
	}

	return db, nil
}

// Pagination use with gorm.DB.Scopes() to calculate total pages and last page
//
// Usage Example:
//
//	db.Scopes(database.Paginate(domain.Book{}, limit, page, total, last)).Find(&books)
func Paginate(limit, page, total, last *int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var totalRows int64

		// Clone the current query (keeps WHERE, JOIN, etc.)
		countDB := db.Session(&gorm.Session{})

		countDB.Count(&totalRows)

		*total = int(totalRows)
		*last = int(math.Ceil(float64(totalRows) / float64(*limit)))

		offset := (*page - 1) * *limit
		return db.Offset(offset).Limit(*limit)
	}
}
