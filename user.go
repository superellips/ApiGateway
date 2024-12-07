package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PostRegisterUser(c *gin.Context) {
	url := "http://" + userHost + "/api/version/users"
	json, err := PostResponseData(url, c.Request.Body)
	if err != nil {
		return
	}
	c.Data(http.StatusOK, "application/json", json)
}

func GetUserById(c *gin.Context) {
	claims, err := ExtractAuthClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization failed"})
		return
	}
	if claims["userId"].(string) != c.Param("id") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	url := "http://" + userHost + "/api/version/user/" + c.Param("id")
	json, err := GetReponseData(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}
	c.Data(http.StatusOK, "application/json", json)
}

func GetUserByName(c *gin.Context) {
	claims, err := ExtractAuthClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization failed"})
		return
	}
	if claims["name"].(string) != c.Param("name") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	url := "http://" + userHost + "/api/version/user/name/" + c.Param("name")
	json, err := GetReponseData(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}
	c.Data(http.StatusOK, "application/json", json)
}

func GetActiveUser(c *gin.Context) {
	claims, err := ExtractAuthClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "no authenticated user"})
		return
	}
	fmt.Println(claims["userId"].(string))
	fmt.Println(claims["name"].(string))
	c.JSON(http.StatusOK, gin.H{"id": claims["userId"].(string), "name": claims["name"].(string)})
}

func PostLoginUser(c *gin.Context) {
	// Add password validation
	loginData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	json, err := UnmarshalJsonData(loginData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	url := "http://" + userHost + "/api/version/user/name/" + json["name"].(string)
	response, err := GetReponseData(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	data, err := UnmarshalJsonData(response)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	if !ValidatePassword("password", "hash") {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "authentication failed"})
		return
	}
	token, err := GenerateToken(data["name"].(string), data["id"].(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.SetCookie("auth", token, 3600, "/", "http://localhost:8888", false, true)
	c.Data(http.StatusOK, "application/json", response)
}
