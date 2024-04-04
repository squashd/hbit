package character

import (
	"context"
	"database/sql"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/rpg/rpgdb"
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
	char, err := s.queries.CreateCharacter(ctx, form.CreateCharacterParams)
	if err != nil {
		return DTO{}, err
	}

	return characterToDto(char), nil
}

func (s *service) GetCharacter(ctx context.Context, form ReadCharacterForm) (DTO, error) {
	char, err := s.queries.ReadCharacter(ctx, form.CharacterId)
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

func (s *service) DeleteCharacter(ctx context.Context, form DeleteCharacterForm) error {
	return s.queries.DeleteCharacter(ctx, form.CharacterId)
}

func (s *service) CleanUp() error {
	return s.db.Close()
}
