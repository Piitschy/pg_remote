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

type DumpRequest struct {
	User     string `json:"user"`
	Database string `json:"database"`
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
	fmt.Println("Dumping...")
	dumpExec := Dump()
	fmt.Println("\n", dumpExec.Output, "\n", dumpExec.File, "\n")
	return c.File(dumpExec.File)
}
