package character

import "github.com/SQUASHD/hbit/rpg/rpgdb"

type (
	CharacterClassDTO struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
)

func classToDto(class rpgdb.CharacterClass) CharacterClassDTO {
	return CharacterClassDTO{
		Name:        class.Name,
		Description: class.Description,
	}
}

func classestoDTOs(classes []rpgdb.CharacterClass) []CharacterClassDTO {
	dtos := make([]CharacterClassDTO, len(classes))
	for i, class := range classes {
		dtos[i] = classToDto(class)
	}
	return dtos
}
