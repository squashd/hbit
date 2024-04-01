package auth

import (
	"github.com/SQUASHD/hbit/auth/database"
)

type (
	AuthDTO struct {
		Username     string `json:"username"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	RevokeTokenParams = database.CreateRevokedTokenParams

	LoginForm struct {
		Username string `json:"username" form:"username"`
		Password string `json:"password" form:"password"`
	}

	RevokeTokenForm = database.CreateRevokedTokenParams

	CreateUserForm struct {
		Username        string `json:"username"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
)

func toDTO(model database.Auth, accessToken, refreshToken string) AuthDTO {
	return AuthDTO{
		Username:     model.Username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
