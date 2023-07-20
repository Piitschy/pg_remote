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
	User     string `json:"user" form:"user" query:"user"`
	Database string `json:"database" form:"database" query:"database"`
}

// Dump DB
// @Router /dump [post]
// @Summary Dump the database.
// @Description dump the database.
// @Tags root
// @Accept */*
// @Produce json
// @Param user formData string true "User"
// @Param database formData string true "Database"
// @Success 200 {file} binary
func DumpRoute(c echo.Context) error {
	fmt.Println("Dumping...")
	req := new(DumpRequest)
	c.Bind(req) //TODO: binding not working
	fmt.Println(req)
	user := req.User
	database := req.Database
	fmt.Println(user, database)
	return c.JSON(http.StatusOK, Response{
		Msg: fmt.Sprintf("Dumping %s database for user %s", database, user),
	})
	//dumpExec := Dump()
	//return c.File(dumpExec.File)
}
