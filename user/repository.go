package user

import (
	"context"
	"encoding/json"

	"github.com/SQUASHD/hbit/user/database"
)

type (
	Repository interface {
		CreateSettings(ctx context.Context, data database.CreateUserSettingsParams) (database.UserSetting, error)
		ReadSettings(ctx context.Context, userId string) (database.UserSetting, error)
		UpdateSettings(ctx context.Context, data database.UpdateUserSettingsParams) (database.UserSetting, error)
		DeleteSettings(ctx context.Context, userId string) error
	}

	service struct {
		repo Repository
	}
)

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// GetSettings implements Service.
func (s *service) GetSettings(ctx context.Context, userId string) (SettingsDTO, error) {
	panic("unimplemented")
}

// CreateSettings implements Service.
func (s *service) CreateSettings(ctx context.Context, data database.CreateUserSettingsParams) (SettingsDTO, error) {
	panic("unimplemented")
}

// UpdateSettings implements Service.
func (s *service) UpdateSettings(ctx context.Context, form UpdateSettingsForm) (SettingsDTO, error) {
	panic("unimplemented")
}

// DeleteSettings implements Service.
func (s *service) DeleteSettings(msg json.RawMessage) error {
	panic("unimplemented")
}
