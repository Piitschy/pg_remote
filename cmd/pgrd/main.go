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
		Name:  "pgrd",
		Usage: "pgrd (postgres remote dump) is a tool for managing remote postgres based databases.\nIt allows you to dump and restore databases remotely.\n\nFor example:\npgrd -H localhost -p 5432 -k yourservice dump -o dump.tar",
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
				Aliases:     []string{"k", "K"},
				Value:       "",
				Usage:       "Path the key from the environment variable KEY",
				Destination: &key,
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "dump",
				Aliases: []string{"d"},
				Usage:   "dump a remote database that has the server connected.\n\nFor example:\npgrd -H localhost -p 5432 -k yourservice dump -o dump.tar",
				Action:  Dump,
				Flags: []cli.Flag{
					// &cli.StringFlag{
					// 	Name:        "format",
					// 	Aliases:     []string{"f"},
					// 	Value:       "tar",
					// 	Usage:       "`FORMAT` of the dump file ('tar'/'t' | 'plain'/'p' | 'custom'/'c')",
					// 	Destination: &format,
					// },
					&cli.StringFlag{
						Name:        "output-file",
						Aliases:     []string{"o"},
						Value:       "dump.tar",
						Usage:       "output file `NAME`",
						Destination: &filename,
					},
				},
			},
			{
				Name:      "restore",
				Aliases:   []string{"r"},
				Usage:     "restore a remote database.\n\nFor example:\npgrd -H localhost -p 5432 -k yourservice restore -i dump.tar",
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

	envs()

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func envs() {
	host = os.Getenv("PGRD_HOST")
	port = os.Getenv("PGRD_PORT")
	key = os.Getenv("PGRD_KEY")

	filename = os.Getenv("PGRD_FILENAME")
	format = os.Getenv("PGRD_FORMAT")
}

func Dump(cCtx *cli.Context) error {

	log.Println("host: ", host)
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

	log.Println("response Status:", resp.Status)
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error: %s", resp.Status)
	}

	ext := filepath.Ext(filename)
	if ext != "" {
		format = ext
	}

	fullFilename := genFilenameFromFormat(filename)
	log.Println("Dump loaded")
	log.Println("writing file:", fullFilename, "...")
	writer := bufio.NewWriter(os.Stdout)
	if filename != "" {
		file, err := os.Create(fullFilename)
		if err != nil {
			return err
		}
		defer file.Close()
		writer = bufio.NewWriter(file)
	}
	defer writer.Flush()

	io.Copy(writer, body)

	log.Println("write completed")

	return nil
}

func Restore(cCtx *cli.Context) error {

	log.Println("host: ", host)

	format = filepath.Ext(filename)

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

	log.Println(resp.Status)
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error: %s", resp.Status)
	}
	log.Println("Restore completed")
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

func genFilenameFromFormat(filename string) string {
	var ext string
	ext = "tar"
	// if format == "tar" || format == "t" || format == "" {
	// 	ext = "tar"
	// } else if format == "plain" || format == "p" || format == "sql" {
	// 	ext = "sql"
	// } else {
	// 	ext = "custom"
	// }
	if filename == "" {
		filename = "dump." + ext
	}
	if strings.Split(filename, ".")[len(strings.Split(filename, "."))-1] != ext {
		filename = filename + "." + ext
	}
	return filename
}
