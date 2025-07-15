package userhandlers

import (
	"car-service/database"
	"car-service/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Profile(c *gin.Context) {
	//  Get user ID from context
	userId, exists := c.MustGet("user_id").(uint)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User ID not found in context"})
		return
	}

	//Find the user from DB
	var user models.Register
	if err := database.DB.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	//Send user data as response
	c.JSON(http.StatusOK, gin.H{"profile": user})
}

func ProfileUpdate(c *gin.Context) {
	var update models.Register

	// Step 1: Bind JSON input
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Step 2: Get user ID from context (JWT middleware must have set it)
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}
	userID := userId.(uint)

	// Step 3: Find the existing user
	var profile models.Register
	if err := database.DB.First(&profile, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Step 4: Update allowed fields
	profile.FirstName = update.FirstName
	profile.LastName = update.LastName
	profile.PhoneNumber = update.PhoneNumber
	profile.Address = update.Address

	// Step 5: Save to DB
	if err := database.DB.Save(&profile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	// Step 6: Send updated profile
	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"profile": profile,
	})
}
