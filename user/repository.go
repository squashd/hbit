package user

import (
	"context"
	"database/sql"

	"github.com/SQUASHD/hbit/user/database"
)

type repository struct {
	queries *database.Queries
}

func NewReposiory(db *sql.DB) Repository {
	queries := database.New(db)
	return &repository{
		queries: queries,
	}
}

// CreateSettings implements Repository.
func (r *repository) CreateSettings(ctx context.Context, data database.CreateUserSettingsParams) (database.UserSetting, error) {
	panic("unimplemented")
}

// ReadSettings implements Repository.
func (r *repository) ReadSettings(ctx context.Context, userId string) (database.UserSetting, error) {
	panic("unimplemented")
}

// UpdateSettings implements Repository.
func (r *repository) UpdateSettings(ctx context.Context, data database.UpdateUserSettingsParams) (database.UserSetting, error) {
	panic("unimplemented")
}

// DeleteSettings implements Repository.
func (r *repository) DeleteSettings(ctx context.Context, userId string) error {
	panic("unimplemented")
}
