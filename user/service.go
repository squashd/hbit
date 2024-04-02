package user

import (
	"context"
	"encoding/json"

	"github.com/SQUASHD/hbit/user/userdb"
)

type (
	Repository interface {
		CreateSettings(ctx context.Context, data userdb.CreateUserSettingsParams) (userdb.UserSetting, error)
		ReadSettings(ctx context.Context, userId string) (userdb.UserSetting, error)
		UpdateSettings(ctx context.Context, data userdb.UpdateUserSettingsParams) (userdb.UserSetting, error)
		DeleteSettings(ctx context.Context, userId string) error
		Cleanup() error
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
func (s *service) CreateSettings(ctx context.Context, data userdb.CreateUserSettingsParams) (SettingsDTO, error) {
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

func (s *service) Cleanup() error {
	return s.repo.Cleanup()
}
