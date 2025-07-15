package models

import "gorm.io/gorm"

type Cars struct {
	gorm.Model
	Brand       string `json:"brand"`
	Name        string `json:"name"`
	Fuel_type   string `json:"fuel_type"`
	Description string `json:"description"`
	Image       string `json:"image"`
}
