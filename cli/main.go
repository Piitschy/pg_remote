package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/urfave/cli/v2"
)

var user string
var host string
var port string
var key string

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
	fmt.Println("user: ", user)
	body := []byte(`{
		"database": "` + database + `",
		"user": "` + user + `",
	}`)
	r, _ := http.NewRequest("POST", "http://"+host+":"+port+"/dump", bytes.NewBuffer(body))
	r.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
	return nil
}
