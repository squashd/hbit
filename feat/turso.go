package feat

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/pressly/goose/v3"
)

func NewDatabase() (*sql.DB, error) {
	url := os.Getenv("ACH_DB_URL")
	token := os.Getenv("ACH_DB_TOKEN")
	connectionStr := fmt.Sprintf("%s?authToken=%s", url, token)
	db, err := sql.Open("libsql", connectionStr)
	if err != nil {
		log.Fatalf("cannot open database: %s", err)
	}
	goose.SetBaseFS(Migrations)
	if err := goose.SetDialect("sqlite"); err != nil {
		log.Fatalf("cannot set dialect: %s", err)
	}
	if err := goose.Up(db, "schemas"); err != nil {
		log.Fatalf("cannot run migrations: %s", err)
	}
	return sql.Open("libsql", connectionStr)
}

func DatabaseDown() error {
	url := os.Getenv("ACH_DB_URL")
	token := os.Getenv("ACH_DB_TOKEN")
	connectionStr := fmt.Sprintf("%s?authToken=%s", url, token)
	db, err := sql.Open("libsql", connectionStr)
	if err != nil {
		return err
	}
	goose.SetBaseFS(Migrations)
	if err := goose.SetDialect("sqlite"); err != nil {
		return err
	}
	if err := goose.DownTo(db, "schemas", 0); err != nil {
		return err
	}
	return nil
}
