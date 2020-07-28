package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/teris-io/shortid"
)

type payload struct {
	LongURL   string `json:"long_url" binding:"required"`
	Domain    string `json:"domain,omitempty"`
	ShortPath string `json:"short_path,omitempty"`
	Overwrite bool   `json:"overwrite,omitempty"`
}

type header struct {
	Authorization string `json:"authorization" binding:"required"`
}

var port string

func init() {
	// Disable log prefixes such as the default timestamp.
	// Prefix text prevents the message from being parsed as JSON.
	// A timestamp is added when shipping logs to Stackdriver.
	log.SetFlags(0)

	// Load environment vars from .env, the system env vars has
	// precedence over .env
	err := godotenv.Load()
	if err == nil {
		log.Print("Loading environment vars from .env file.")
	}

	// HTTP listen port
	port = ":8080"
	if value, ok := os.LookupEnv("HTTP_PORT"); ok {
		port = value
	}
	if port[0] != ':' {
		port = ":" + port
	}
}

func main() {
	r := gin.Default()
	r.GET("/:path", func(c *gin.Context) {
		var route, err = getRoute(c.Param("path"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		c.Redirect(http.StatusMovedPermanently, route)
	})

	r.POST("/v1/shorten", func(c *gin.Context) {
		var auth header
		if err := c.BindHeader(&auth); err != nil || !isKeyValid(auth.Authorization) {
			c.JSON(http.StatusUnauthorized, gin.H{"Unauthorized": "Invalid authorization token on header"})
			return
		}

		var data payload
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var id, err = processAndSave(&data)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"ShortPath": id})
	})
	log.Fatal(r.Run(port))
}

//TODO: Implement Later
func isKeyValid(key string) bool {
	return true
}

//TODO: Implement Later
func processAndSave(data *payload) (string, error) {
	return shortid.Generate()
}

//TODO: Implement Later
func getRoute(path string) (string, error) {
	return "https://google.com/?q=" + path, nil
}
