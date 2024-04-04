package rpg

import "github.com/SQUASHD/hbit/rpg/rpgdb"

func calculateExperienceToNextLevel(level int) int {
	return 100 * level
}

type TaskDifficulty string

const (
	Easy   TaskDifficulty = "easy"
	Medium TaskDifficulty = "medium"
	Hard   TaskDifficulty = "hard"
)

func calculateReward(character rpgdb.CharacterState, difficulty TaskDifficulty) int {
	switch difficulty {
	case Easy:
		return 100
	case Medium:
		return 200
	case Hard:
		return 300
	default:
		return 0
	}
}

func getExperienceGain(difficulty TaskDifficulty) int {
	switch difficulty {
	case Easy:
		return 10
	case Medium:
		return 20
	case Hard:
		return 30
	default:
		return 0
	}
}

func calculateRewardModifier(character rpgdb.CharacterState) float64 {
	return 1.0 + (float64(character.CharacterLevel) * 0.1)
}
