package auth

import (
	"context"

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

	Service interface {
		Login(ctx context.Context, form LoginForm) (AuthDTO, error)
		Register(ctx context.Context, form CreateUserForm) (AuthDTO, error)
		AuthenticateUser(ctx context.Context, accessToken string) (userId string, err error)
		RefreshToken(ctx context.Context, refreshToken string) (accessToken, userId string, err error)
		RevokeToken(ctx context.Context, form RevokeTokenForm) error
		IsAdmin(ctx context.Context, userId string) (bool, error)
	}
)

func toDTO(model database.Auth, accessToken, refreshToken string) AuthDTO {
	return AuthDTO{
		Username:     model.Username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
