package main

import (
	"github.com/dblaq/buzzycash/http"
	"github.com/dblaq/buzzycash/docs"
	"github.com/dblaq/buzzycash/internal/config"
	"github.com/dblaq/buzzycash/server"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	config.LoadConfig()
	config.InitDB()
	defer config.CloseDB()

	r := server.NewServer()

	// Swagger setup
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Host = "localhost:5005"
	url := ginSwagger.URL("/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	server.HealthCheck(r)
	http.RegisterRoutes(r, config.DB)

	server.StartServer(r)
}
