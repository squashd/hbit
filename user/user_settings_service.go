package user

import (
	"context"

	"github.com/SQUASHD/hbit/user/userdb"
)

type UserSettingsService interface {
	GetSettings(ctx context.Context, userId string) (SettingsDTO, error)
	UpdateSettings(ctx context.Context, form UpdateSettingsForm) (SettingsDTO, error)
}

func (s *service) GetSettings(ctx context.Context, userId string) (SettingsDTO, error) {
	settings, err := s.queries.ReadUserSettings(ctx, userId)
	if err != nil {
		return SettingsDTO{}, err
	}

	return toSettingsDTO(settings), nil
}

type UpdateSettingsForm struct {
	userdb.UpdateUserSettingsParams
	RequestedById string
}

func (s *service) UpdateSettings(ctx context.Context, form UpdateSettingsForm) (SettingsDTO, error) {
	settings, err := s.queries.UpdateUserSettings(ctx, form.UpdateUserSettingsParams)
	if err != nil {
		return SettingsDTO{}, err
	}

	return toSettingsDTO(settings), nil
}
