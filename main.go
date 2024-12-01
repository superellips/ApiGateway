package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var guestbookHost string = os.Getenv("GUESTBOOK_HOST")
var userHost string = os.Getenv("USER_HOST")
var aclHost string = os.Getenv("ACL_HOST")

func PostNewGuestbook(c *gin.Context) {
	url := "http://" + guestbookHost + "/api/version/guestbook"
	// domain
	// requireApproval
	// ownerId
	response, err := http.Post(url, "application/json", c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create guestbook"})
		return
	}
	defer response.Body.Close()
	jsonData, err := io.ReadAll(response.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response data"})
		return
	}
	c.Data(response.StatusCode, "application/json", jsonData)
}

func GetMessages(c *gin.Context) {
	id := c.Param("gbId")
	if !IsValidOrigin(c, id) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden origin."})
		return
	}
	url := "http://" + guestbookHost + "/api/ersion/guestbook/" + id + "/messages"
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
	url := "http://" + guestbookHost + "/api/version/guestbook/" + id + "/message"
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

func PostRegisterUser(c *gin.Context) {
	url := "http://" + userHost + "/api/version/users"
	response, err := http.Post(url, "application/json", c.Request.Body)
	if err != nil {
		return
	}
	defer response.Body.Close()
	json, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}
	c.Data(response.StatusCode, "application/json", json)
}

func GetUserById(c *gin.Context) {
	id := c.Param("id")
	url := "http://" + userHost + "/api/version/user/" + id
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

func GetUserByName(c *gin.Context) {
	name := c.Param("name")
	url := "http://" + userHost + "/api/version/user/name/" + name
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
	url := "http://" + guestbookHost + "/api/version/guestbook/" + id
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
	hostname := "http://" + os.Getenv("GUESTBOOK_ROOT_DOMAIN")

	router := gin.Default()

	relaxedCors := cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:    []string{"Content-Type", "Origin"},
		MaxAge:          12 * time.Hour,
	})

	strictCors := cors.New(cors.Config{
		AllowOrigins:     []string{hostname},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type", "Origin"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})

	public := router.Group("/public")
	public.Use(relaxedCors)
	{
		public.GET("/messages/:gbId", GetMessages)
		public.POST("/messages/:gbId", PostMessage)
	}

	private := router.Group("/")
	private.Use(strictCors)
	{
		private.GET("/messages/:gbId", GetMessages)
		private.POST("/messages/:gbId", PostMessage)

		private.POST("/user/register", PostRegisterUser)
		private.GET("/user/id/:id", GetUserById)
		private.GET("/user/name/:name", GetUserByName)

		private.POST("/guestbook/new", PostNewGuestbook)
	}

	router.Run()
}
