package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/noverde/pequi/pkg/store"
)

type payload struct {
	URL string `json:"long_url" binding:"required"`
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

func normalizeURLString(url string) string {
	if url == "" {
		url = "localhost"
	}
	if !strings.HasSuffix(url, "/") {
		url = url + "/"
	}
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}
	return url
}

func main() {
	driver := os.Getenv("STORAGE_DRIVER")
	if driver == "" {
		driver = "memory"
	}

	store, err := store.New(driver)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	address := os.Getenv("SERVER_ADDRESS")
	if address == "" {
		address = "http://localhost" + port
	}
	address = normalizeURLString(address)

	r := gin.Default()
	r.GET("/:hash", func(c *gin.Context) {
		var route, err = store.Get(c.Param("hash"))
		if err != nil {
			c.Data(http.StatusNotFound, "text/html; charset=utf-8", []byte("404 Not Found\n"))
			return
		}
		c.Redirect(http.StatusMovedPermanently, route)
	})

	r.POST("/v1/shorten", func(c *gin.Context) {
		if !store.Auth(c.GetHeader("authorization")) {
			c.JSON(http.StatusUnauthorized, gin.H{"Unauthorized": "Invalid authorization token on header"})
			return
		}

		var data payload
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		slug, err := store.Set(data.URL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"slug": slug, "link": address + slug})
	})
	log.Fatal(r.Run(port))
}
