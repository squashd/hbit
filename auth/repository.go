package auth

import (
	"context"
	"database/sql"
	"errors"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/auth/database"
)

type repository struct {
	queries *database.Queries
}

// DeleteUser implements Repository.
func (r *repository) DeleteUser(userId string) error {
	return r.queries.DeleteUser(context.Background(), userId)
}

func NewRepository(db *sql.DB) Repository {
	queries := database.New(db)
	return &repository{
		queries: queries,
	}
}
func (r *repository) CreateAuth(ctx context.Context, data database.CreateAuthParams) (database.Auth, error) {
	return r.queries.CreateAuth(ctx, data)
}

func (r *repository) FindUserByUsername(ctx context.Context, username string) (database.Auth, error) {
	return r.queries.FindUserByUsername(ctx, username)
}

func (r *repository) FindRevokeToken(ctx context.Context, token string) error {
	_, err := r.queries.FindRevokedToken(ctx, token)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) RevokeToken(ctx context.Context, params database.CreateRevokedTokenParams) error {
	return r.queries.CreateRevokedToken(ctx, params)
}

func (r *repository) IsAdmin(ctx context.Context, userId string) (bool, error) {
	_, err := r.queries.IsAdmin(ctx, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, &hbit.Error{Code: hbit.EINTERNAL, Message: err.Error()}
	}
	return true, nil
}
