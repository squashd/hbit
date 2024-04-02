package character

import "github.com/SQUASHD/hbit/rpg/rpgdb"

type (
	CharacterClass    = rpgdb.CharacterClass
	CharacterClasses  = []CharacterClass
	CharacterClassDTO struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	CreateCharacterClassData = rpgdb.CreateCharacterClassParams
	UpdateCharacterClassData = rpgdb.UpdateCharacterClassParams
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
