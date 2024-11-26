package main

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetMessages(c *gin.Context) {
	id := c.Param("gbId")

	response, err := http.Get("guestbookservice/api/version/guestbook/" + id + "/messages")
	if err != nil {
		return
	}
	defer response.Body.Close()
	json, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}
	c.IndentedJSON(http.StatusOK, json)
}

func main() {
	router := gin.Default()

	router.GET("/messages/:gbId")
}
