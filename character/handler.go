package character

import (
	"context"
)

type (
	CreateCharacterForm struct {
		CreateCharacterData
		RequestedById string `json:"requested_by_id"`
	}

	ReadCharacterForm struct {
		RequestedById string `json:"requested_by_id"`
		CharacterId   string `json:"character_id"`
	}

	UpdateCharacterForm struct {
		UpdateCharacterData
		RequestedById string `json:"requested_by_id"`
		CharacterId   string `json:"character_id"`
	}

	DeleteCharacterForm struct {
		RequestedById string `json:"requested_by_id"`
		CharacterId   string `json:"character_id"`
	}

	Service interface {
		List(ctx context.Context) ([]DTO, error)
		Create(ctx context.Context, form CreateCharacterForm) (DTO, error)
		Read(ctx context.Context, form ReadCharacterForm) (DTO, error)
		Update(ctx context.Context, form UpdateCharacterForm) (DTO, error)
		Delete(ctx context.Context, form DeleteCharacterForm) error
	}
)
