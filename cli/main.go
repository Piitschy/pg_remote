package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/urfave/cli/v2"
)

var user string
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
						Name:        "User",
						Aliases:     []string{"U"},
						Value:       "postgres",
						Usage:       "User to connect to the database",
						Destination: &user,
					},
					&cli.StringFlag{
						Name:        "filename",
						Aliases:     []string{"f"},
						Value:       "",
						Usage:       "File to save the dump",
						Destination: &filename,
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func Dump(cCtx *cli.Context) error {
	database := cCtx.Args().First()
	if database == "" {
		return fmt.Errorf("database name is required")
	}
	fmt.Println("host: ", host)
	fmt.Println("user: ", user)
	body := []byte(`{
		"database": "` + database + `",
		"user": "` + user + `",
	}`)
	r, _ := http.NewRequest("POST", "http://"+host+":"+port+"/dump", bytes.NewReader(body))
	fmt.Println("request", r)
	r.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("response", resp)
	fmt.Println("response Status:", resp.Status)

	out := os.Stdout
	if filename != "" {
		out, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer out.Close()
	}
	io.Copy(out, resp.Body)

	return nil
}
