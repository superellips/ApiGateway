package main

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func PostNewGuestbook(c *gin.Context) {
	claims, err := ExtractAuthClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization failed"})
		return
	}
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request"})
		return
	}
	dataNewGb, err := UnmarshalJsonData(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request"})
		return
	}
	dataNewGb["ownerId"] = claims["userId"]
	data, err = MarshalJsonData(dataNewGb)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write request"})
		return
	}
	url := "http://" + guestbookHost + "/api/version/guestbook"

	jsonData, err := PostResponseData(url, bytes.NewBuffer(data))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response data"})
		return
	}
	c.Data(http.StatusOK, "application/json", jsonData)
}

func GetMessages(c *gin.Context) {
	id := c.Param("gbId")
	if !isValidOrigin(c, id) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden origin."})
		return
	}
	url := "http://" + guestbookHost + "/api/version/guestbook/" + id + "/messages"
	json, err := GetReponseData(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}
	c.Data(http.StatusOK, "application/json", json)
}

func PostMessage(c *gin.Context) {
	id := c.Param("gbId")
	if !isValidOrigin(c, id) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden origin."})
		return
	}
	url := "http://" + guestbookHost + "/api/version/guestbook/" + id + "/message"
	json, err := PostResponseData(url, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}
	c.Data(http.StatusOK, "application/json", json)
}

func OptionsMessage(c *gin.Context) {
	id := c.Param("gbId")
	if !isValidOrigin(c, id) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden origin"})
		return
	}
}

func GetGuestbookAllowedDomain(id string) (string, error) {
	url := "http://" + guestbookHost + "/api/version/guestbook/" + id
	jsonData, err := GetReponseData(url)
	if err != nil {
		panic(err)
	}
	data, err := UnmarshalJsonData(jsonData)
	if err != nil {
		panic(err)
	}
	domain, exists := data["domain"].(string)
	if !exists {
		return "", errors.New("allowed domain missing")
	}
	return domain, nil
}

func GetGuestbook(c *gin.Context) {
	url := "http://" + guestbookHost + "/api/version/guestbook/" + c.Param("id")
	json, err := GetReponseData(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.Data(http.StatusOK, "application/json", json)
}

func isValidOrigin(c *gin.Context, id string) bool {
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
