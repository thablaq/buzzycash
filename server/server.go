package server

import (
	"fmt"
	"net/http"

	"github.com/dblaq/buzzycash/internal/config"
	"github.com/dblaq/buzzycash/internal/middlewares"
	"github.com/gin-gonic/gin"
)

// NewServer initializes the Gin engine and middleware
func NewServer() *gin.Engine {
	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middlewares.RecoveryAndErrorMiddleware())

	// Serve static files
	r.Static("/uploads/profile-pictures", "./uploads/profile-pictures")

	return r
}

// StartServer runs the server
func StartServer(r *gin.Engine) {
	fmt.Println("🚀 Server started on :" + config.AppConfig.Port)
	r.Run(":" + config.AppConfig.Port)
}

// HealthCheck adds a basic welcome route
func HealthCheck(r *gin.Engine) {
	r.GET("/api/v1/welcome", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Welcome to BuzzyCash API")
	})
}
