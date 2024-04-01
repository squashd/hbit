package quest

import (
	"context"
	"database/sql"

	quest "github.com/SQUASHD/hbit/quest/database"
)

type repository struct {
	queries *quest.Queries
}

func NewRepository(db *sql.DB) Repository {
	queries := quest.New(db)
	return &repository{queries: queries}
}

// Create implements Repository.
func (r *repository) Create(ctx context.Context, data quest.CreateQuestParams) (quest.Quest, error) {
	panic("unimplemented")
}

// Delete implements Repository.
func (r *repository) Delete(ctx context.Context, id string) error {
	panic("unimplemented")
}

// Get implements Repository.
func (r *repository) Get(ctx context.Context, id string) (quest.Quest, error) {
	panic("unimplemented")
}

// List implements Repository.
func (r *repository) List(ctx context.Context, userId string) ([]quest.Quest, error) {
	panic("unimplemented")
}

// Update implements Repository.
func (r *repository) Update(ctx context.Context, data quest.UpdateQuestParams) (quest.Quest, error) {
	panic("unimplemented")
}
