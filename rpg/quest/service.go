package quest

import (
	"context"
	"database/sql"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/rpg/rpgdb"
)

type (
	service struct {
		db      *sql.DB
		queries *rpgdb.Queries
	}
)

func NewService(
	db *sql.DB,
	queries *rpgdb.Queries,
) QuestService {
	return &service{
		db:      db,
		queries: queries,
	}
}

func (s *service) ListQuests(ctx context.Context, userId string) ([]QuestDTO, error) {
	quests, err := s.queries.ListQuests(ctx)
	if err != nil {
		return nil, err
	}

	return questsToDTOs(quests), nil
}

func (s *service) ReadQuest(ctx context.Context, id string) (QuestDTO, error) {
	return QuestDTO{}, &hbit.Error{Code: hbit.EINTERNAL, Message: "Not implemented"}
}

func (s *service) CreateQuest(ctx context.Context, data any) (QuestDTO, error) {
	return QuestDTO{}, &hbit.Error{Code: hbit.EINTERNAL, Message: "Not implemented"}
}

func (s *service) UpdateQuest(ctx context.Context, data any) (QuestDTO, error) {
	return QuestDTO{}, &hbit.Error{Code: hbit.EINTERNAL, Message: "Not implemented"}

}

func (s *service) DeleteQuest(ctx context.Context, id string) error {
	return &hbit.Error{Code: hbit.EINTERNAL, Message: "Not implemented"}

}

func (s *service) CleanUp() error {
	return s.db.Close()
}

func (s *service) DeleteUserQuests(userId string) error {
	return &hbit.Error{Code: hbit.EINTERNAL, Message: "Not implemented"}
}
