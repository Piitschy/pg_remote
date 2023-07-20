package main

import (
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
// @Success 200 {object} Response
func DumpRoute(c echo.Context) error {
	//Dump(db)
	return c.JSON(http.StatusOK, Response{
		Msg: "Dump success",
	})
}
