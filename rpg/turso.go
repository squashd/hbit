package rpg

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/pressly/goose/v3"
)

func NewDatabase() (*sql.DB, error) {
	url := os.Getenv("RPG_DB_URL")
	token := os.Getenv("RPG_DB_TOKEN")
	connectionStr := fmt.Sprintf("%s?authToken=%s", url, token)
	db, err := sql.Open("libsql", connectionStr)
	if err != nil {
		return nil, err
	}
	goose.SetBaseFS(Migrations)
	if err := goose.SetDialect("sqlite"); err != nil {
		return nil, err
	}
	if err := goose.Up(db, "schemas"); err != nil {
		return nil, err
	}
	return sql.Open("libsql", connectionStr)
}

func DatabaseDown() error {
	url := os.Getenv("RPG_DB_URL")
	token := os.Getenv("RPG_DB_TOKEN")
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
