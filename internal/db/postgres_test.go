package db_test

import (
	"context"
	"testing"

	"github.com/Piitschy/postgress-dump-tool/internal/db"
	"github.com/testcontainers/testcontainers-go"
	testpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
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
	data := map[string]string{
		"db":       "postgres",
		"username": "postgres",
		"password": "password",
	}

	ctx := context.Background()

	pgContainer, err := testpostgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15.3-alpine"),
		testpostgres.WithDatabase(data["db"]),
		testpostgres.WithUsername(data["username"]),
		testpostgres.WithPassword(data["password"]),
		// testcontainers.WithWaitStrategy(
		// 	wait.ForLog("database system is ready to accept connections").
		// 		WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)

	if err != nil {
		t.Errorf("error while creating testconatiner: Expected nil, got %s", err)
	}

	defer t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Errorf("error while terminating testcontainer: Expected nil, got %s", err)
		}
	})

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Errorf("Expected nil, got %s", err)

	}

	db, err := db.NewPostgresFromConnString(connStr)
	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}

	err = db.TestConnection()
	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}
}
