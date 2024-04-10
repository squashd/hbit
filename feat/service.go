package feat

import (
	"context"
)

type (
	UserFeatService interface {
		GetUserFeats(ctx context.Context, userID string) ([]UserFeatsDTO, error)
	}
	EventService interface {
	}
)

func (s *service) GetUserFeats(ctx context.Context, userId string) ([]UserFeatsDTO, error) {
	feats, err := s.queries.ListUserFeats(ctx, userId)
	if err != nil {
		return []UserFeatsDTO{}, err
	}

	dtos := toDTOs(feats)
	return dtos, nil
}
