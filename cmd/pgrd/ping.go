package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/urfave/cli/v2"
)

func Ping(cCtx *cli.Context) error {
	c := NewConfigFromContextOrEnv(cCtx)
	start := time.Now()
	fmt.Println("host: ", c.host)
	r, _ := http.NewRequest("GET", "http://"+c.host+":"+c.port+"/", bytes.NewReader([]byte{}))
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Key", c.key)
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
