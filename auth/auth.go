package auth

import (
	"context"

	"github.com/SQUASHD/hbit/auth/authdb"
)

type (
	Service interface {
		UserAuth
		JwtAuth
		IsAdmin(ctx context.Context, userId string) (bool, error)
		Cleanup() error
	}

	UserAuth interface {
		Login(ctx context.Context, form LoginForm) (AuthDTO, error)
		Register(ctx context.Context, form CreateUserForm) (AuthDTO, error)
		DeleteUser(userId string) error
	}

	// JwtAuth is the interface for handling JWT tokens
	JwtAuth interface {
		AuthenticateUser(ctx context.Context, accessToken string) (userId string, err error)
		RefreshToken(ctx context.Context, refreshToken string) (accessToken, userId string, err error)
		RevokeToken(ctx context.Context, form RevokeTokenForm) error
	}

	LoginForm struct {
		Username string `json:"username" form:"username"`
		Password string `json:"password" form:"password"`
	}

	RevokeTokenForm struct {
		authdb.CreateRevokedTokenParams
		RequesterId string `json:"requester_id"`
	}

	CreateUserForm struct {
		Username        string `json:"username"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
)
