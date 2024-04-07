package http

import (
	"encoding/json"
	"net/http"

	"github.com/SQUASHD/hbit"
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
	respondWithJSON(w, http.StatusOK, tasks)
}

func (h *taskHandler) Get(w http.ResponseWriter, r *http.Request, requestedById string) {
	todos, err := h.taskSvc.List(r.Context(), requestedById)
	if err != nil {
		Error(w, r, err)
		return
	}
	respondWithJSON(w, http.StatusOK, todos)

}

func (h *taskHandler) Create(w http.ResponseWriter, r *http.Request, requestedById string) {
	var data task.CreateTaskRequest

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		Error(w, r, &hbit.Error{Message: "Invalid JSON Body", Code: hbit.EINVALID})
		return
	}

	taskState := task.CreateTaskForm{
		CreateTaskRequest: data,
		RequestedById:     requestedById,
	}

	todo, err := h.taskSvc.Create(r.Context(), taskState)
	if err != nil {
		Error(w, r, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, todo)
}

func (h *taskHandler) Update(w http.ResponseWriter, r *http.Request, requestedById string) {
	taskId := r.PathValue("taskId")
	var data taskdb.UpdateTaskParams

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		Error(w, r, err)
		return
	}

	taskState := task.UpdateTaskForm{
		UpdateTaskParams: data,
		TaskId:           taskId,
		RequestedById:    requestedById,
	}

	todo, err := h.taskSvc.Update(r.Context(), taskState)
	if err != nil {
		Error(w, r, err)
		return
	}

	respondWithJSON(w, http.StatusOK, todo)
}

func (h *taskHandler) Delete(w http.ResponseWriter, r *http.Request, requestedById string) {
	taskId := r.PathValue("taskId")

	taskState := task.DeleteTaskForm{
		TaskId:        taskId,
		RequestedById: requestedById,
	}

	err := h.taskSvc.Delete(r.Context(), taskState)
	if err != nil {
		Error(w, r, err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (h *taskHandler) Done(w http.ResponseWriter, r *http.Request) {
	var taskState task.TaskStatePayload

	if err := json.NewDecoder(r.Body).Decode(&taskState); err != nil {
		Error(w, r, &hbit.Error{Code: hbit.EINVALID, Message: "Invalid JSON Body"})
		return
	}

	task, err := h.taskSvc.TaskDone(r.Context(), taskState)
	if err != nil {
		Error(w, r, err)
		return
	}

	respondWithJSON(w, http.StatusOK, task)
}

func (h *taskHandler) Undone(w http.ResponseWriter, r *http.Request) {
	var taskState task.TaskStatePayload

	if err := json.NewDecoder(r.Body).Decode(&taskState); err != nil {
		Error(w, r, &hbit.Error{Code: hbit.EINVALID, Message: "Invalid JSON Body"})
		return
	}

	task, err := h.taskSvc.TaskUndone(r.Context(), taskState)
	if err != nil {
		Error(w, r, err)
		return
	}

	respondWithJSON(w, http.StatusOK, task)
}
