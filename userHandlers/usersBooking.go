package userhandlers

import (
	"car-service/database"
	"car-service/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Request body struct
type BookingRequest struct {
	VehicleNumber string `json:"vehicle_number" binding:"required"`
	Kilometers    string `json:"kilometers" binding:"required"`
	FuelType      string `json:"fuel_type" binding:"required"`
}

// BookingHandler handles user booking requests
func BookingHandler(c *gin.Context) {
	// Get car ID from URL
	carIDParam := c.Param("id")
	carID, err := strconv.Atoi(carIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid car ID"})
		return
	}

	// Fetch car from DB
	var car models.Cars
	if err := database.DB.First(&car, carID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Car not found"})
		return
	}

	// Get user ID from context (set by middleware)
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userID := userIDRaw.(uint)

	// Parse request body
	var input BookingRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create booking
	booking := models.Booking{
		User_ID:       userID,
		CarID:         uint(carID),
		VehicleNumber: input.VehicleNumber,
		Kilometers:    input.Kilometers,
		FuelType:      input.FuelType,
		Status:        "Pending",
	}

	if err := database.DB.Create(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
		return
	}

	// booking response
	response := models.BookingResponse{
		ID:            booking.ID,
		UserID:        booking.User_ID,
		CarID:         booking.CarID,
		CarBrand:      car.Brand,
		CarModel:      car.Name,
		VehicleNumber: booking.VehicleNumber,
		Kilometers:    booking.Kilometers,
		FuelType:      booking.FuelType,
		Status:        booking.Status,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Booking successful for Tomorrow. Please arrive after 9 AM.",
		"booking": response,
	})
}
