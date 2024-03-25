package utils

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var secretKey = os.Getenv("SECRET_KEY")

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		token, err := ExtractTokenID(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Set("userID", token)
		c.Next()
	}
}

func GenerateToken(userID uint) (string, error) {
	token_lifespan, err := strconv.Atoi(os.Getenv("TOKEN_LIFESPAN"))
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * time.Duration(token_lifespan)).Unix(),
	})
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ExtractTokenID(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return 0, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := uint(claims["user_id"].(float64))
		return userID, nil
	}
	return 0, fmt.Errorf("Invalid token")
}

func TokenCookie(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
	c.SetCookie("token", tokenString, 0, "/", "localhost", false, true)
}

//// ExtractTokenID extracts the user ID from the JWT token
//func ExtractTokenID(tokenString string) (uint, error) {
//	// Parse the token
//	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
//		// Check the signing method
//		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
//		}
//
//		// Return the secret key
//		return []byte(secretKey), nil
//	})
//	if err != nil {
//		return 0, err
//	}
//
//	// Check if the token is valid
//	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
//		userID := uint(claims["user_id"].(float64))
//		return userID, nil
//	}
//
//	return 0, fmt.Errorf("Invalid token")
//}
