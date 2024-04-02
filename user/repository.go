package user

import (
	"context"
	"database/sql"

	"github.com/SQUASHD/hbit/user/userdb"
)

type repository struct {
	queries *userdb.Queries
	db      *sql.DB
}

func NewReposiory(db *sql.DB) Repository {
	queries := userdb.New(db)
	return &repository{
		queries: queries,
		db:      db,
	}
}

// CreateSettings implements Repository.
func (r *repository) CreateSettings(ctx context.Context, data userdb.CreateUserSettingsParams) (userdb.UserSetting, error) {
	panic("unimplemented")
}

// ReadSettings implements Repository.
func (r *repository) ReadSettings(ctx context.Context, userId string) (userdb.UserSetting, error) {
	panic("unimplemented")
}

// UpdateSettings implements Repository.
func (r *repository) UpdateSettings(ctx context.Context, data userdb.UpdateUserSettingsParams) (userdb.UserSetting, error) {
	panic("unimplemented")
}

// DeleteSettings implements Repository.
func (r *repository) DeleteSettings(ctx context.Context, userId string) error {
	panic("unimplemented")
}

func (r *repository) Cleanup() error {
	return r.db.Close()
}
