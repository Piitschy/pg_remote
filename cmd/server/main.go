package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/Piitschy/pgrd/cmd/server/docs"
	db "github.com/Piitschy/pgrd/internal/db"
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
	// Load db from env
	// DB_HOST
	// DB_DATABASE
	// DB_USER
	// DB_PASSWORD
	// DB_PORT
	db := db.NewPostgresFromEnv()

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Logger.SetLevel(logLevelFromEnv())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	keyAuth := middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "header:Key",
		Validator: func(key string, c echo.Context) (bool, error) {
			return key == os.Getenv("KEY"), nil
		},
	})

	dir := dumpDir()
	err := os.Mkdir(dir, os.ModePerm)
	if err != nil {
		panic(fmt.Sprintf("Could not create dumps directory: %s", err))
	}

	// Routes
	e.GET("/", HealthCheck, keyAuth)
	e.POST("/dump", DumpRoute(db, dir), keyAuth)
	e.POST("/restore", RestoreRoute(db), keyAuth)
	e.GET("/docs", func(c echo.Context) error {
		return c.Redirect(301, "/docs/index.html")
	})
	e.GET("/docs/*", echoSwagger.WrapHandler)

	// Start server
	e.Logger.Fatal(e.Start(":3000"))
}

func dumpDir() string {
	dir := os.Getenv("DUMP_DIR")
	if dir == "" {
		dir = filepath.Join(".", "dumps")
	}
	return dir
}

func logLevelFromEnv() log.Lvl {
	ll := os.Getenv("LOG_LEVEL")
	switch strings.ToUpper(ll) {
	case "DEBUG":
		return log.DEBUG
	case "INFO":
		return log.INFO
	case "WARN":
		return log.WARN
	case "ERROR":
		return log.ERROR
	default:
		return log.WARN
	}
}
