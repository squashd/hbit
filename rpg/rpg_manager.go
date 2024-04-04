package rpg

import (
	"github.com/SQUASHD/hbit/rpg/rpgdb"
	"github.com/SQUASHD/hbit/task"
)

const (
	BASE_GOLD_REWARD = 10
	BASE_EXP_REWARD  = 10
	BASE_MANA_REWARD = 10
	BASE_DMG_DONE    = 10
)

func calculateDamageDone(char rpgdb.CharacterState, taskDifficulty task.TaskDifficulty) int {
	var difficultyMultiplier float64
	switch taskDifficulty {
	case task.EASY:
		difficultyMultiplier = 0.5
	case task.MEDIUM:
		difficultyMultiplier = 1
	case task.HARD:
		difficultyMultiplier = 1.5
	case task.EPIC:
		difficultyMultiplier = 2
	default:
		difficultyMultiplier = 1

	}
	dexModifier := 1 + (float64(char.Dexterity) / 100)
	strModifier := 1 + (float64(char.Strength) / 100)
	levelModifier := 1 + (float64(char.CharacterLevel) / 100)

	return int(float64(BASE_DMG_DONE) * dexModifier * strModifier * levelModifier * difficultyMultiplier)

}

func determineReward(char rpgdb.CharacterState, taskDifficulty task.TaskDifficulty) TaskRewardPayload {
	var difficultyMultiplier float64
	switch taskDifficulty {
	case task.EASY:
		difficultyMultiplier = 0.5
	case task.MEDIUM:
		difficultyMultiplier = 1
	case task.HARD:
		difficultyMultiplier = 1.5
	case task.EPIC:
		difficultyMultiplier = 2
	default:
		difficultyMultiplier = 1
	}

	dexModifier := 1 + (float64(char.Dexterity) / 100)
	levelModifier := 1 + (float64(char.CharacterLevel) / 100)
	intModifier := 1 + (float64(char.Intelligence) / 100)

	goldReward := int(float64(BASE_GOLD_REWARD) * dexModifier * levelModifier * float64(difficultyMultiplier))
	expReward := int(float64(BASE_EXP_REWARD) * intModifier * levelModifier * float64(difficultyMultiplier))
	manaReward := int(float64(BASE_MANA_REWARD) * intModifier * levelModifier * float64(difficultyMultiplier))

	return TaskRewardPayload{
		Gold: goldReward,
		Exp:  expReward,
		Mana: manaReward,
	}
}
