package db

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	pg "github.com/habx/pg-commands"
)

// Postgres is a wrapper around pg.postgres
// ENV variables:
// DB_HOST
// DB_DATABASE
// DB_USER
// DB_PASSWORD
// DB_PORT

type path string

type Postgres struct {
	pg.Postgres
}

func NewPostgresFromEnv() *Postgres {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	db := os.Getenv("DB_DATABASE")
	if db == "" {
		db = "postgres"
	}
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	port_s := os.Getenv("DB_PORT")
	if port_s == "" {
		port_s = "5432"
	}
	port, err := strconv.Atoi(port_s)
	if err != nil {
		panic(err)
	}
	return NewPostgres(host, db, user, password, port)
}

func NewPostgres(host, db, user, password string, port int) *Postgres {
	postgres := pg.Postgres{
		Host:     host,
		Port:     port,
		DB:       db,
		Username: user,
		Password: password,
	}
	return &Postgres{postgres}
}

func (db *Postgres) Dump(format string) pg.Result {
	if format == "" {
		format = "t"
	}
	if format != "t" && format != "p" {
		log.Fatal("Format must be t or p")
	}
	now := time.Now().Format("2006-01-02T15:04:05")
	filename := "dump_" + now + "." + ext(format)
	dump, err := pg.NewDump(&db.Postgres)
	if err != nil {
		panic(err)
	}
	dump.SetFileName(filename)
	dump.SetupFormat(format)
	dumpExec := dump.Exec(pg.ExecOptions{StreamPrint: false})
	if dumpExec.Error != nil {
		fmt.Println(dumpExec.Error.Err)
		fmt.Println(dumpExec.Output)
	}
	fmt.Println(dumpExec.Output)
	fmt.Println("Dump success")
	return dumpExec
}

func (db *Postgres) Restore(path string) error {
	restore, err := pg.NewRestore(&db.Postgres)
	if err != nil {
		panic(err)
	}

	restore.Options = append(restore.Options, "-Ft")

	restoreExec := restore.Exec(path, pg.ExecOptions{StreamPrint: false})
	if restoreExec.Error != nil {
		fmt.Println(restoreExec.Error.Err)
		fmt.Println(restoreExec.Output)
		return restoreExec.Error.Err
	}
	fmt.Println("Restore success")
	fmt.Println(restoreExec.Output)
	return nil
}

func ext(format string) string {
	if format == "p" {
		return "sql"
	}
	return "tar"
}
