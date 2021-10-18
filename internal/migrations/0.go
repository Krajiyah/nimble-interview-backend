package migrations

import (
	"github.com/Krajiyah/nimble-interview-backend/internal/models"
	"gorm.io/gorm"
)

func initializeDB(db *gorm.DB) error {
	m := db.Migrator()

	if !m.HasTable(&models.User{}) {
		if err := m.CreateTable(&models.User{}); err != nil {
			return err
		}
	}

	if !m.HasTable(&models.Message{}) {
		if err := m.CreateTable(&models.Message{}); err != nil {
			return err
		}
	}

	return nil
}
