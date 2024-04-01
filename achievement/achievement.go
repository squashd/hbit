package achievement

import (
	"context"
)

type (
	Service interface {
		GetUserAchievements(ctx context.Context, userID string) (UserAchievements, error)
		TaskDone(userId string) error
		UserDeleted(userId string) error
	}
)
