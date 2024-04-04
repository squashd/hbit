package feat

import (
	"context"

	"github.com/SQUASHD/hbit"
)

type (
	Service interface {
		UserFeatService
		Cleanup() error
	}

	EventService interface {
		hbit.Publisher
		hbit.UserDataHandler
	}

	UserFeatService interface {
		GetUserFeats(ctx context.Context, userID string) ([]UserFeatsDTO, error)
	}
)
