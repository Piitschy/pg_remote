package db

import pg "github.com/habx/pg-commands"

type Database interface {
	Dump(path, format string) (pg.Result, error)
	Restore(path string) error
}
