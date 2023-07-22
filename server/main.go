package main

import (
	"fmt"
	"os"

	pg "github.com/habx/pg-commands"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/Piitschy/postgress-dump-tool/server/docs"
)

var DB *pg.Postgres

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

	DB = NewPostgres()

	fmt.Println("DB_HOST:", DB.Host, "DB_PORT:", DB.Port, "DB_NAME:", DB.DB, "DB_USER:", DB.Username)

	// Echo instance
	e := echo.New()
	e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "header:Key",
		Validator: func(key string, c echo.Context) (bool, error) {
			return key == os.Getenv("KEY"), nil
		},
	}))

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())

	// Routes
	e.GET("/", HealthCheck)
	e.POST("/dump", DumpRoute)
	e.GET("/docs/*", echoSwagger.WrapHandler)

	// Start server
	e.Logger.Fatal(e.Start(":3000"))
}
