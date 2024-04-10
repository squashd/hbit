package auth

import (
	"context"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/auth/authdb"
)

type JwtAuth interface {
	AuthenticateUser(ctx context.Context, accessToken string) (userId string, err error)
	RefreshToken(ctx context.Context, refreshToken string) (accessToken, userId string, err error)
	RevokeToken(ctx context.Context, form RevokeTokenForm) error
}

func (s *service) AuthenticateUser(ctx context.Context, accessToken string) (userId string, err error) {
	id, err := ValidateJWT(accessToken, s.jwtConfig.JwtSecret)
	if err != nil {
		return "", &hbit.Error{Code: hbit.EUNAUTHORIZED, Message: "invalid token"}
	}
	return id, nil
}

func (s *service) RefreshToken(ctx context.Context, refreshToken string) (accessToken, userId string, err error) {
	id, err := ValidateJWT(refreshToken, s.jwtConfig.JwtSecret)
	if err != nil {
		return "", "", &hbit.Error{Code: hbit.EUNAUTHORIZED, Message: "invalid token"}
	}

	_, err = s.queries.FindRevokedToken(ctx, refreshToken)
	if err == nil {
		return "", "", &hbit.Error{Code: hbit.EUNAUTHORIZED, Message: "invalid token"}
	}

	token, err := s.makeAccessToken(id)
	if err != nil {
		return "", "", err
	}

	return token, id, nil
}

func (s *service) makeAccessToken(userId string) (string, error) {
	accessToken, err := MakeJWT(userId, s.jwtConfig.JwtSecret, s.jwtConfig.AccessIssuer, s.jwtConfig.AccessDuration)
	if err != nil {
		return "", &hbit.Error{Code: hbit.EINTERNAL, Message: "failed to generate access token"}
	}
	return accessToken, nil
}

func (s *service) makeRefreshToken(userId string) (string, error) {
	accessToken, err := MakeJWT(userId, s.jwtConfig.JwtSecret, s.jwtConfig.RefreshIssuer, s.jwtConfig.RefreshDuration)
	if err != nil {
		return "", &hbit.Error{Code: hbit.EINTERNAL, Message: "failed to generate access token"}
	}
	return accessToken, nil
}

type RevokeTokenForm struct {
	authdb.CreateRevokedTokenParams
	RequesterId string `json:"requester_id"`
}

func (s *service) RevokeToken(ctx context.Context, form RevokeTokenForm) error {
	err := s.queries.CreateRevokedToken(ctx, form.CreateRevokedTokenParams)
	if err != nil {
		return err
	}
	return nil
}
