package middleware

import (
	"car-service/database"
	"car-service/models"
	"car-service/tokens"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMid(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		const prefix = "Bearer "

		if !strings.HasPrefix(authHeader, prefix) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing or malformed"})
			ctx.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, prefix)
		claims := &tokens.Claims{}

		// Parse and validate token
		tokenObj, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return tokens.SECRET_KEY, nil
		})

		if err != nil || tokenObj == nil || !tokenObj.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid or expired token",
				"details": err.Error(),
			})
			ctx.Abort()
			return
		}

		// 🔐 Check token exists in DB (i.e., not logged out)
		var storedToken models.RefreshToken
		if err := database.DB.Where("token = ?", tokenStr).First(&storedToken).Error; err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token has been revoked or does not exist"})
			ctx.Abort()
			return
		}

		// ✅ Role check
		roleAllowed := false
		for _, val := range allowedRoles {
			if strings.EqualFold(claims.Role, val) {
				roleAllowed = true
				break
			}
		}

		if !roleAllowed {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "No permission granted"})
			ctx.Abort()
			return
		}

		// ✅ Set context
		ctx.Set("user_id", claims.UserID)
		ctx.Set("role", claims.Role)
		ctx.Next()
	}
}
