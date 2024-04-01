package achievement

import (
	"context"
	"database/sql"

	"github.com/SQUASHD/hbit/achievement/database"
)

type repository struct {
	queries *database.Queries
}

// CreateTaskDoneEvent implements Repository.
func (r *repository) CreateTaskDoneEvent(ctx context.Context, userId string) error {
	panic("unimplemented")
}

// ListUserAchievements implements Repository.
func (r *repository) ListUserAchievements(ctx context.Context, userId string) ([]database.UserAchievement, error) {
	panic("unimplemented")
}

func NewRepository(db *sql.DB) Repository {
	queries := database.New(db)
	return &repository{queries}
}
