package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Piitschy/postgress-dump-tool/internal/db"
	"github.com/urfave/cli/v2"
)

func Restore(cCtx *cli.Context) error {
	c := NewConfigFromContextOrEnv(cCtx)
	filename := cCtx.String("input-file")

	log.Println("host: ", c.host)

	// format := filepath.Ext(filename)

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", filepath.Base(file.Name()))
	io.Copy(part, file)
	writer.Close()

	r, _ := http.NewRequest("POST", "http://"+c.host+":"+c.port+"/restore", body)
	r.Header.Add("Content-Type", writer.FormDataContentType()) //TODO: Transmit format
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Key", c.key)

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return err
	}

	log.Println(resp.Status)
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error: %s", resp.Status)
	}
	log.Println("Restore completed")
	return nil
}

func LocalRestore(cCtx *cli.Context) error {
	c := NewConfigFromContextOrEnv(cCtx)
	port, err := strconv.Atoi(c.port)
	if err != nil {
		return err
	}
	db := db.NewPostgres(c.host, port, *c.dbName, *c.dbUser, *c.dbPassword)

	filename := cCtx.String("input-file")
	log.Println("host: ", c.host)

	err = db.Restore(filename)
	if err != nil {
		return err
	}
	log.Println("Restore completed")
	return nil
}
