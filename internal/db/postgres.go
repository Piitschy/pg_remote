package db

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Piitschy/postgress-dump-tool/internal/utils"
	"github.com/jackc/pgx/v5"

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
	return NewPostgres(host, port, db, user, password)
}

// postgres://username:password@localhost:32781/database?sslmode=disable
func NewPostgresFromConnString(connString string) (*Postgres, error) {
	var err error = nil
	sections := strings.Split(strings.Split(connString, "://")[1], "@")
	username := strings.Split(sections[0], ":")[0]
	password := strings.Split(sections[0], ":")[1]
	host := strings.Split(sections[1], ":")[0]
	portStr := strings.Split(strings.Split(sections[1], ":")[1], "/")[0]
	db := strings.Split(strings.Split(sections[1], "/")[1], "?")[0]

	port, err := strconv.Atoi(portStr)
	return NewPostgres(host, port, db, username, password), err
}

func NewPostgres(host string, port int, db, user, password string) *Postgres {
	postgres := pg.Postgres{
		Host:     host,
		Port:     port,
		DB:       db,
		Username: user,
		Password: password,
	}
	return &Postgres{postgres}
}

func (db *Postgres) GetUrl() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", db.Username, db.Password, db.Host, db.Port, db.DB)
}

func (db *Postgres) TestConnection() error {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, db.GetUrl())
	if err != nil {
		return err
	}
	conn.Close(ctx)
	return nil
}

// Dump creates a dump of the database in the specified path
// format can be "sql" or "tar" (or "p" or "t")
func (db *Postgres) Dump(path string, format string) (pg.Result, error) {
	format, err := utils.FormatFlag(format)
	if err != nil {
		return pg.Result{}, err
	}
	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}

	now := time.Now().Format("2006-01-02T15:04:05")
	filename := filepath.Join(path, "dump_"+now+"."+utils.Ext(format))
	dump, err := pg.NewDump(&db.Postgres)
	if err != nil {
		return pg.Result{}, err
	}
	dump.SetFileName(filename)
	dump.SetupFormat(format)
	dumpExec := dump.Exec(pg.ExecOptions{StreamPrint: false})
	if dumpExec.Error != nil {
		return dumpExec, dumpExec.Error.Err
	}
	return dumpExec, nil
}

func (db *Postgres) Restore(path string) error {
	restore, err := pg.NewRestore(&db.Postgres)
	if err != nil {
		return err
	}

	ext := filepath.Ext(path)
	f, _ := utils.FormatFlag(ext)
	if f == "c" {
		restore.Options = append(restore.Options, "-c", "-U", db.Username, "-d", db.DB)
	} else {
		restore.Options = append(restore.Options, "-x", "-F"+f)
		restore.Options = append(restore.Options, "-d", db.DB)
	}

	if restore.Role == "" && restore.Username != "" {
		restore.Role = db.Username
	}

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
