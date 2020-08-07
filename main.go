package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	storeInit()
	defer storeClose() // Without gin graceful stop it is not really working

	r := gin.Default()
	r.GET("/:hash", func(c *gin.Context) {
		var route, err = storeGet(c.Param("hash"))
		if err != nil {
			c.Data(http.StatusNotFound, "text/html; charset=utf-8", []byte("404 Not Found\n"))
			return
		}
		c.Redirect(http.StatusMovedPermanently, route)
	})

	r.POST("/v1/shorten", func(c *gin.Context) {
		if !storeAuth(c.GetHeader("authorization")) {
			c.JSON(http.StatusUnauthorized, gin.H{"Unauthorized": "Invalid authorization token on header"})
			return
		}

		var data payload
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var id, err = storePut(&data)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"ShortPath": id})
	})
	log.Fatal(r.Run(port))
}
