package character

import (
	"context"
	"database/sql"

	"github.com/SQUASHD/hbit/rpg/rpgdb"
)

type repository struct {
	queries *rpgdb.Queries
	db      *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	queries := rpgdb.New(db)
	return &repository{queries: queries, db: db}
}

func (r *repository) ListCharacters(ctx context.Context) ([]rpgdb.Character, error) {
	return r.queries.ListCharacters(ctx)
}

func (r *repository) CreateChracter(ctx context.Context, data rpgdb.CreateCharacterParams) (rpgdb.Character, error) {
	return r.queries.CreateCharacter(ctx, data)
}

func (r *repository) ReadCharacter(ctx context.Context, characterId string) (rpgdb.Character, error) {
	return r.queries.ReadCharacter(ctx, characterId)
}

func (r *repository) UpdateCharacter(ctx context.Context, data rpgdb.UpdateCharacterParams) (rpgdb.Character, error) {
	return r.queries.UpdateCharacter(ctx, data)
}

func (r *repository) DeleteCharacter(ctx context.Context, characterId string) error {
	return r.queries.DeleteCharacter(ctx, characterId)
}

func (r *repository) Cleanup() error {
	return r.db.Close()
}
