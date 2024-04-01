package quest

import "github.com/SQUASHD/hbit/quest/database"

type (
	QuestItemDTO struct {
		QuestID          string `json:"quest_id"`
		ItemID           string `json:"item_id"`
		QuantityRequired int64  `json:"quantity_required"`
	}
)

func questItemToDTO(questItem database.QuestItem) QuestItemDTO {
	return QuestItemDTO{
		QuestID:          questItem.QuestID,
		ItemID:           questItem.ItemID,
		QuantityRequired: questItem.QuantityRequired,
	}
}

func questItemToDTOs(questItems []database.QuestItem) []QuestItemDTO {
	dtos := make([]QuestItemDTO, len(questItems))
	for i, questItem := range questItems {
		dtos[i] = questItemToDTO(questItem)
	}
	return dtos
}
