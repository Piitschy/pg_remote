package main

import (
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

type Config struct {
	host       string
	port       string
	key        string
	dbUser     *string
	dbPassword *string
	dbName     *string
}

func NewConfigFromContextOrEnv(cCtx *cli.Context) *Config {
	envhost := os.Getenv("PGRD_HOST")
	envport := os.Getenv("PGRD_PORT")
	envkey := os.Getenv("PGRD_KEY")
	envdbuser := os.Getenv("PGRD_LOCAL_DB_USER")
	envdbpassword := os.Getenv("PGRD_LOCAL_DB_PASSWORD")
	envdbname := os.Getenv("PGRD_LOCAL_DB_NAME")

	c := &Config{}
	if cCtx.String("host") != "" {
		c.host = cCtx.String("host")
	} else if envhost != "" {
		c.host = envhost
	} else {
		c.host = "localhost"
	}

	if cCtx.String("port") != "" {
		c.port = cCtx.String("port")
	} else if envport != "" {
		c.port = envport
	} else {
		c.port = "5432"
	}

	if cCtx.String("key") != "" {
		c.key = cCtx.String("key")
	} else if envkey != "" {
		c.key = envkey
	} else {
		c.key = ""
	}

	if cCtx.String("localdb-user") != "" {
		dbUser := cCtx.String("localdb-user")
		c.dbUser = &dbUser
	} else if envdbuser != "" {
		c.dbUser = &envdbuser
	} else {
		c.dbUser = nil
	}

	if cCtx.String("localdb-password") != "" {
		dbPassword := cCtx.String("localdb-password")
		c.dbPassword = &dbPassword
	} else if envdbpassword != "" {
		c.dbPassword = &envdbpassword
	} else {
		c.dbPassword = nil
	}

	if cCtx.String("localdb-name") != "" {
		dbName := cCtx.String("localdb-name")
		c.dbName = &dbName
	} else if envdbname != "" {
		c.dbName = &envdbname
	} else {
		c.dbName = nil
	}

	return c

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
