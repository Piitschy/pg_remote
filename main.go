package main

import (
	"net/http"
	"os"
	"strconv"

	pg "github.com/habx/pg-commands"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/Piitschy/postgress-dump-tool/docs"
)

// @title Echo Swagger Example API
// @version 1.0
// @description This is a sample server server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /
// @schemes http
func main() {
	// Get env variables
	db := &pg.Postgres{
		Host:     os.Getenv("DB_HOST"),
		Port:     strconv.Atoi(os.Getenv("DB_PORT")),
		DB:       os.Getenv("DB_NAME"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
	}

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())

	// Routes
	e.GET("/", HealthCheck)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Start server
	e.Logger.Fatal(e.Start(":3000"))
}

// HealthCheck godoc
// @Router / [get]
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} Message
func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, Message{
		Data: "Server is up and running",
	})
}

func Dump(c echo.Context) error {
	return c.JSON(http.StatusOK, Message{
		Data: "Server is up and running",
	})
}
