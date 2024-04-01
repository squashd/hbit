package auth

import (
	"context"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/auth/database"
	"github.com/SQUASHD/hbit/config"
	"github.com/wagslane/go-rabbitmq"
)

type (
	Repository interface {
		FindUserByUsername(ctx context.Context, username string) (database.Auth, error)
		CreateAuth(ctx context.Context, data database.CreateAuthParams) (database.Auth, error)
		FindRevokeToken(ctx context.Context, token string) error
		DeleteUser(userId string) error
		RevokeToken(ctx context.Context, form RevokeTokenParams) error
		IsAdmin(ctx context.Context, userId string) (bool, error)
	}

	service struct {
		jwtConfig config.JwtConfig
		repo      Repository
		publisher *rabbitmq.Publisher
	}
)

func NewService(
	repo Repository,
	jwtConfig config.JwtConfig,
	publisher *rabbitmq.Publisher,
) Service {
	return &service{
		jwtConfig: jwtConfig,
		repo:      repo,
		publisher: publisher,
	}
}

func (s *service) Register(ctx context.Context, form CreateUserForm) (AuthDTO, error) {
	var errs []*hbit.Error
	_, err := s.repo.FindUserByUsername(ctx, form.Username)
	if err == nil {
		errs = append(errs, &hbit.Error{Code: hbit.ECONFLICT, Message: "Username already exists"})
	}

	errs = append(errs, validateUsername(form.Username)...)
	errs = append(errs, validatePassword(form.Password, form.ConfirmPassword)...)

	if len(errs) > 0 {
		return AuthDTO{}, &hbit.MultiError{Errors: errs}
	}

	hashedPassword, err := HashPassword(form.Password)
	if err != nil {
		return AuthDTO{}, &hbit.Error{Code: hbit.EINTERNAL, Message: "Failed to hash password"}
	}

	userData := convertUserFormToModel(form, hashedPassword)

	user, err := s.repo.CreateAuth(ctx, userData)
	if err != nil {
		return AuthDTO{}, err
	}

	accessToken, err := s.makeAccessToken(user.UserID)
	if err != nil {
		return AuthDTO{}, &hbit.Error{Code: hbit.EINTERNAL, Message: "failed to create tokens"}
	}

	refreshtoken, err := s.makeRefreshToken(user.UserID)
	if err != nil {
		return AuthDTO{}, &hbit.Error{Code: hbit.EINTERNAL, Message: "failed to create tokens"}
	}

	dto := toDTO(user, accessToken, refreshtoken)

	return dto, nil
}

func validateUsername(username string) []*hbit.Error {
	var errs []*hbit.Error
	if len(username) < 5 {
		errs = append(errs, &hbit.Error{Code: hbit.EINVALID, Message: "username must be at least 5 characters long"})
	}
	// TODO: add more validation rules for username
	return errs
}

func validatePassword(password, confirmPassword string) []*hbit.Error {
	var errs []*hbit.Error
	if password != confirmPassword {
		errs = append(errs, &hbit.Error{Code: hbit.EINVALID, Message: "passwords do not match"})
	}

	if len(password) < 8 {
		errs = append(errs, &hbit.Error{Code: hbit.EINVALID, Message: "password must be at least 8 characters long"})
	}

	// TODO: add entropy check for password

	return errs
}

func convertUserFormToModel(form CreateUserForm, password string) database.CreateAuthParams {
	return database.CreateAuthParams{
		Username:       form.Username,
		HashedPassword: password,
	}
}

func (s *service) Login(ctx context.Context, form LoginForm) (AuthDTO, error) {
	user, err := s.repo.FindUserByUsername(ctx, form.Username)
	if err != nil {
		return AuthDTO{}, &hbit.Error{Code: hbit.EUNAUTHORIZED, Message: "invalid username or password"}
	}

	err = CheckPasswordHash(form.Password, user.HashedPassword)
	if err != nil {
		return AuthDTO{}, &hbit.Error{Code: hbit.EUNAUTHORIZED, Message: "invalid username or password"}
	}

	accessToken, err := s.makeAccessToken(user.UserID)
	if err != nil {
		return AuthDTO{}, &hbit.Error{Code: hbit.EINTERNAL, Message: "failed to create tokens"}
	}

	refreshtoken, err := s.makeRefreshToken(user.UserID)
	if err != nil {
		return AuthDTO{}, &hbit.Error{Code: hbit.EINTERNAL, Message: "failed to create tokens"}
	}

	dto := toDTO(user, accessToken, refreshtoken)

	return dto, nil
}

func (s *service) SignOut(ctx context.Context) error {
	panic("implement me")
}

func (s *service) AuthenticateUser(ctx context.Context, accessToken string) (userId string, err error) {
	id, err := ValidateJWT(accessToken, s.jwtConfig.JwtSecret, s.jwtConfig.AccessIssuer)
	if err != nil {
		return "", &hbit.Error{Code: hbit.EUNAUTHORIZED, Message: "invalid token"}
	}
	return id, nil
}

func (s *service) RefreshToken(ctx context.Context, refreshToken string) (accessToken, userId string, err error) {
	id, err := ValidateJWT(refreshToken, s.jwtConfig.JwtSecret, s.jwtConfig.RefreshIssuer)
	if err != nil {
		return "", "", &hbit.Error{Code: hbit.EUNAUTHORIZED, Message: "invalid token"}
	}

	err = s.repo.FindRevokeToken(ctx, refreshToken)
	if err == nil {
		return "", "", &hbit.Error{Code: hbit.EUNAUTHORIZED, Message: "invalid token"}
	}

	token, err := s.makeAccessToken(userId)
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

func (s *service) RevokeToken(ctx context.Context, form RevokeTokenForm) error {
	err := s.repo.RevokeToken(ctx, form)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) IsAdmin(ctx context.Context, userId string) (bool, error) {
	_, err := s.repo.IsAdmin(ctx, userId)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *service) DeleteUser(userId string) error {
	err := s.repo.DeleteUser(userId)
	if err != nil {
		return err
	}
	return nil
}
