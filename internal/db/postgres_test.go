package db_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/Piitschy/pgrd/internal/db"
	"github.com/Piitschy/pgrd/internal/testhelpers"
)

func TestNewPostgres(t *testing.T) {
	db := db.NewPostgres("localhost", 5432, "postgres", "postgres", "password")
	if db.Host != "localhost" {
		t.Errorf("Expected localhost, got %s", db.Host)
	}
	if db.Port != 5432 {
		t.Errorf("Expected 5432, got %d", db.Port)
	}
	if db.DB != "postgres" {
		t.Errorf("Expected postgres, got %s", db.DB)
	}
	if db.Username != "postgres" {
		t.Errorf("Expected postgres, got %s", db.Username)
	}
	if db.Password != "password" {
		t.Errorf("Expected password, got %s", db.Password)
	}
}

func TestNewPostgresFromConnString(t *testing.T) {
	connString := "postgres://username:password@localhost:32781/database?sslmode=disable"
	db, err := db.NewPostgresFromConnString(connString)
	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}
	if db.Host != "localhost" {
		t.Errorf("Expected localhost, got %s", db.Host)
	}
	if db.Port != 32781 {
		t.Errorf("Expected 32781, got %d", db.Port)
	}
	if db.DB != "database" {
		t.Errorf("Expected database, got %s", db.DB)
	}
	if db.Username != "username" {
		t.Errorf("Expected username, got %s", db.Username)
	}
	if db.Password != "password" {
		t.Errorf("Expected password, got %s", db.Password)
	}
}

func TestPostgresConnection(t *testing.T) {
	ctx := context.Background()
	pgContainer, err := testhelpers.CreatePostgresContainer(ctx)
	if err != nil {
		t.Errorf("error while container creation, got %s", err)
	}

	defer t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Errorf("error while terminating testcontainer: Expected nil, got %s", err)
		}
	})

	db, err := db.NewPostgresFromConnString(pgContainer.ConnectionString)
	if err != nil {
		t.Errorf("Expected nil while parsing connection string, got %s", err)
	}

	err = db.TestConnection()
	if err != nil {
		t.Errorf("Expected nil while testing connection to postgres, got %s", err)
	}
}

func TestDump(t *testing.T) {
	ctx := context.Background()
	pgContainer, err := testhelpers.CreatePostgresContainer(ctx)
	if err != nil {
		t.Errorf("error while container creation, got %s", err)
	}
	defer t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Errorf("error while terminating testcontainer: Expected nil, got %s", err)
		}
	})

	db, err := db.NewPostgresFromConnString(pgContainer.ConnectionString)
	// db, err := db.NewPostgresFromConnString("postgres://pg:password@localhost:5432/database")
	if err != nil {
		t.Errorf("Expected nil while parsing connection string, got %s", err)
	}

	err = db.TestConnection()
	if err != nil {
		t.Errorf("Expected nil while testing connection to postgres, got %s", err)
	}

	// create folder for test dumps
	dir := "./tmp"
	err = os.Mkdir(dir, os.ModePerm)
	if err != nil {
		t.Log(err)
	}
	// Cleanup
	defer func() {
		err = os.RemoveAll(dir)
		if err != nil {
			t.Errorf("Expected nil while removing directory, got %s", err)
		}
	}()

	_, err = db.Dump(dir, "t")
	if err != nil {
		t.Errorf("Expected nil while dumping table, got %s", err)
	}
}

func TestDumpAndRestore(t *testing.T) {
	ctx := context.Background()
	pgContainerBase, err := testhelpers.CreatePostgresContainer(ctx)
	pgContainerTarget, err := testhelpers.CreatePostgresContainer(ctx)
	if err != nil {
		t.Errorf("error while container creation, got %s", err)
	}
	defer t.Cleanup(func() {
		if err := pgContainerBase.Terminate(ctx); err != nil {
			t.Errorf("error while terminating testcontainer: Expected nil, got %s", err)
		}
		if err := pgContainerTarget.Terminate(ctx); err != nil {
			t.Errorf("error while terminating testcontainer: Expected nil, got %s", err)
		}
	})

	dbBase, err := db.NewPostgresFromConnString(pgContainerBase.ConnectionString)
	dbTarget, err := db.NewPostgresFromConnString(pgContainerBase.ConnectionString)
	if err != nil {
		t.Errorf("Expected nil while parsing connection string, got %s", err)
	}
	_, _, err = pgContainerTarget.Exec(ctx, []string{"psql", "-U", dbTarget.Username, "-d", dbTarget.DB, "-c", fmt.Sprintf("DROP DATABASE %s;", dbTarget.DB)})
	if err != nil {
		t.Errorf("error while dropping target database, got %s", err)
	}
	// db, err := db.NewPostgresFromConnString("postgres://pg:password@localhost:5432/database")

	err = dbBase.TestConnection()
	if err != nil {
		t.Errorf("Expected nil while testing connection to postgres, got %s", err)
	}

	// create folder for test dumps
	dir := "./tmp"
	err = os.Mkdir(dir, os.ModePerm)
	if err != nil {
		t.Log(err)
	}
	// Cleanup
	defer func() {
		err = os.RemoveAll(dir)
		if err != nil {
			t.Errorf("Expected nil while removing directory, got %s", err)
		}
	}()

	dump, err := dbBase.Dump(dir, "t")
	if err != nil {
		t.Errorf("Expected nil while dumping table, got %s", err)
	}

	err = dbTarget.Restore(dump.File)
	if err != nil {
		t.Errorf("Expected nil while restoring dump, got %s", err)
	}
}
