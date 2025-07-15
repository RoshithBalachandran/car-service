package models

import (
	"gorm.io/gorm"
)

// registur struct
type Register struct {
	gorm.Model
	FirstName   string    `json:"first_name" binding:"required"`
	LastName    string    `json:"last_name" binding:"required"`
	PhoneNumber string    `json:"phone_number" binding:"required" gorm:"unique"`
	Email       string    `json:"email" binding:"required,email" gorm:"unique"`
	Password    string    `json:"password" binding:"required"`
	Address     string    `json:"address" binding:"required"`
	Role        string    `json:"role"`
	// Bookings    []Booking `gorm:"foreignKey:UserID" json:"bookings"`
}
