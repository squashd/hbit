package auth

import (
	"context"
	"database/sql"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/auth/authdb"
	"github.com/SQUASHD/hbit/config"
	"github.com/wagslane/go-rabbitmq"
)

type (
	Service interface {
		UserAuth
		JwtAuth
		hbit.Publisher
		IsAdmin(ctx context.Context, userId string) (bool, error)
		Cleanup() error
	}
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
