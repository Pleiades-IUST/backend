package webservice

import (
	"net/http"

	"github.com/Pleiades-IUST/backend/utils/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	gin.SetMode(config.GetGinMode())
	e := gin.Default()

	// Add CORS middleware
	e.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Ping test
	e.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	return e
}
