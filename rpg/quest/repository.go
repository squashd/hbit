package quest

import (
	"context"
	"database/sql"

	"github.com/SQUASHD/hbit/rpg/rpgdb"
)

type repository struct {
	queries *rpgdb.Queries
	db      *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	queries := rpgdb.New(db)
	return &repository{queries: queries, db: db}
}

// Create implements Repository.
func (r *repository) Create(ctx context.Context, data rpgdb.CreateQuestParams) (rpgdb.Quest, error) {
	panic("unimplemented")
}

// Delete implements Repository.
func (r *repository) Delete(ctx context.Context, id string) error {
	panic("unimplemented")
}

// DeleteUserQuests implements Repository.
func (r *repository) DeleteUserQuests(ctx context.Context, userId string) error {
	panic("unimplemented")
}

// GetUserQuests implements Repository.
func (r *repository) GetUserQuests(ctx context.Context, userId string) ([]rpgdb.Quest, error) {
	panic("unimplemented")
}

// List implements Repository.
func (r *repository) List(ctx context.Context, userId string) ([]rpgdb.Quest, error) {
	panic("unimplemented")
}

// Read implements Repository.
func (r *repository) Read(ctx context.Context, id string) (rpgdb.Quest, error) {
	panic("unimplemented")
}

// Update implements Repository.
func (r *repository) Update(ctx context.Context, data rpgdb.UpdateQuestParams) (rpgdb.Quest, error) {
	panic("unimplemented")
}

func (r *repository) Cleanup() error {
	return r.db.Close()
}
