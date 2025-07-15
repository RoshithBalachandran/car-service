package userhandlers

// import (
// 	"car-service/database"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// )

// func History(c *gin.Context) {
// 	user, exist := c.Get("user_id")

// 	if !exist{
// 		c.JSON(http.StatusUnauthorized,gin.H{"Error":"unauthorized user"})
// 		return
// 	}
// 	userid ,ok := user.(uint)

// 	if !ok{
// 		c.JSON(http.StatusBadRequest,gin.H{"error":"invalid user"})
// 		return
// 	}

// 	if err:=database.DB.Where("user_id=?",userid).Preload("Car").Order("created_at~")	
		
// }