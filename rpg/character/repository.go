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

func (r *repository) List(ctx context.Context) (Characters, error) {
	return r.queries.ListCharacters(ctx)
}

func (r *repository) Create(ctx context.Context, data rpgdb.CreateCharacterParams) (Character, error) {
	return r.queries.CreateCharacter(ctx, data)
}

func (r *repository) Read(ctx context.Context, characterId string) (Character, error) {
	return r.queries.ReadCharacter(ctx, characterId)
}

func (r *repository) Update(ctx context.Context, data rpgdb.UpdateCharacterParams) (Character, error) {
	return r.queries.UpdateCharacter(ctx, data)
}

func (r *repository) Delete(ctx context.Context, characterId string) error {
	return r.queries.DeleteCharacter(ctx, characterId)
}

func (r *repository) Cleanup() error {
	return r.db.Close()
}
