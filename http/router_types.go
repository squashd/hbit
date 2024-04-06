package http

import (
	"net/http"

	"github.com/SQUASHD/hbit/rpg/character"
	"github.com/SQUASHD/hbit/task"
)

// NewTypesRouter returns a new http.ServeMux with routes to generate JSON
// responses for the various types
// internal
func NewTypesRouter() *http.ServeMux {
	r := http.NewServeMux()
	r.HandleFunc("GET /task", func(w http.ResponseWriter, r *http.Request) {
		taskDTO := task.DTO{}
		respondWithJSON(w, http.StatusOK, taskDTO)
	})
	r.HandleFunc("GET /character", func(w http.ResponseWriter, r *http.Request) {
		characterDTO := character.DTO{
			Level:        0,
			Experience:   0,
			Health:       0,
			Mana:         0,
			Strength:     0,
			Dexterity:    0,
			Intelligence: 0,
		}
		respondWithJSON(w, http.StatusOK, characterDTO)
	})
	r.HandleFunc("GET /taskpayload", func(w http.ResponseWriter, r *http.Request) {
		var taskDoneReq struct {
			task.UpdateTaskForm
			task.TaskDonePayload
		}
		respondWithJSON(w, http.StatusOK, taskDoneReq)
	})

	return r
}
