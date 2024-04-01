package quest

import (
	"context"

	"github.com/SQUASHD/hbit/quest/database"
)

type (
	QuestDTO struct {
		ID           string  `json:"id"`
		Title        string  `json:"title"`
		Description  string  `json:"description"`
		Requirements *string `json:"requirements"`
		Rewards      *string `json:"rewards"`
	}

	Service interface {
		ListQuests(ctx context.Context, userId string) ([]QuestDTO, error)
		CreateQuest(ctx context.Context, data database.CreateQuestParams) (QuestDTO, error)
		ReadQuest(ctx context.Context, id string) (QuestDTO, error)
		UpdateQuest(ctx context.Context, data database.UpdateQuestParams) (QuestDTO, error)
		DeleteQuest(ctx context.Context, id string) error
	}
)

func questToDTO(quest database.Quest) QuestDTO {
	return QuestDTO{
		ID:           quest.ID,
		Title:        quest.Title,
		Description:  quest.Description,
		Requirements: quest.Requirements,
		Rewards:      quest.Rewards,
	}
}

func questsToDTOs(quests []database.Quest) []QuestDTO {
	dtos := make([]QuestDTO, len(quests))
	for i, quest := range quests {
		dtos[i] = questToDTO(quest)
	}
	return dtos
}
