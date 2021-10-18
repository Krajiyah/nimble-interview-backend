package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string `json:"username" gorm:"unique_index"`
	PasswordHash string `json:"-"`
}

type Message struct {
	gorm.Model
	Data     string `json:"data"`
	Username string `json:"username"`
}

func GetUserByID(db *gorm.DB, id uint) (*User, error) {
	result := &User{}
	if err := db.First(result, id).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func GetUserByUsername(db *gorm.DB, username string) (*User, error) {
	result := &User{}
	if err := db.Where("username", username).First(result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func NewUser(db *gorm.DB, result *User) (*User, error) {
	return result, db.Create(result).Error
}

func NewMessage(db *gorm.DB, result *Message) (*Message, error) {
	return result, db.Create(result).Error
}
