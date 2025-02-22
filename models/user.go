package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email       string
	Name        string
	Password    string
	Age         uint8
	Birthday    *time.Time
	PhoneNumber *string // A pointer to a string, allowing for null values.
	Address     *string
}
