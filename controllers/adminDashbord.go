package controllers

import (
	"car-service/database"
	"car-service/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func AdminDashbord(c *gin.Context) {

	var todayBooking []models.Booking
	// Get today's date in correct format
	today := time.Now().Format("2006-01-02")

	// Fetch today's bookings
	if err := database.DB.
		Preload("User").
		Preload("Car").
		Where("DATE(created_at) = ?", today).
		Find(&todayBooking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch today's bookings"})
		return
	}

	// Send response
	c.JSON(http.StatusOK, gin.H{
		"message":         "Bookings fetched successfully",
		"todays_bookings": todayBooking,
		"today_count":     len(todayBooking),
	})
}
