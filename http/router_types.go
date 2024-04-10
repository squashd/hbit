package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/rpg"
	"github.com/SQUASHD/hbit/rpg/character"
	"github.com/SQUASHD/hbit/rpg/quest"
	"github.com/SQUASHD/hbit/rpg/rpgdb"
	"github.com/SQUASHD/hbit/task"
)

// NewTypesRouter returns a new http.ServeMux with routes to generate JSON
// responses for the various types
// internal

type rpgTestHandler struct {
	rpgSvc   rpg.Service
	questSvc quest.Service
	charSvc  character.Service
	taskSvc  task.Service
}

func newTestHandler(
	rpgSvc rpg.Service,
	questSvc quest.Service,
	charSvc character.Service,
	taskSvc task.Service,
) *rpgTestHandler {
	return &rpgTestHandler{
		rpgSvc:   rpgSvc,
		questSvc: questSvc,
		charSvc:  charSvc,
		taskSvc:  taskSvc,
	}
}

func NewTypesRouter(
	rpgSvc rpg.Service,
	questSvc quest.Service,
	charSvc character.Service,
	taskSvc task.Service,
) *http.ServeMux {
	r := http.NewServeMux()

	handler := newTestHandler(rpgSvc, questSvc, charSvc, taskSvc)
	rpgHandler := newRPGHandler(rpgSvc)
	taskResolver := newTaskResolutionHandler(taskSvc)

	r.HandleFunc("GET /task", taskDTO)
	r.HandleFunc("GET /character", charDTO)
	r.HandleFunc("GET /taskpayload", taskPayload)
	r.HandleFunc("GET /queststate", handler.getQuestState)
	r.HandleFunc("GET /taskreward", AuthMiddleware(taskRewardRequest))
	r.HandleFunc("GET /questdetails", questDetails)

	r.HandleFunc("GET /", getUserRPGState)
	r.HandleFunc("POST /questdone", rpgHandler.CalculateRewards)
	r.HandleFunc("GET /taskdone", taskResolver.Done)

	return r
}

func getUserRPGState(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, rpgdb.GetUserRpgStateRow{})
}

func (rh *rpgTestHandler) getQuestState(w http.ResponseWriter, r *http.Request) {
	bossHealth := new(int)
	itemsNeeded := new(int)

	*itemsNeeded = 5
	details := quest.QuestDetails{
		ItemsNeeded: itemsNeeded,
		BossHealth:  bossHealth,
		Rewards: quest.QuestRewards{
			Gold: 10,
			Exp:  10,
		},
	}

	detailsStr, err := json.Marshal(details)
	if err != nil {
		Error(w, r, err)
	}

	qs, err := rh.questSvc.GetQuestState(rpgdb.UserQuest{
		UserID:    "",
		QuestID:   "",
		Completed: false,
		Timestamp: time.Time{},
		EventID:   "",
		Details:   string(detailsStr),
	})
	if err != nil {
		Error(w, r, err)
		return
	}
	respondWithJSON(w, http.StatusOK, qs)
}

func taskDTO(w http.ResponseWriter, r *http.Request) {
	taskDTO := task.DTO{
		ID:          hbit.NewUUID(),
		UserID:      hbit.NewUUID(),
		Title:       "My first task",
		Text:        "This is longer text in the task",
		IsCompleted: false,
		Difficulty:  string(task.EASY),
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{}}
	respondWithJSON(w, http.StatusOK, taskDTO)
}

func charDTO(w http.ResponseWriter, r *http.Request) {
	characterDTO := character.CharacterDTO{
		Level:        0,
		Experience:   0,
		Health:       0,
		Mana:         0,
		Strength:     0,
		Dexterity:    0,
		Intelligence: 0,
	}
	respondWithJSON(w, http.StatusOK, characterDTO)
}
func taskRewardRequest(w http.ResponseWriter, r *http.Request, userId string) {
	res := rpg.TaskRewardRequest{
		Difficulty: task.EASY,
		UserId:     userId,
	}

	respondWithJSON(w, http.StatusOK, res)
}

func taskPayload(w http.ResponseWriter, r *http.Request) {
	taskDoneReq := hbit.TaskOrchestrationRequest{
		Difficulty: string(task.HARD),
	}
	respondWithJSON(w, http.StatusOK, taskDoneReq)
}

func questDetails(w http.ResponseWriter, r *http.Request) {
	itemsNeeded := new(int)
	bossHealth := new(int)
	rewards := quest.QuestRewards{
		Gold: 10,
		Exp:  150,
	}

	res := quest.QuestDetails{
		ItemsNeeded: itemsNeeded,
		BossHealth:  bossHealth,
		Rewards:     rewards,
	}

	respondWithJSON(w, http.StatusOK, res)
}
