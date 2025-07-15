package controllers

import (
	"car-service/constant"
	"car-service/database"
	"car-service/models"
	"car-service/tokens"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// staff signin
func StaffSignin(c *gin.Context) {
	var staff struct {
		Name        string `json:"name" binding:"required"`
		Email       string `json:"email" binding:"required,email"`
		Password    string `json:"password" binding:"required"`
		PhoneNumber string `json:"phonenumber" binding:"required"`
	}

	if err := c.ShouldBindJSON(&staff); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	// Check for duplicate user
	var exist models.Register
	err := database.DB.Where("email = ? OR phone_number = ?", staff.Email, staff.PhoneNumber).First(&exist).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists with email or phone number"})
		return
	}
	if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//  Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(staff.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to generate password"})
		return
	}

	// Save to DB
	staffs := models.Register{
		FirstName:   staff.Name,
		Email:       staff.Email,
		Password:    string(hash),
		PhoneNumber: staff.PhoneNumber,
		Role:        constant.Service,
	}

	if err := database.DB.Create(&staffs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to register"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "registration successful", "user": staffs})
}

// staff login
func StaffLogin(c *gin.Context) {
	var staflogin struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&staflogin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid user name or password"})
		return
	}

	var staff models.Register
	if err := database.DB.Where("email=?", staflogin.Email).First(&staff).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "invalid user name or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(staff.Password), []byte(staflogin.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to generate password"})
		return
	}
	acces, err := tokens.GenerateAccessToken(staff.ID, staff.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "failed to generate acces token"})
		return
	}

	ref, err := tokens.GenerateRefreshToken(staff.ID,staff.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to generate Refresh Token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "login sucessfully", "user": staff.FirstName, "Access_token": acces, "Refresh_token": ref})
}

// staff logout
func StaffLogout(c *gin.Context) {
	var input struct {
		Stafflogout string `json:"stafflogout"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	var token models.RefreshToken
	if err := database.DB.Where("token=?", input.Stafflogout).First(&token).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := database.DB.Delete(&token).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "failed to delete token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "logout sucessfully"})
}
