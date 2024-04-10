package user

import (
	"context"

	"github.com/SQUASHD/hbit/user/userdb"
)

type InternalUserSettingsService interface {
	CreateSettings(ctx context.Context, data userdb.CreateUserSettingsParams) (SettingsDTO, error)
}

// TODO: Orchestrate with registration
func (s *service) CreateSettings(ctx context.Context, data userdb.CreateUserSettingsParams) (SettingsDTO, error) {
	settings, err := s.queries.CreateUserSettings(ctx, data)
	if err != nil {
		return SettingsDTO{}, err
	}

	return toSettingsDTO(settings), nil
}
