package quest

import (
	"context"

	"github.com/SQUASHD/hbit/rpg/rpgdb"
)

type (
	Repository interface {
		CRUDRepository
		UserRepository
		Cleanup() error
	}

	CRUDRepository interface {
		List(ctx context.Context, userId string) ([]rpgdb.Quest, error)
		Create(ctx context.Context, data rpgdb.CreateQuestParams) (rpgdb.Quest, error)
		Read(ctx context.Context, id string) (rpgdb.Quest, error)
		Update(ctx context.Context, data rpgdb.UpdateQuestParams) (rpgdb.Quest, error)
		Delete(ctx context.Context, id string) error
	}

	UserRepository interface {
		GetUserQuests(ctx context.Context, userId string) ([]rpgdb.Quest, error)
		DeleteUserQuests(ctx context.Context, userId string) error
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
	quest, err := s.repo.Read(ctx, id)
	if err != nil {
		return QuestDTO{}, err
	}

	return questToDTO(quest), nil
}

func (s *service) CreateQuest(ctx context.Context, data rpgdb.CreateQuestParams) (QuestDTO, error) {
	quest, err := s.repo.Create(ctx, data)
	if err != nil {
		return QuestDTO{}, err
	}

	return questToDTO(quest), nil
}

func (s *service) UpdateQuest(ctx context.Context, data rpgdb.UpdateQuestParams) (QuestDTO, error) {
	quest, err := s.repo.Update(ctx, data)
	if err != nil {
		return QuestDTO{}, err
	}

	return questToDTO(quest), nil
}

func (s *service) DeleteQuest(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) Cleanup() error {
	return s.repo.Cleanup()
}

func (s *service) DeleteUserQuests(userId string) error {
	return s.repo.DeleteUserQuests(context.Background(), userId)
}
