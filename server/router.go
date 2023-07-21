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
// @Accept json
// @Produce json
// @Success 200 {file} binary
/*
// @Param data body DumpRequest true "dump params"
//dumpRequest := new(DumpRequest)
//c.Bind(dumpRequest)
*/
func DumpRoute(c echo.Context) error {
	c.Logger().Info("Dumping...")
	dumpExec := Dump()
	c.Logger().Info(dumpExec.File)
	return c.File(dumpExec.File)
}
