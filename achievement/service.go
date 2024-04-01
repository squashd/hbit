package achievement

import (
	"context"
	"encoding/json"
)

type (
	Repository interface {
		ListUserAchievements(ctx context.Context, userId string) (UserAchievements, error)
	}

	service struct {
		repo Repository
	}
)

func NewService(r Repository) Service {
	return &service{
		repo: r,
	}
}

func (s *service) GetUserAchievements(ctx context.Context, userId string) (UserAchievements, error) {
	return s.repo.ListUserAchievements(ctx, userId)
}

func (s *service) TaskComplet(msg json.RawMessage) (string, error) {
	panic("unimplemented")
}

func (s *service) UserDeleted(msg json.RawMessage) (string, error) {
	panic("unimplemented")
}
