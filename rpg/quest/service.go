package quest

import (
	"context"
	"database/sql"

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
	quest, err := s.queries.ReadQuest(ctx, id)
	if err != nil {
		return QuestDTO{}, err
	}

	return questToDTO(quest), nil
}

func (s *service) CreateQuest(ctx context.Context, data rpgdb.CreateQuestParams) (QuestDTO, error) {
	quest, err := s.queries.CreateQuest(ctx, data)
	if err != nil {
		return QuestDTO{}, err
	}

	return questToDTO(quest), nil
}

func (s *service) UpdateQuest(ctx context.Context, data rpgdb.UpdateQuestParams) (QuestDTO, error) {
	quest, err := s.queries.UpdateQuest(ctx, data)
	if err != nil {
		return QuestDTO{}, err
	}

	return questToDTO(quest), nil
}

func (s *service) DeleteQuest(ctx context.Context, id string) error {
	return s.queries.DeleteQuest(ctx, id)
}

func (s *service) CleanUp() error {
	return s.db.Close()
}

func (s *service) DeleteUserQuests(userId string) error {
	panic("not implemented")
}
