package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	pg "github.com/habx/pg-commands"
)

type path string

func NewPostgres() *pg.Postgres {
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
	return &pg.Postgres{
		Host:     host,
		Port:     port,
		DB:       db,
		Username: user,
		Password: password,
	}
}

func Dump() pg.Result {
	now := time.Now().Format("2006-01-02T15:04:05")
	filename := "dump_" + now + ".tar"
	dump, err := pg.NewDump(DB)
	if err != nil {
		panic(err)
	}
	dump.SetFileName(filename)
	dump.SetupFormat("t")
	dumpExec := dump.Exec(pg.ExecOptions{StreamPrint: false})
	if dumpExec.Error != nil {
		fmt.Println(dumpExec.Error.Err)
		fmt.Println(dumpExec.Output)
	}
	fmt.Println(dumpExec.Output)
	fmt.Println("Dump success")
	return dumpExec
}

func Restore(path string) error {
	restore, err := pg.NewRestore(DB)
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
