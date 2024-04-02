package auth

import (
	"context"
	"database/sql"
	"errors"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/auth/authdb"
)

type repository struct {
	queries *authdb.Queries
	db      *sql.DB
}

// DeleteUser implements Repository.
func (r *repository) DeleteUser(userId string) error {
	return r.queries.DeleteUser(context.Background(), userId)
}

func NewRepository(db *sql.DB) Repository {
	queries := authdb.New(db)
	return &repository{
		queries: queries,
		db:      db,
	}
}
func (r *repository) CreateAuth(ctx context.Context, data authdb.CreateAuthParams) (authdb.Auth, error) {
	return r.queries.CreateAuth(ctx, data)
}

func (r *repository) FindUserByUsername(ctx context.Context, username string) (authdb.Auth, error) {
	return r.queries.FindUserByUsername(ctx, username)
}

func (r *repository) FindRevokeToken(ctx context.Context, token string) error {
	_, err := r.queries.FindRevokedToken(ctx, token)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) RevokeToken(ctx context.Context, form RevokeTokenForm) error {
	admin, err := r.queries.IsAdmin(ctx, form.RequesterId)
	if err != nil {
		return err
	}
	if admin == "" {
		return &hbit.Error{Code: hbit.EFORBIDDEN, Message: "You are not authorized to perform this action"}
	}
	return r.queries.CreateRevokedToken(ctx, form.CreateRevokedTokenParams)
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

func (r *repository) Cleanup() error {
	return r.db.Close()
}
