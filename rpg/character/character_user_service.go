package character

import (
	"context"
	"database/sql"
	"errors"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/rpg/rpgdb"
)

func (s *service) List(ctx context.Context) ([]CharacterDTO, error) {
	return []CharacterDTO{}, &hbit.Error{Code: hbit.EINTERNAL, Message: "Not implemented"}
}

type CreateCharacterForm struct {
	rpgdb.CreateCharacterParams
	RequestedById string `json:"requested_by_id"`
}

func (s *service) CreateCharacter(ctx context.Context, form CreateCharacterForm) (CharacterDTO, error) {
	form.EventID = string(hbit.NewEventIdWithTimestamp("rpg"))
	tx, err := s.db.Begin()
	if err != nil {
		return CharacterDTO{}, &hbit.Error{Code: hbit.EINTERNAL, Message: "failed to start transaction"}
	}
	defer tx.Rollback()
	_, err = s.queries.ReadCharacter(ctx, form.CreateCharacterParams.UserID)
	if err == nil && err != sql.ErrNoRows {
		return CharacterDTO{}, &hbit.Error{Code: hbit.ECONFLICT, Message: "Character already exists"}
	}
	char, err := s.queries.CreateCharacter(ctx, form.CreateCharacterParams)
	if err != nil {
		return CharacterDTO{}, err
	}
	if err = tx.Commit(); err != nil {
		return CharacterDTO{}, &hbit.Error{Code: hbit.EINTERNAL, Message: "failed to commit transaction"}
	}

	return characterToDto(char), nil
}

func (s *service) GetCharacter(ctx context.Context, userId string) (CharacterDTO, error) {
	char, err := s.queries.ReadCharacter(ctx, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return CharacterDTO{}, &hbit.Error{Code: hbit.EINVALID, Message: "Character not found"}
		}
		return CharacterDTO{}, err
	}

	return characterToDto(char), nil
}

type UpdateCharacterForm struct {
	rpgdb.UpdateCharacterParams
	RequestedById string `json:"requested_by_id"`
	CharacterId   string `json:"character_id"`
}

func (s *service) UpdateCharacter(ctx context.Context, form UpdateCharacterForm) (CharacterDTO, error) {
	// TODO: Validate
	char, err := s.queries.UpdateCharacter(ctx, form.UpdateCharacterParams)
	if err != nil {
		return CharacterDTO{}, err
	}

	return characterToDto(char), nil
}

type DeleteCharacterForm struct {
	RequestedById string `json:"requested_by_id"`
	CharacterId   string `json:"character_id"`
}

func (s *service) DeleteCharacter(ctx context.Context, userId string) error {
	return s.queries.DeleteCharacter(ctx, userId)
}
