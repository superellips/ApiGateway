package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var guestbookHost string = os.Getenv("GUESTBOOK_HOST")
var userHost string = os.Getenv("USER_HOST")
var aclHost string = os.Getenv("ACL_HOST")

func UnmarshalJsonData(jsonData []byte) (map[string]interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		return nil, err
	}
	return data, nil
}

func MarshalJsonData(jsonMap map[string]interface{}) ([]byte, error) {
	data, err := json.Marshal(jsonMap)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GetReponseData(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	jsonData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return jsonData, err
}

func PostResponseData(url string, body io.Reader) ([]byte, error) {
	response, err := http.Post(url, "application/json", body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	json, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return json, nil
}

func main() {
	hostname := os.Getenv("GUESTBOOK_ROOT_DOMAIN")

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
		private.POST("/guestbook/new", PostNewGuestbook)
		private.GET("/guestbook/:id", GetGuestbook)

		private.POST("/user/register", PostRegisterUser)
		private.GET("/user/id/:id", GetUserById)
		private.GET("/user/name/:name", GetUserByName)
		private.POST("/user/login", PostLoginUser)

	}

	router.Run()
}
