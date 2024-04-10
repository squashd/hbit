package user

import (
	"context"
	"database/sql"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/user/userdb"
)

type (
	Service interface {
		UserSettingsService
		InternalUserSettingsService
		hbit.UserDataHandler
		Cleanup() error
	}

	service struct {
		db      *sql.DB
		queries *userdb.Queries
	}
)

func NewService(
	db *sql.DB,
	queries *userdb.Queries,
) Service {
	return &service{
		db:      db,
		queries: queries,
	}
}

func (s *service) DeleteData(ctx context.Context, userId string) error {
	return s.queries.DeleteUserData(ctx, userId)
}

func (s *service) Cleanup() error {
	return s.db.Close()
}
