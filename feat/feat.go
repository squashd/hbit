package feat

import "context"

type (
	Service interface {
		UserFeatService
		Cleanup() error
	}

	UserFeatService interface {
		GetUserFeats(ctx context.Context, userID string) ([]UserFeatsDTO, error)
	}
)
