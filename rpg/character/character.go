package character

import (
	"context"

	"github.com/SQUASHD/hbit/rpg/rpgdb"
)

type (
	CharacterService interface {
		List(ctx context.Context) ([]DTO, error)
		UserCharacterService
		CleanUp() error
	}

	UserCharacterService interface {
		CreateCharacter(ctx context.Context, form CreateCharacterForm) (DTO, error)
		GetCharacter(ctx context.Context, form ReadCharacterForm) (DTO, error)
		UpdateCharacter(ctx context.Context, form UpdateCharacterForm) (DTO, error)
		DeleteCharacter(ctx context.Context, form DeleteCharacterForm) error
	}

	CreateCharacterForm struct {
		rpgdb.CreateCharacterParams
		RequestedById string `json:"requested_by_id"`
	}

	ReadCharacterForm struct {
		RequestedById string `json:"requested_by_id"`
		CharacterId   string `json:"character_id"`
	}

	UpdateCharacterForm struct {
		rpgdb.UpdateCharacterParams
		RequestedById string `json:"requested_by_id"`
		CharacterId   string `json:"character_id"`
	}

	DeleteCharacterForm struct {
		RequestedById string `json:"requested_by_id"`
		CharacterId   string `json:"character_id"`
	}
)
