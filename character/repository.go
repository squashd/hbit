package character

import (
	"context"
	"database/sql"

	character "github.com/SQUASHD/hbit/character/database"
)

type repository struct {
	queries *character.Queries
}

func NewRepository(db *sql.DB) Repository {
	queries := character.New(db)
	return &repository{queries: queries}
}

func (r *repository) List(ctx context.Context) (Characters, error) {
	return r.queries.ListCharacters(ctx)
}

func (r *repository) Create(ctx context.Context, data CreateCharacterData) (Character, error) {
	return r.queries.CreateCharacter(ctx, data)
}

func (r *repository) Read(ctx context.Context, characterId string) (Character, error) {
	return r.queries.ReadCharacter(ctx, characterId)
}

func (r *repository) Update(ctx context.Context, data UpdateCharacterData) (Character, error) {
	return r.queries.UpdateCharacter(ctx, data)
}

func (r *repository) Delete(ctx context.Context, characterId string) error {
	return r.queries.DeleteCharacter(ctx, characterId)
}
