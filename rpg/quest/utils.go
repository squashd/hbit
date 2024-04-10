package quest

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/SQUASHD/hbit/rpg/rpgdb"
)

type InternalUserQuestUtils interface {
	GetQuestState(rpgdb.UserQuest) (QuestState, error)
	HandleQuestChange(rpgdb.UserQuest, QuestProgression) (UserQuestDTO, error)
}

type QuestDetails struct {
	ItemsNeeded *int         `json:"items_needed"`
	BossHealth  *int         `json:"boss_health"`
	Rewards     QuestRewards `json:"rewards"`
}

func GetQuestDetails(jsonStr string) (QuestDetails, error) {
	var questDetails QuestDetails

	if err := json.Unmarshal([]byte(jsonStr), &questDetails); err != nil {
		return QuestDetails{}, err
	}

	return questDetails, nil
}

type QuestState struct {
	Active       bool         `json:"active"`
	DropsNeeded  *int         `json:"drops_needed"`
	DamageNeeded *int         `json:"damage_needed"`
	QuestRewards QuestRewards `json:"quest_rewards"`
}

func (s *service) GetQuestState(quest rpgdb.UserQuest) (QuestState, error) {
	var questDetails QuestDetails

	if err := json.Unmarshal([]byte(quest.Details), &questDetails); err != nil {
		fmt.Println("error unmarshalling quest details: ", err)
		return QuestState{}, err
	}

	if bool, _ := strconv.ParseBool(os.Getenv("DEBUG")); bool {
		fmt.Printf("questDetails: %+v\n", *questDetails.ItemsNeeded)
	}

	// updatedQuest := s.queries.

	return QuestState{
		Active:       !quest.Completed,
		DropsNeeded:  questDetails.ItemsNeeded,
		DamageNeeded: questDetails.BossHealth,
		QuestRewards: questDetails.Rewards,
	}, nil
}

type QuestProgression struct {
	Completed        bool
	DropChange       int
	BossHealthChange int
}

func (s *service) HandleQuestChange(userQuest rpgdb.UserQuest, progess QuestProgression) (UserQuestDTO, error) {
	return UserQuestDTO{}, nil
}
