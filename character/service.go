package character

import (
	"context"
)

type (
	Repository interface {
		List(ctx context.Context) (Characters, error)
		Create(ctx context.Context, form CreateCharacterData) (Character, error)
		Read(ctx context.Context, characterId string) (Character, error)
		Update(ctx context.Context, data UpdateCharacterData) (Character, error)
		Delete(ctx context.Context, characterId string) error
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
	characters, err := s.charRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	dtos := charactersToDtos(characters)

	return dtos, nil
}

func (s *service) Create(ctx context.Context, form CreateCharacterForm) (DTO, error) {
	char, err := s.charRepo.Create(ctx, form.CreateCharacterData)
	if err != nil {
		return DTO{}, err
	}

	return characterToDto(char), nil
}

func (s *service) Read(ctx context.Context, form ReadCharacterForm) (DTO, error) {
	char, err := s.charRepo.Read(ctx, form.CharacterId)
	if err != nil {
		return DTO{}, err
	}

	return characterToDto(char), nil
}

func (s *service) Update(ctx context.Context, form UpdateCharacterForm) (DTO, error) {
	char, err := s.charRepo.Update(ctx, form.UpdateCharacterData)
	if err != nil {
		return DTO{}, err
	}

	return characterToDto(char), nil
}

func (s *service) Delete(ctx context.Context, form DeleteCharacterForm) error {
	return s.charRepo.Delete(ctx, form.CharacterId)
}
