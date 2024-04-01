package quest

import (
	"context"

	"github.com/SQUASHD/hbit/quest/database"
)

type (
	Repository interface {
		List(ctx context.Context, userId string) ([]database.Quest, error)
		Get(ctx context.Context, id string) (database.Quest, error)
		Create(ctx context.Context, data database.CreateQuestParams) (database.Quest, error)
		Update(ctx context.Context, data database.UpdateQuestParams) (database.Quest, error)
		Delete(ctx context.Context, id string) error
	}

	service struct {
		repo Repository
	}
)

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) ListQuests(ctx context.Context, userId string) ([]QuestDTO, error) {
	quests, err := s.repo.List(ctx, userId)
	if err != nil {
		return nil, err
	}

	return questsToDTOs(quests), nil
}

func (s *service) ReadQuest(ctx context.Context, id string) (QuestDTO, error) {
	quest, err := s.repo.Get(ctx, id)
	if err != nil {
		return QuestDTO{}, err
	}

	return questToDTO(quest), nil
}

func (s *service) CreateQuest(ctx context.Context, data database.CreateQuestParams) (QuestDTO, error) {
	quest, err := s.repo.Create(ctx, data)
	if err != nil {
		return QuestDTO{}, err
	}

	return questToDTO(quest), nil
}

func (s *service) UpdateQuest(ctx context.Context, data database.UpdateQuestParams) (QuestDTO, error) {
	quest, err := s.repo.Update(ctx, data)
	if err != nil {
		return QuestDTO{}, err
	}

	return questToDTO(quest), nil
}

func (s *service) DeleteQuest(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
