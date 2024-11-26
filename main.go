package main

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetMessages(c *gin.Context) {
	id := c.Param("gbId")
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

func main() {
	router := gin.Default()

	router.GET("/messages/:gbId", GetMessages)
	router.POST("/messages/:gbId", PostMessage)

	router.Run()
}
