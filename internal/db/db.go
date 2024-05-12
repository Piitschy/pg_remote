package db

import pg "github.com/habx/pg-commands"

type Database interface {
	Dump(format string) pg.Result
	Restore(path string) error
}
