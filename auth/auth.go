package auth

import (
	"context"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/auth/authdb"
)

type (
	Service interface {
		UserAuth
		JwtAuth
		hbit.Publisher
		IsAdmin(ctx context.Context, userId string) (bool, error)
		Cleanup() error
	}

	UserAuth interface {
		Login(ctx context.Context, form LoginForm) (AuthDTO, error)
		Register(ctx context.Context, form CreateUserForm) (AuthDTO, error)
		// DeleteUser deletes a user and all associated data and also publsihes an event
		// Currently only services is orchestrated with registration, but it may be necessary to
		// orchestrate with other services in the future
		DeleteUser(ctx context.Context, userId string) error
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

	// Since we are orchestrating registration we need to supply the user id
	CreateUserForm struct {
		UserID          string `json:"userId"`
		Username        string `json:"username"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
)
