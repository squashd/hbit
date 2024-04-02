package feat

import (
	"context"

	"github.com/SQUASHD/hbit/feat/featdb"
)

type (
	Repository interface {
		UserRepository
		Cleanup() error
	}

	UserRepository interface {
		ListFeats(ctx context.Context, userId string) ([]featdb.ListUserFeatsRow, error)
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

func (s *service) GetUserFeats(ctx context.Context, userID string) ([]UserFeatsDTO, error) {
	feats, err := s.repo.ListFeats(ctx, userID)
	if err != nil {
		return []UserFeatsDTO{}, err
	}

	dtos := toDTOs(feats)
	return dtos, nil
}

func (s *service) Cleanup() error {
	return s.repo.Cleanup()
}
