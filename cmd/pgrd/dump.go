package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Piitschy/postgress-dump-tool/internal/db"
	"github.com/urfave/cli/v2"
)

func Dump(cCtx *cli.Context) error {
	c := NewConfigFromContextOrEnv(cCtx)
	// format := cCtx.String("format")
	filename := cCtx.String("output-file")

	log.Println("host: ", c.host)
	r, _ := http.NewRequest("POST", "http://"+c.host+":"+c.port+"/dump", bytes.NewReader([]byte{}))
	r.Header.Set("Content-Type", "application/json") //TODO: Transmit format
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Key", c.key)

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

	// ext := filepath.Ext(filename)
	// if ext != "" {
	// 	format = ext
	// }

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

func LocalDump(cCtx *cli.Context) error {
	c := NewConfigFromContextOrEnv(cCtx)
	port, err := strconv.Atoi(c.port)
	if err != nil {
		return err
	}
	db := db.NewPostgres(c.host, port, *c.dbName, *c.dbUser, *c.dbPassword)

	filename := cCtx.String("output-file")

	dumpExec, err := db.Dump(filename, "t")
	if err != nil {
		return err
	}

	log.Println(dumpExec.Output)
	return nil
}
