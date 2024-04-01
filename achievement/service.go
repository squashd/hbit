package achievement

import (
	"context"
)

type (
	Repository interface {
		ListUserAchievements(ctx context.Context, userId string) (UserAchievements, error)
		CreateTaskDoneEvent(ctx context.Context, userId string) error
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

// TaskDone implements Service.
func (s *service) TaskDone(userId string) error {
	return s.repo.CreateTaskDoneEvent(context.Background(), userId)
}

// UserDeleted implements Service.
func (s *service) UserDeleted(userId string) error {
	panic("unimplemented")
}

func (s *service) GetUserAchievements(ctx context.Context, userId string) (UserAchievements, error) {
	return s.repo.ListUserAchievements(ctx, userId)
}
