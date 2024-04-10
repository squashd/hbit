package quest

import (
	"time"

	"github.com/SQUASHD/hbit/rpg/rpgdb"
)

type (
	// QuestItemDTO struct {
	// 	QuestID          string `json:"quest_id"`
	// 	ItemID           string `json:"item_id"`
	// 	QuantityRequired int64  `json:"quantity_required"`
	// }

	QuestDTO struct {
		QuestID     string    `json:"quest_id"`
		QuestType   string    `json:"quest_type"`
		Description string    `json:"description"`
		Title       string    `json:"title"`
		Details     string    `json:"details"`
		UpdatedAt   time.Time `json:"updated_at"`
	}

	QuestRewards struct {
		Gold int `json:"gold"`
		Exp  int `json:"exp"`
	}

	QuestProgress struct {
		CurrItems      *int `json:"curr_items"`
		CurrBossHealth *int `json:"curr_boss_health"`
	}

	UserQuestDTO struct {
		Completed bool         `json:"completed"`
		Details   QuestDetails `json:"details"`
	}
)

func userQuestToDTO(quest rpgdb.UserQuest) (UserQuestDTO, error) {
	details, err := GetQuestDetails(quest.Details)
	if err != nil {
		return UserQuestDTO{}, err
	}

	dto := UserQuestDTO{
		Completed: quest.Completed,
		Details:   details,
	}

	return dto, nil
}

type QuestType string

const (
	BossQuest QuestType = "boss"
	ItemQuest QuestType = "item"
)

// func questItemToDTO(questItem rpgdb.QuestItem) QuestItemDTO {
// 	return QuestItemDTO{
// 		QuestID:          questItem.QuestID,
// 		ItemID:           questItem.ItemID,
// 		QuantityRequired: questItem.QuantityRequired,
// 	}
// }
//
// func questItemToDTOs(questItems []rpgdb.QuestItem) []QuestItemDTO {
// 	dtos := make([]QuestItemDTO, len(questItems))
// 	for i, questItem := range questItems {
// 		dtos[i] = questItemToDTO(questItem)
// 	}
// 	return dtos
// }

func questToDTO(quest rpgdb.Quest) QuestDTO {
	return QuestDTO{
		QuestID:     quest.QuestID,
		QuestType:   quest.QuestType,
		Description: quest.Description,
		Title:       quest.Title,
		Details:     quest.Details,
		UpdatedAt:   quest.UpdatedAt,
	}
}

func questsToDTOs(quests []rpgdb.Quest) []QuestDTO {
	dtos := make([]QuestDTO, len(quests))
	for i, quest := range quests {
		dtos[i] = questToDTO(quest)
	}
	return dtos
}
