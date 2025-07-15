package models

import "gorm.io/gorm"

//booking struct
type Booking struct {
	gorm.Model

	User_ID uint     `json:"user_id"`                            
	User    Register `gorm:"foreignKey:User_ID"`                

	CarID uint `json:"car_id"`                                   
	Car   Cars `gorm:"foreignKey:CarID"`                         

	VehicleNumber string `json:"vehicle_number"`
	Kilometers    string `json:"kilometers"`
	FuelType      string `json:"fuel_type"`
	Status        string `json:"status" gorm:"default:'Pending'"`
}


//booking responce struct
type BookingResponse struct {
	ID            uint   `json:"id"`
	UserID        uint   `json:"user_id"`
	CarID         uint   `json:"car_id"`
	CarBrand      string `json:"car_brand"`
	CarModel      string `json:"car_model"`
	VehicleNumber string `json:"vehicle_number"`
	Kilometers    string `json:"kilometers"`
	FuelType      string `json:"fuel_type"`
	Status        string `json:"status"`
}