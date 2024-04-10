package user

import (
	"time"

	"github.com/SQUASHD/hbit/user/userdb"
)

type SettingsDTO struct {
	UserID             string    `json:"user_id"`
	Theme              string    `json:"theme"`
	DisplayName        string    `json:"display_name"`
	Email              string    `json:"email"`
	EmailNotifications bool      `json:"email_notifications"`
	ResetTime          string    `json:"reset_time"`
	UpdatedAt          time.Time `json:"updated_at"`
}

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
