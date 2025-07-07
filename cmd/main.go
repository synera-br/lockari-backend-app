package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	// This is the main entry point for the application.
	// You can initialize your application here, set up routes, etc.
	// For example, you might want to start a web server or connect to a database.

	// Example: Start a web server
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, Lockari!")
	})

	// Start the server on port 8080
	if err := router.Run(":8080"); err != nil {
		panic("Failed to start server: " + err.Error())
	}

}
