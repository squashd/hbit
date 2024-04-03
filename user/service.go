package user

import (
	"context"
	"database/sql"

	"github.com/SQUASHD/hbit/user/userdb"
)

type (
	service struct {
		db      *sql.DB
		queries *userdb.Queries
	}
)

func NewService(
	db *sql.DB,
	queries *userdb.Queries,
) Service {
	return &service{
		db:      db,
		queries: queries,
	}
}

// GetSettings implements Service.
func (s *service) GetSettings(ctx context.Context, userId string) (SettingsDTO, error) {
	settings, err := s.queries.ReadUserSettings(ctx, userId)
	if err != nil {
		return SettingsDTO{}, err
	}

	return toSettingsDTO(settings), nil
}

// CreateSettings implements Service.
func (s *service) CreateSettings(ctx context.Context, data userdb.CreateUserSettingsParams) (SettingsDTO, error) {
	settings, err := s.queries.CreateUserSettings(ctx, data)
	if err != nil {
		return SettingsDTO{}, err
	}

	return toSettingsDTO(settings), nil
}

// UpdateSettings implements Service.
func (s *service) UpdateSettings(ctx context.Context, form UpdateSettingsForm) (SettingsDTO, error) {
	settings, err := s.queries.UpdateUserSettings(ctx, form.UpdateUserSettingsParams)
	if err != nil {
		return SettingsDTO{}, err
	}

	return toSettingsDTO(settings), nil
}

func (s *service) DeleteSettings(userId string) error {
	return s.queries.DeleteUserSettings(context.Background(), userId)
}

func (s *service) Cleanup() error {
	return s.db.Close()
}
