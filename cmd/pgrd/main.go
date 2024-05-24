package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "pgrd",
		Usage: "pgrd (postgres remote dump) is a tool for managing remote postgres based databases.\nIt allows you to dump and restore databases remotely.\n\nFor example:\npgrd -H localhost -p 5432 -k yourservice dump -o dump.tar",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "host",
				Aliases: []string{"H"},
				Value:   "localhost",
				Usage:   "Host to connect to the database",
				EnvVars: []string{"PGRD_HOST"},
			},
			&cli.StringFlag{
				Name:    "port",
				Aliases: []string{"P", "p"},
				Value:   "5432",
				Usage:   "Port to connect to the database",
				EnvVars: []string{"PGRD_PORT"},
			},

			&cli.StringFlag{
				Name:    "key",
				Aliases: []string{"k", "K"},
				Value:   "",
				Usage:   "Path the key from the environment variable KEY",
				EnvVars: []string{"PGRD_KEY"},
			},

			&cli.StringFlag{
				Name:    "localdb-user",
				Aliases: []string{"u"},
				Value:   "",
				Usage:   "Local database user",
				EnvVars: []string{"PGRD_LOCAL_DB_USER"},
			},
			&cli.StringFlag{
				Name:    "localdb-password",
				Aliases: []string{"pw"},
				Value:   "",
				Usage:   "Local database password",
				EnvVars: []string{"PGRD_LOCAL_DB_PASSWORD"},
			},
			&cli.StringFlag{
				Name:    "localdb-name",
				Aliases: []string{"db"},
				Value:   "",
				Usage:   "Local database name",
				EnvVars: []string{"PGRD_LOCAL_DB_NAME"},
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "remotedump",
				Aliases: []string{"dump", "rd", "rdump"},
				Usage:   "dump a remote database that has the server connected.\n\nFor example:\npgrd --host localhost -p 5432 -k yourservice dump -o dump.tar\n",
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
						Name:      "output-file",
						Aliases:   []string{"o"},
						Value:     "dump.tar",
						Usage:     "output file `NAME`",
						TakesFile: true,
					},
				},
			},
			{
				Name:    "remoterestore",
				Aliases: []string{"rrestore", "restore", "rr"},
				Usage:   "restore a remote database.\n\nFor example:\npgrd --host localhost -p 5432 -k yourservice restore -i dump.tar\n",
				// ArgsUsage: "[source file]",
				Action: Restore,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:      "input-file",
						Aliases:   []string{"i"},
						Value:     "dump",
						Usage:     "input file `NAME`",
						TakesFile: true,
					},
				},
			},
			{
				Name:    "localdump",
				Aliases: []string{"ldump", "ld"},
				Usage:   "dump a local database.\nlocaldb-user, localdb-password and localdb-name are required.\n\nFor example:\npgrd -u user --pw password --db dbname -p 5432 localdump -o dump.tar\n",
				Action:  LocalDump,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:      "output-file",
						Aliases:   []string{"o"},
						Value:     "dump.tar",
						Usage:     "output file `NAME`",
						TakesFile: true,
					},
				},
			},
			{
				Name:    "localrestore",
				Aliases: []string{"lrestore", "lr"},
				Usage:   "restore a local database.\nlocaldb-user, localdb-password and localdb-name are required.\n\nFor example:\npgrd -u user --pw password --db dbname -p 5432 localrestore -i dump.tar\n",
				Action:  LocalRestore,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:      "input-file",
						Aliases:   []string{"i"},
						Value:     "dump",
						Usage:     "input file `NAME`",
						TakesFile: true,
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
