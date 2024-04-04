package quest

import (
	"context"
	"time"

	"github.com/SQUASHD/hbit/rpg/rpgdb"
)

type (
	// Top level service for instantiation
	QuestService interface {
		QuestManagement
		EventQuestService
		InternalQuestService
	}

	QuestManagement interface {
		ListQuests(ctx context.Context, userId string) ([]QuestDTO, error)
	}

	EventQuestService interface {
		DeleteUserQuests(userId string) error
	}

	InternalQuestService interface {
		CleanUp() error
	}

	QuestDTO struct {
		QuestID     string    `json:"quest_id"`
		QuestType   string    `json:"quest_type"`
		Description string    `json:"description"`
		Title       string    `json:"title"`
		Details     string    `json:"details"`
		UpdatedAt   time.Time `json:"updated_at"`
	}

	BossQuestDetails struct {
		BossID string `json:"boss_id"`
	}

	ItemQuestDetails struct {
		ItemID      string `json:"item_id"`
		DropsNeeded int    `json:"drops_needed"`
	}
)

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
