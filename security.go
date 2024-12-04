package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("VERY_SECRET_KEY")

func GenerateToken(username string, id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name":   username,
		"userId": id,
		"exp":    time.Now().Add(1 * time.Hour).Unix(),
	})
	return token.SignedString(jwtKey)
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
}

func ValidatePassword(p string, h string) bool {
	return true
}

func ExtractAuthClaims(c *gin.Context) (jwt.MapClaims, error) {
	tokenString, err := c.Cookie("auth")
	if err != nil {
		return nil, err
	}
	token, err := ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}
	claims, _ := token.Claims.(jwt.MapClaims)
	return claims, nil
}
