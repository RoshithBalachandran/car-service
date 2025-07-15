package controllers

import (
	"car-service/database"
	"car-service/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Add new cars

func AddCar(c *gin.Context) {
	var car models.Cars
	if err := c.ShouldBindJSON(&car); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&car).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add car"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Car added successfully", "car": car})
}

// Get All cars
func AdminGetAll(c *gin.Context) {
	var cars []models.Cars
	if err := database.DB.Find(&cars).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "Failed to Get cars"})
		return
	}
	c.JSON(http.StatusOK, cars)
}

func UpdateCar(c *gin.Context) {
	id := c.Param("id")

	var Car models.Cars

	if err := database.DB.Find(&Car, &id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Data not found"})
		return
	}

	var updateCar models.Cars

	if err := c.ShouldBindJSON(&updateCar); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	Car.Name = updateCar.Name
	Car.Description = updateCar.Description
	Car.Fuel_type = updateCar.Fuel_type
	Car.Image = updateCar.Image
	
	//updates on databases
	if err := database.DB.Save(&Car).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to update car"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Updated sucessfully",
		"car":     Car,
	})
}

func DeleteCar(c *gin.Context) {
	id := c.Param("id")

	if err := database.DB.Delete(&models.Cars{}, &id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed To delete car"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Data deleted sucessfully", "id": id})
}
