package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func GetMessages(c *gin.Context) {
	id := c.Param("gbId")
	if !IsValidOrigin(c, id) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden origin."})
		return
	}
	url := "http://guestbookservice:8080/api/version/guestbook/" + id + "/messages"
	response, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data"})
		return
	}
	defer response.Body.Close()
	json, err := io.ReadAll(response.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}
	c.Data(response.StatusCode, "application/json", json)
}

func PostMessage(c *gin.Context) {
	id := c.Param("gbId")
	if !IsValidOrigin(c, id) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden origin."})
		return
	}
	url := "http://guestbookservice:8080/api/version/guestbook/" + id + "/message"
	response, err := http.Post(url, "application/json", c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to access service."})
		return
	}
	defer response.Body.Close()
	json, err := io.ReadAll(response.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}
	c.Data(response.StatusCode, "application/json", json)
}

func IsValidOrigin(c *gin.Context, id string) bool {
	origin := c.GetHeader("Origin")
	if origin == "" {
		origin = c.GetHeader("Referer")
	}
	parsedOrigin, err := url.Parse(origin)
	if err != nil {
		return false
	}
	validHostname, err := GetGuestbookAllowedDomain(id)
	if err != nil {
		return false
	}
	usedHostname := parsedOrigin.Hostname()
	return validHostname == usedHostname
}

func GetGuestbookAllowedDomain(id string) (string, error) {
	url := "http://guestbookservice:8080/api/version/guestbook/" + id
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	jsonData, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		panic(err)
	}

	domain, exists := data["domain"].(string)
	if !exists {
		return "", errors.New("allowed domain missing")
	}
	return domain, nil
}

func main() {
	router := gin.Default()

	defaultCors := cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	})

	router.Use(defaultCors)

	router.GET("/messages/:gbId", GetMessages)
	router.POST("/messages/:gbId", PostMessage)

	router.Run()
}
