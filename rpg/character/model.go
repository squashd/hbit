package character

import "github.com/SQUASHD/hbit/rpg/rpgdb"

type (
	CharacterDTO struct {
		Level        int64 `json:"level"`
		Experience   int64 `json:"experience"`
		Health       int64 `json:"health"`
		Mana         int64 `json:"mana"`
		Strength     int64 `json:"strength"`
		Dexterity    int64 `json:"dexterity"`
		Intelligence int64 `json:"intelligence"`
	}
)

func characterToDto(char rpgdb.CharacterState) CharacterDTO {
	return CharacterDTO{
		Level:        char.CharacterLevel,
		Experience:   char.Experience,
		Health:       char.Health,
		Mana:         char.Mana,
		Strength:     char.Strength,
		Dexterity:    char.Dexterity,
		Intelligence: char.Intelligence,
	}
}

func charactersToDtos(characters []rpgdb.CharacterState) []CharacterDTO {
	dtos := make([]CharacterDTO, len(characters))
	for i, char := range characters {
		dtos[i] = characterToDto(char)
	}
	return dtos
}
