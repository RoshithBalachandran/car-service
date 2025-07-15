package controllers

import (
	"car-service/constant"
	"car-service/database"
	"car-service/models"
	"car-service/tokens"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Registeration(c *gin.Context) {
	var input struct {
		FirstName   string `json:"first_name" binding:"required"`
		LastName    string `json:"last_name" binding:"required"`
		PhoneNumber string `json:"phone_number" binding:"required"`
		Email       string `json:"email" binding:"required,email"`
		Address     string `json:"address" binding:"required"`
		Password    string `json:"password" binding:"required"`
	}

	// Validate input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create user record
	user := models.Register{
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		PhoneNumber: input.PhoneNumber,
		Email:       input.Email,
		Address:     input.Address,
		Password:    string(hashedPassword),
		Role:        constant.User, // default role
	}

	// Save user
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Registration successful",
		"user": gin.H{
			"first_name":   user.FirstName,
			"last_name":    user.LastName,
			"address":      user.Address,
			"email":        user.Email,
			"phone_number": user.PhoneNumber,
			"role":         user.Role,
		},
	})
}

// login section and compare user input data

func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	// Validate input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch user by email
	var user models.Register
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Use constants for role check
	var role string
	switch user.Role {
	case constant.Admin:
		role = "admin"
	case constant.User:
		role = "user"
	case constant.Service:
		role = "staff"
	default:
		role = "unknown"
	}

	// Generate JWT tokens
	accessToken, err := tokens.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := tokens.GenerateRefreshToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}
	//save the data to the database
	database.DB.Create(&models.RefreshToken{UserID: user.ID, Token: refreshToken})
	
	// Final response
	c.JSON(http.StatusOK, gin.H{
		"message":       "Login successful",
		"user":          user.FirstName,
		"user_id":       user.ID,
		"role":          role,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func Home(c *gin.Context) {
	var cars []models.Cars
	if err := database.DB.Find(&cars).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "Failed to Get cars"})
		return
	}
	c.JSON(http.StatusOK, cars)
}

// logout
func LogoutUserOrAdmin(c *gin.Context) {
	var input struct {
		UserLogout string `json:"userlogout"`
	}

	// Attempt to get token from Authorization header
	authHeader := c.GetHeader("Authorization")
	var tokenString string

	if !strings.HasPrefix(authHeader, "Bearer ") {
		tokenString = strings.TrimPrefix(authHeader, "Bearer ")
	} else {
		// Fallback to token from JSON body
		if err := c.ShouldBindJSON(&input); err != nil || input.UserLogout == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Logout token is required"})
			return
		}
		tokenString = input.UserLogout
	}

	fmt.Println("Token received:", tokenString)

	// Check if token exists
	var token models.RefreshToken
	if err := database.DB.Where("token = ?", tokenString).First(&token).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	// Delete the token
	if err := database.DB.Delete(&token).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
