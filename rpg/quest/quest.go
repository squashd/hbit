package quest

import (
	"context"

	"github.com/SQUASHD/hbit/rpg/rpgdb"
)

type (
	Service interface {
		UserQuestService
		EventQuestService
		InternalQuestService
	}

	UserQuestService interface {
		ListQuests(ctx context.Context, userId string) ([]QuestDTO, error)
		ReadQuest(ctx context.Context, id string) (QuestDTO, error)
	}

	EventQuestService interface {
		DeleteUserQuests(userId string) error
	}

	InternalQuestService interface {
		UpdateQuest(ctx context.Context, data rpgdb.UpdateQuestParams) (QuestDTO, error)
		CreateQuest(ctx context.Context, data rpgdb.CreateQuestParams) (QuestDTO, error)
		DeleteQuest(ctx context.Context, id string) error
		Cleanup() error
	}

	QuestDTO struct {
		ID           string  `json:"id"`
		Title        string  `json:"title"`
		Description  string  `json:"description"`
		Requirements *string `json:"requirements"`
		Rewards      *string `json:"rewards"`
	}
)

func questToDTO(quest rpgdb.Quest) QuestDTO {
	return QuestDTO{
		ID:           quest.ID,
		Title:        quest.Title,
		Description:  quest.Description,
		Requirements: quest.Requirements,
		Rewards:      quest.Rewards,
	}
}

func questsToDTOs(quests []rpgdb.Quest) []QuestDTO {
	dtos := make([]QuestDTO, len(quests))
	for i, quest := range quests {
		dtos[i] = questToDTO(quest)
	}
	return dtos
}
