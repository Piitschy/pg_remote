package main

type Response struct {
	Msg string `json:"message"`
}

type DumpRequest struct {
	User     string `json:"user"`
	Database string `json:"database"`
}

type Empty struct{}
