package achievement

import (
	"context"
	"database/sql"

	"github.com/SQUASHD/hbit/achievement/database"
)

type repository struct {
	queries *database.Queries
}

// ListAchievements implements Repository.
func (r *repository) ListAchievements(ctx context.Context) ([]database.Achievement, error) {
	panic("unimplemented")
}

// ListUserAchievements implements Repository.
func (r *repository) ListUserAchievements(ctx context.Context, userID string) ([]database.UserAchievement, error) {
	panic("unimplemented")
}

// ReadAchievement implements Repository.
func (r *repository) ReadAchievement(ctx context.Context, id string) (database.Achievement, error) {
	panic("unimplemented")
}

// ReadUserAchievement implements Repository.
func (r *repository) ReadUserAchievement(ctx context.Context, userID string, achievementID string) (database.UserAchievement, error) {
	panic("unimplemented")
}

func NewRepository(db *sql.DB) Repository {
	queries := database.New(db)
	return &repository{queries}
}
