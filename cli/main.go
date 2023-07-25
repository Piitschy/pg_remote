package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

var host string
var port string
var key string
var format string
var filename string

func main() {

	app := &cli.App{
		Name:  "pg_remote",
		Usage: "pg_remote is a tool for managing remote postgresql databases",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "host",
				Aliases:     []string{"H"},
				Value:       "localhost",
				Usage:       "Host to connect to the database",
				Destination: &host,
			},
			&cli.StringFlag{
				Name:        "port",
				Aliases:     []string{"P", "p"},
				Value:       "5432",
				Usage:       "Port to connect to the database",
				Destination: &port,
			},

			&cli.StringFlag{
				Name:        "key",
				Aliases:     []string{"k"},
				Value:       "",
				Usage:       "Path the key from the environment variable KEY",
				Destination: &key,
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "dump",
				Aliases: []string{"d"},
				Usage:   "dump a remote database",
				Action:  Dump,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "format",
						Aliases:     []string{"f"},
						Value:       "tar",
						Usage:       "`FORMAT` of the dump file ('tar' | 'pain')",
						Destination: &format,
					},
					&cli.StringFlag{
						Name:        "output-file",
						Aliases:     []string{"o"},
						Value:       "dump",
						Usage:       "output file `NAME`",
						Destination: &filename,
					},
				},
			},
			{
				Name:      "restore",
				Aliases:   []string{"r"},
				Usage:     "restore a remote database",
				ArgsUsage: "[source file]",
				Action:    Restore,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "input-file",
						Aliases:     []string{"i"},
						Value:       "dump",
						Usage:       "input file `NAME`",
						Destination: &filename,
					},
				},
			},
			{
				Name:   "ping",
				Usage:  "Test connection to server",
				Action: Ping,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func Dump(cCtx *cli.Context) error {

	fmt.Println("host: ", host)
	r, _ := http.NewRequest("POST", "http://"+host+":"+port+"/dump", bytes.NewReader([]byte{}))
	r.Header.Set("Content-Type", "application/json") //TODO: Transmit format
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Key", key)

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body := bufio.NewReader(resp.Body)

	fmt.Println("response Status:", resp.Status)

	writer := bufio.NewWriter(os.Stdout)
	if filename != "" {
		file, err := os.Create(extFilename(filename))
		if err != nil {
			return err
		}
		defer file.Close()
		writer = bufio.NewWriter(file)
	}
	defer writer.Flush()

	io.Copy(writer, body)

	return nil
}

func Restore(cCtx *cli.Context) error {

	fmt.Println("host: ", host)

	format = detectFormat(filename)

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

	r, _ := http.NewRequest("POST", "http://"+host+":"+port+"/restore", body)
	r.Header.Add("Content-Type", writer.FormDataContentType()) //TODO: Transmit format
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Key", key)

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return err
	}

	fmt.Println("response Status:", resp.Status)

	return nil
}

func Ping(cCtx *cli.Context) error {
	start := time.Now()
	fmt.Println("host: ", host)
	r, _ := http.NewRequest("GET", "http://"+host+":"+port+"/", bytes.NewReader([]byte{}))
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Key", key)
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("response time:", time.Now().Sub(start))
	fmt.Println("response Status:", resp.Status)
	return nil
}

func extFilename(filename string) string {
	var ext string
	if format == "tar" {
		ext = "tar"
	} else {
		ext = "sql"
	}
	if filename == "" {
		filename = "dump." + ext
	}
	if strings.Split(filename, ".")[len(strings.Split(filename, "."))-1] != ext {
		filename = filename + "." + ext
	}
	return filename
}

func detectFormat(filename string) string {
	ext := strings.Split(filename, ".")[len(strings.Split(filename, "."))-1]
	if ext == "tar" {
		return "tar"
	}
	if ext == "sql" || ext == "txt" {
		return "plain"
	}
	return "custom"
}
