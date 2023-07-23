package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

var host string
var port string
var key string
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
				Name:      "dump",
				Aliases:   []string{"d"},
				Usage:     "dump a remote database",
				ArgsUsage: "[database]",
				Action:    Dump,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "filename",
						Aliases:     []string{"f"},
						Value:       "",
						Usage:       "File to save the dump",
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
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Key", key)

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("response", resp)
	fmt.Println("response Status:", resp.Status)

	writer := bufio.NewWriter(os.Stdout)
	if filename != "" {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close()
		writer = bufio.NewWriter(file)
	}
	defer writer.Flush()

	io.Copy(writer, resp.Body)

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
