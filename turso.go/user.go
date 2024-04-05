package user

import (
	"context"
	"time"

	"github.com/SQUASHD/hbit/user/userdb"
)

type (
	Service interface {
		GetSettings(ctx context.Context, userId string) (SettingsDTO, error)
		UpdateSettings(ctx context.Context, form UpdateSettingsForm) (SettingsDTO, error)
		CreateSettings(ctx context.Context, data userdb.CreateUserSettingsParams) (SettingsDTO, error)
		DeleteSettings(userId string) error
		Cleanup() error
	}

	SettingsDTO struct {
		UserID             string    `json:"user_id"`
		Theme              string    `json:"theme"`
		DisplayName        string    `json:"display_name"`
		Email              string    `json:"email"`
		EmailNotifications bool      `json:"email_notifications"`
		ResetTime          string    `json:"reset_time"`
		UpdatedAt          time.Time `json:"updated_at"`
	}

	UpdateSettingsForm struct {
		userdb.UpdateUserSettingsParams
		RequestedById string
	}
)

func toSettingsDTO(s userdb.UserSetting) SettingsDTO {
	return SettingsDTO{
		UserID:             s.UserID,
		Theme:              s.Theme,
		DisplayName:        s.DisplayName,
		Email:              s.Email,
		EmailNotifications: s.EmailNotifications,
		ResetTime:          s.ResetTime,
		UpdatedAt:          s.UpdatedAt,
	}
}
