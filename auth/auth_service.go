package auth

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/auth/authdb"
	"github.com/SQUASHD/hbit/config"
	"github.com/wagslane/go-rabbitmq"
)

type (
	service struct {
		jwtConfig config.JwtOptions
		db        *sql.DB
		queries   *authdb.Queries
		publisher *rabbitmq.Publisher
	}
)

func NewService(
	publisher *rabbitmq.Publisher,
	jwtConfig config.JwtOptions,
	db *sql.DB,
	queries *authdb.Queries,
) Service {
	return &service{
		jwtConfig: jwtConfig,
		publisher: publisher,
		db:        db,
		queries:   queries,
	}
}

func (s *service) Register(ctx context.Context, form CreateUserForm) (AuthDTO, error) {
	var errs []*hbit.Error
	tx, er := s.db.Begin()
	if er != nil {
		return AuthDTO{}, &hbit.Error{Code: hbit.EINTERNAL, Message: "failed to start transaction"}
	}
	defer tx.Rollback()

	qtx := s.queries.WithTx(tx)

	_, err := qtx.FindUserByUsername(ctx, form.Username)
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

	user, err := qtx.CreateAuth(ctx, userData)
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

	if err := tx.Commit(); err != nil {
		return AuthDTO{}, &hbit.Error{Code: hbit.EINTERNAL, Message: "failed to commit transaction"}
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

func convertUserFormToModel(form CreateUserForm, password string) authdb.CreateAuthParams {
	return authdb.CreateAuthParams{
		UserID:         form.UserID,
		Username:       form.Username,
		HashedPassword: password,
	}
}

func (s *service) Login(ctx context.Context, form LoginForm) (AuthDTO, error) {
	user, err := s.queries.FindUserByUsername(ctx, form.Username)
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

// TODO: Implement
func (s *service) SignOut(ctx context.Context) error {
	return nil
}

func (s *service) IsAdmin(ctx context.Context, userId string) (bool, error) {
	_, err := s.queries.IsAdmin(ctx, userId)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *service) DeleteUser(ctx context.Context, userId string) error {
	err := s.queries.DeleteUser(ctx, userId)
	if err != nil {
		return err
	}
	event, err := hbit.NewEventMessage(
		hbit.AUTHDELETE,
		hbit.UserId(userId),
		hbit.NewEventIdWithTimestamp("auth"),
		nil,
	)
	if err != nil {
		return err
	}
	err = s.Publish(event, []string{"auth.delete"})
	if err != nil {
		return err
	}
	return nil
}

// Publish implements the hbit.Publisher interface
func (s *service) Publish(event hbit.EventMessage, routingKeys []string) error {

	msg, err := json.Marshal(&event)
	if err != nil {
		return err
	}

	return s.publisher.Publish(
		msg,
		routingKeys,
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsExchange("events"),
	)
}

// Cleanup closes the publisher and the database connection
func (s *service) Cleanup() error {
	s.publisher.Close()
	return s.db.Close()
}
