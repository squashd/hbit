package auth

import "context"

type (
	Service interface {
		Login(ctx context.Context, form LoginForm) (AuthDTO, error)
		Register(ctx context.Context, form CreateUserForm) (AuthDTO, error)
		AuthenticateUser(ctx context.Context, accessToken string) (userId string, err error)
		RefreshToken(ctx context.Context, refreshToken string) (accessToken, userId string, err error)
		RevokeToken(ctx context.Context, form RevokeTokenForm) error
		IsAdmin(ctx context.Context, userId string) (bool, error)
		DeleteUser(userId string) error
	}
)
