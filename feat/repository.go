package feat

import (
	"context"
	"database/sql"

	"github.com/SQUASHD/hbit/feat/featdb"
)

type repository struct {
	queries *featdb.Queries
	db      *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	queries := featdb.New(db)
	return &repository{queries: queries, db: db}
}

func (r *repository) ListFeats(ctx context.Context, userId string) ([]featdb.ListUserFeatsRow, error) {
	return r.queries.ListUserFeats(ctx, userId)
}

func (r *repository) Cleanup() error {
	return r.db.Close()
}
