package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	sslMode = "disable" // would enable this in production
	tries   = 10
)

func NewProdDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
		sslMode,
	)

	i := 0
	var err error
	for i < tries {
		var db *gorm.DB
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			return db, err
		}

		fmt.Println("Attempted to connect to db...retrying")
		time.Sleep(2 * time.Second)

		i++
	}

	return nil, err
}

func NewUnitDB() (*gorm.DB, string, error) {
	u, _ := uuid.NewV4()
	fileName := fmt.Sprintf("unit-test-%s.db", u.String())
	db, err := gorm.Open(sqlite.Open(fileName), &gorm.Config{})
	return db, fileName, err
}
