package character

import character "github.com/SQUASHD/hbit/character/database"

type (
	CharacterClass    = character.CharacterClass
	CharacterClasses  = []CharacterClass
	CharacterClassDTO struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	CreateCharacterClassData = character.CreateCharacterClassParams
	UpdateCharacterClassData = character.UpdateCharacterClassParams
)

func classToDto(class CharacterClass) CharacterClassDTO {
	return CharacterClassDTO{
		Name:        class.Name,
		Description: class.Description,
	}
}

func classestoDTOs(classes CharacterClasses) []CharacterClassDTO {
	dtos := make([]CharacterClassDTO, len(classes))
	for i, class := range classes {
		dtos[i] = classToDto(class)
	}
	return dtos
}
