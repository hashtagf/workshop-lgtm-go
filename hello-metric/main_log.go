package main

// https://github.com/samber/slog-gin
import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func callFunc(number int) {
	logger.Info("Calling some function", "function", "callFunc", "number", number)
}

func main() {
	// Create a new Gin router
	router := gin.New()
	router.Use(sloggin.New(logger))
	router.Use(gin.Recovery())

	// Define a simple ping endpoint
	router.GET("/ping", func(c *gin.Context) {
		logger.Info("Get /ping", "method", c.Request.Method, "path", c.Request.URL.Path, "status", http.StatusOK)
		callFunc(1)
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.Run(":1234")
}
