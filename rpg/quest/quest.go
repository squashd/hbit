package quest

import (
	"context"
	"database/sql"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/rpg/rpgdb"
)

type (
	// Top level service for instantiation
	Service interface {
		UserQuestService
		InternalUserQuestUtils
		hbit.UserDataHandler
		CleanUp() error
	}

	service struct {
		db      *sql.DB
		queries *rpgdb.Queries
	}
)

func NewService(
	db *sql.DB,
	queries *rpgdb.Queries,
) Service {
	return &service{
		db:      db,
		queries: queries,
	}
}

func (s *service) DeleteData(ctx context.Context, userId string) error {
	return s.queries.DeleteUserQuestData(ctx, userId)
}

func (s *service) CleanUp() error {
	return s.db.Close()
}
