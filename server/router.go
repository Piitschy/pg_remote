package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// HealthCheck godoc
// @Router / [get]
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} Response
func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, Response{
		Msg: "Server is up and running",
	})
}

// Dump DB
// @Router /dump [post]
// @Summary Dump the database.
// @Description dump the database.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {file} binary
func DumpRoute(c echo.Context) error {
	user := c.FormValue("user")
	database := c.FormValue("database")
	fmt.Println("Dumping database " + database + " for user " + user)
	dumpExec := Dump()
	return c.File(dumpExec.File)
}
