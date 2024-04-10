package rpg

import (
	"github.com/SQUASHD/hbit/rpg/rpgdb"
)

func convertRPGStateToCharAndQuest(state rpgdb.GetUserRpgStateRow) (rpgdb.CharacterState, rpgdb.UserQuest) {
	charState := rpgdb.CharacterState{
		EventID:        state.EventID,
		UserID:         state.UserID,
		ClassID:        state.ClassID,
		CharacterLevel: state.CharacterLevel,
		Experience:     state.Experience,
		Health:         state.Health,
		Mana:           state.Mana,
		Strength:       state.Strength,
		Dexterity:      state.Dexterity,
		Intelligence:   state.Intelligence,
		Timestamp:      state.Timestamp,
	}
	var details string

	quest := rpgdb.UserQuest{
		UserID:    state.UserID,
		QuestID:   state.UserQuestID,
		Completed: state.QuestCompleted,
		Details:   details,
	}

	return charState, quest
}
