package http

import (
	"net/http"

	"github.com/SQUASHD/hbit/task"
	"github.com/SQUASHD/hbit/task/taskdb"
)

type taskHandler struct {
	taskSvc task.UserTaskService
}

func newTaskHandler(svc task.UserTaskService) *taskHandler {
	return &taskHandler{taskSvc: svc}
}

func (h *taskHandler) FindAll(w http.ResponseWriter, r *http.Request, userId string) {
	tasks, err := h.taskSvc.List(r.Context(), userId)
	if err != nil {
		Error(w, r, err)
		return
	}
	RespondWithJSON(w, http.StatusOK, tasks)
}

func (h *taskHandler) Get(w http.ResponseWriter, r *http.Request, requestedById string) {
	todos, err := h.taskSvc.List(r.Context(), requestedById)
	if err != nil {
		Error(w, r, err)
		return
	}
	RespondWithJSON(w, http.StatusOK, todos)

}

func (h *taskHandler) Create(w http.ResponseWriter, r *http.Request, requestedById string) {
	var data taskdb.CreateTaskParams

	if err := Decode(r, &data); err != nil {
		Error(w, r, err)
		return
	}

	form := task.CreateTaskForm{
		CreateTaskParams: data,
		RequestedById:    requestedById,
	}

	todo, err := h.taskSvc.Create(r.Context(), form)
	if err != nil {
		Error(w, r, err)
		return
	}

	RespondWithJSON(w, http.StatusCreated, todo)
}

func (h *taskHandler) Update(w http.ResponseWriter, r *http.Request, requestedById string) {
	id := r.PathValue("id")
	var data taskdb.UpdateTaskParams

	if err := Decode(r, &data); err != nil {
		Error(w, r, err)
		return
	}

	form := task.UpdateTaskForm{
		UpdateTaskParams: data,
		TaskId:           id,
		RequestedById:    requestedById,
	}

	todo, err := h.taskSvc.Update(r.Context(), form)
	if err != nil {
		Error(w, r, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, todo)
}

func (h *taskHandler) Delete(w http.ResponseWriter, r *http.Request, requestedById string) {
	id := r.PathValue("id")

	form := task.DeleteTaskForm{
		TaskId:        id,
		RequestedById: requestedById,
	}

	err := h.taskSvc.Delete(r.Context(), form)
	if err != nil {
		Error(w, r, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
