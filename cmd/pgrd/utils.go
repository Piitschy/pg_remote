package main

import (
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
	c := &Config{
		host:       cCtx.String("host"),
		port:       cCtx.String("port"),
		key:        cCtx.String("key"),
		dbUser:     nil,
		dbPassword: nil,
		dbName:     nil,
	}
	if cCtx.String("localdb-user") != "" {
		u := cCtx.String("localdb-user")
		c.dbUser = &u
	}
	if cCtx.String("localdb-password") != "" {
		pw := cCtx.String("localdb-password")
		c.dbPassword = &pw
	}
	if cCtx.String("localdb-name") != "" {
		n := cCtx.String("localdb-name")
		c.dbName = &n
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
