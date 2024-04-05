package character

import (
	"context"
	"database/sql"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/rpg/rpgdb"
	"github.com/wagslane/go-rabbitmq"
)

type (
	service struct {
		db      *sql.DB
		queries *rpgdb.Queries
	}
)

func NewService(
	db *sql.DB,
	queries *rpgdb.Queries,
	publisher *rabbitmq.Publisher,
) CharacterService {
	return &service{
		db:      db,
		queries: queries,
	}
}

func (s *service) List(ctx context.Context) ([]DTO, error) {
	return []DTO{}, &hbit.Error{Code: hbit.EINTERNAL, Message: "Not implemented"}
}

func (s *service) CreateCharacter(ctx context.Context, form CreateCharacterForm) (DTO, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return DTO{}, &hbit.Error{Code: hbit.EINTERNAL, Message: "failed to start transaction"}
	}
	defer tx.Rollback()
	_, err = s.queries.ReadCharacter(ctx, form.CreateCharacterParams.UserID)
	if err == nil {
		return DTO{}, &hbit.Error{Code: hbit.ECONFLICT, Message: "Character already exists"}
	}
	char, err := s.queries.CreateCharacter(ctx, form.CreateCharacterParams)
	if err != nil {
		return DTO{}, err
	}
	if err = tx.Commit(); err != nil {
		return DTO{}, &hbit.Error{Code: hbit.EINTERNAL, Message: "failed to commit transaction"}
	}

	return characterToDto(char), nil
}

func (s *service) GetCharacter(ctx context.Context, userId string) (DTO, error) {
	char, err := s.queries.ReadCharacter(ctx, userId)
	if err != nil {
		return DTO{}, err
	}

	return characterToDto(char), nil
}

func (s *service) UpdateCharacter(ctx context.Context, form UpdateCharacterForm) (DTO, error) {
	char, err := s.queries.UpdateCharacter(ctx, form.UpdateCharacterParams)
	if err != nil {
		return DTO{}, err
	}

	return characterToDto(char), nil
}

func (s *service) DeleteCharacter(ctx context.Context, userId string) error {
	return s.queries.DeleteCharacter(ctx, userId)
}

func (s *service) CleanUp() error {
	return s.db.Close()
}
