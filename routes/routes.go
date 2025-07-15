package routes

import (
	"car-service/constant"
	"car-service/controllers"
	"car-service/middleware"
	userhandlers "car-service/userHandlers"

	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	// Public routes
	r.POST("/register", controllers.Registeration)
	r.POST("/login", controllers.Login)
	r.GET("/home", controllers.Home)

	// Protected user routes (role = constant.User)
	userGroup := r.Group("/user")
	userGroup.Use(middleware.AuthMid(constant.User))
	{
		userGroup.POST("/booking/:id", userhandlers.BookingHandler)
		// userGroup.GET("/history", controllers.HistoryHandler)
		userGroup.GET("/profile", userhandlers.Profile)
		userGroup.PUT("/update", userhandlers.ProfileUpdate)
		userGroup.POST("/logout", controllers.LogoutUserOrAdmin)
	}

	// Admin routes
	adminGroup := r.Group("/admin")
	adminGroup.POST("/login", controllers.Login) // public
	adminGroup.Use(middleware.AuthMid(constant.Admin))
	{
		adminGroup.GET("/dashbord",controllers.AdminDashbord)
		adminGroup.POST("/car", controllers.AddCar)
		adminGroup.GET("/cars", controllers.AdminGetAll)
		adminGroup.PUT("/car/:id", controllers.UpdateCar)
		adminGroup.DELETE("/car/:id", controllers.DeleteCar)
		adminGroup.POST("/logout", controllers.LogoutUserOrAdmin)
	}

	staffGroup := r.Group("/staff")
	staffGroup.POST("/signin", controllers.StaffSignin)
	staffGroup.POST("/login", controllers.StaffLogin)
	staffGroup.Use(middleware.AuthMid(constant.Service))
	{
		staffGroup.GET("/dashbord")
		staffGroup.GET("/service/:id",)
		staffGroup.POST("/logout",controllers.StaffLogout)
	}
}
