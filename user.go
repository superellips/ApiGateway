package main

import (
	"encoding/json"
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
	userUrl := "http://" + userHost + "/api/version/user/" + c.Param("id")
	userData, err := GetReponseData(userUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}
	userJson, err := UnmarshalJsonData(userData)
	if err != nil {
		return
	}
	// Get the users available guestbooks too please
	aclUrl := "http://" + aclHost + "/api/version/acls/user/" + claims["userId"].(string)
	aclData, err := GetReponseData(aclUrl)
	if err != nil {
		return
	}
	var aclList []map[string]string
	if err := json.Unmarshal(aclData, &aclList); err != nil {
		return
	}
	guestbooks := make(map[string]string)
	for _, val := range aclList {
		gbUrl := "http://" + guestbookHost + "/api/version/guestbook/" + val["guestbookId"]
		gbData, err := GetReponseData(gbUrl)
		if err != nil {
			return
		}
		gbJson, err := UnmarshalJsonData(gbData)
		if err != nil {
			return
		}
		guestbooks[val["guestbookId"]] = gbJson["domain"].(string)
	}
	userJson["guestbooks"] = guestbooks
	returnData, err := MarshalJsonData(userJson)
	if err != nil {
		return
	}
	c.Data(http.StatusOK, "application/json", returnData)
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
