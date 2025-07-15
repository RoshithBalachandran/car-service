package tokens

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// getting the env secret key
var SECRET_KEY = []byte(os.Getenv("JWT_SECRET"))

// claims struct
type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateAccessToken creates a short-lived JWT (15 mins)
func GenerateAccessToken(userID uint, role string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "jwt-auth-app",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SECRET_KEY)
}

// GenerateRefreshToken creates a long-lived token (7 days)
func GenerateRefreshToken(userID uint, Role string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Role:   Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "jwt-auth-app",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SECRET_KEY)
}

//validate token

func Validatetoken(tokenstr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenstr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return SECRET_KEY, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return claims, nil
}
