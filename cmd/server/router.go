package main

import (
	"io"
	"net/http"
	"os"

	db "github.com/Piitschy/pgrd/internal/db"
	"github.com/labstack/echo/v4"
)

type Response struct {
	Msg string `json:"message"`
}

type DumpRequest struct {
	User     string `json:"user"`
	Database string `json:"database"`
}

type Empty struct{}

// HealthCheck godoc
// @Router / [get]
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags root
// @Accept */*
// @Produce json
// @Param Key header string true "Key from environment"
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
// @Param Key header string true "Key from environment"
// @Param data body Empty true "future: dump params"
/* //TODO: Add body params
// @Param data body DumpRequest true "dump params"
//dumpRequest := new(DumpRequest)
//c.Bind(dumpRequest)
*/
func DumpRoute(db db.Database, path string) func(c echo.Context) error {
	return func(c echo.Context) error {
		c.Logger().Info("Dumping...")
		dumpExec, err := db.Dump(path, "tar")
		if err != nil {
			return err
		}
		c.Logger().Info(dumpExec.Output)
		c.Logger().Info("Dump success")
		c.Logger().Info(dumpExec.File)
		return c.File(dumpExec.File)
	}
}

// Restore DB
// @Router /restore [post]
// @Summary Restore the database.
// @Description Restore the database.
// @Tags root
// @Accept json
// @Produce json
// @Success 200 {file} binary
// @Param Key header string true "Key from environment"
func RestoreRoute(db db.Database) func(c echo.Context) error {
	return func(c echo.Context) error {
		c.Logger().Info("Restoring...")
		file, err := c.FormFile("file")
		if err != nil {
			return err
		}
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		dst, err := os.Create(file.Filename)
		if err != nil {
			return err
		}
		defer dst.Close()

		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			return err
		}

		err = db.Restore(file.Filename)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Response{
				Msg: "Error restoring",
			})
		}
		return c.JSON(http.StatusOK, Response{
			Msg: "Restored successfully",
		})
	}
}
