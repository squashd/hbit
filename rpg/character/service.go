package character

import (
	"context"

	"github.com/SQUASHD/hbit/rpg/rpgdb"
)

type (
	Repository interface {
		ListCharacters(ctx context.Context) ([]rpgdb.Character, error)
		CreateChracter(ctx context.Context, data rpgdb.CreateCharacterParams) (rpgdb.Character, error)
		ReadCharacter(ctx context.Context, characterId string) (rpgdb.Character, error)
		UpdateCharacter(ctx context.Context, data rpgdb.UpdateCharacterParams) (rpgdb.Character, error)
		DeleteCharacter(ctx context.Context, characterId string) error
		Cleanup() error
	}

	UserRepository interface {
		FindCharacter(ctx context.Context, userId string) (rpgdb.Character, error)
	}

	service struct {
		charRepo Repository
	}
)

func NewService(charRepo Repository) Service {
	return &service{
		charRepo: charRepo,
	}
}

func (s *service) List(ctx context.Context) ([]DTO, error) {
	characters, err := s.charRepo.ListCharacters(ctx)
	if err != nil {
		return nil, err
	}

	dtos := charactersToDtos(characters)

	return dtos, nil
}

func (s *service) CreateCharacter(ctx context.Context, form CreateCharacterForm) (DTO, error) {
	char, err := s.charRepo.CreateChracter(ctx, form.CreateCharacterParams)
	if err != nil {
		return DTO{}, err
	}

	return characterToDto(char), nil
}

func (s *service) GetCharacter(ctx context.Context, form ReadCharacterForm) (DTO, error) {
	char, err := s.charRepo.ReadCharacter(ctx, form.CharacterId)
	if err != nil {
		return DTO{}, err
	}

	return characterToDto(char), nil
}

func (s *service) UpdateCharacter(ctx context.Context, form UpdateCharacterForm) (DTO, error) {
	char, err := s.charRepo.UpdateCharacter(ctx, form.UpdateCharacterParams)
	if err != nil {
		return DTO{}, err
	}

	return characterToDto(char), nil
}

func (s *service) DeleteCharacter(ctx context.Context, form DeleteCharacterForm) error {
	return s.charRepo.DeleteCharacter(ctx, form.CharacterId)
}

func (s *service) Cleanup() error {
	return s.charRepo.Cleanup()
}
