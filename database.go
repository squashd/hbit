package hbit

import (
	"database/sql"
	"io/fs"

	"github.com/pressly/goose/v3"
)

type DbDriver string

const (
	DbDriverLibsql DbDriver = "libsql"
)

type NewDbParams struct {
	ConnectionStr string
	Driver        DbDriver
}

// NewDatabase is a wrapper around sql.Open
func NewDatabase(
	params NewDbParams,
) (*sql.DB, error) {
	db, err := sql.Open(string(params.Driver), params.ConnectionStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}

type MigrationData struct {
	FS      fs.FS
	Dialect string
	Dir     string
}

// DBMigrateUp is a wrapper around goose.Up
func DBMigrateUp(
	db *sql.DB,
	migrations MigrationData,
) error {
	goose.SetBaseFS(migrations.FS)
	if err := goose.SetDialect(migrations.Dialect); err != nil {
		return err
	}
	if err := goose.Up(db, migrations.Dir); err != nil {
		return err
	}
	return nil
}

// DBMigratedDown is a wrapper around goose.DownTo
func DBMigratedDownTo(
	db *sql.DB,
	target int64,
	migrations MigrationData,
) error {
	goose.SetBaseFS(migrations.FS)
	if err := goose.SetDialect(migrations.Dialect); err != nil {
		return err
	}
	if err := goose.DownTo(db, migrations.Dir, target); err != nil {
		return err
	}
	return nil
}
