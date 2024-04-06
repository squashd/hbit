package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/rpg"
	"github.com/SQUASHD/hbit/task"
)

type (
	TaskOrchestrator interface {
		OrchestrateTaskDone(w http.ResponseWriter, r *http.Request, userId string)
		OrchestrateTaskUndo(w http.ResponseWriter, r *http.Request, userId string)
	}

	orchestratorTask struct {
		taskSvcUrl string
		rpgSvcUrl  string
		client     *http.Client
	}
)

func NewTaskOrchestrator(
	taskSvcUrl, rpgSvcUrl string,
	client *http.Client,
) TaskOrchestrator {
	return &orchestratorTask{
		taskSvcUrl: taskSvcUrl,
		rpgSvcUrl:  rpgSvcUrl,
		client:     client,
	}
}

func (o *orchestratorTask) OrchestrateTaskDone(w http.ResponseWriter, r *http.Request, userId string) {
	taskId := r.PathValue("id")
	var taskDoneReq struct {
		task.UpdateTaskForm
		task.TaskDonePayload
	}
	err := json.NewDecoder(r.Body).Decode(&taskDoneReq)
	if err != nil {
		Error(w, r, &hbit.Error{Code: hbit.EINVALID, Message: "Invalid request payload"})
		return
	}

	taskDoneReq.UpdateTaskForm.ID = taskId
	taskDoneReq.UpdateTaskForm.RequestedById = userId

	var wg sync.WaitGroup

	var taskRes, rpgRes *http.Response
	var taskErr, rpgErr error

	wg.Add(1)
	go func() {
		defer wg.Done()
		taskRes, taskErr = o.callTaskDone(taskDoneReq.UpdateTaskForm)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		rpgRes, rpgErr = o.callRPGTaskDone(taskDoneReq.TaskDonePayload)
	}()

	wg.Wait()

	if rpgErr != nil && taskErr == nil {
		go o.callTaskUndone(taskDoneReq.UpdateTaskForm)
		Error(w, r, rpgErr)
		return
	}

	if taskErr != nil && rpgErr == nil {
		go o.callRPGTaskUndone(taskDoneReq.TaskDonePayload)
		Error(w, r, taskErr)
		return
	}

	if taskErr != nil && rpgErr != nil {
		Error(w, r, taskErr)
		return
	}

	var taskDTO task.DTO
	err = json.NewDecoder(taskRes.Body).Decode(&taskDTO)
	defer taskRes.Body.Close()
	if err != nil {
		Error(w, r, err)
		return
	}

	var rpgPayload rpg.TaskRewardPayload
	err = json.NewDecoder(rpgRes.Body).Decode(&rpgPayload)
	defer rpgRes.Body.Close()
	if err != nil {
		Error(w, r, err)
		return
	}

	taskDoneRes := struct {
		Task   task.DTO              `json:"task"`
		Reward rpg.TaskRewardPayload `json:"reward"`
	}{
		Task:   taskDTO,
		Reward: rpgPayload,
	}

	respondWithJSON(w, http.StatusOK, taskDoneRes)

}

func (o *orchestratorTask) callTaskDone(task task.UpdateTaskForm) (*http.Response, error) {
	taskData, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/%s/%s", o.taskSvcUrl, task.ID, "done")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(taskData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (o *orchestratorTask) callTaskUndone(task task.UpdateTaskForm) (*http.Response, error) {
	taskData, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/%s/%s", o.taskSvcUrl, task.ID, "undone")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(taskData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}

func (o *orchestratorTask) callRPGTaskDone(task task.TaskDonePayload) (*http.Response, error) {
	taskData, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/%s/%s", o.rpgSvcUrl, task.TaskId, "undone")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(taskData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}

func (o *orchestratorTask) callRPGTaskUndone(task task.TaskDonePayload) (*http.Response, error) {
	taskData, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/%s/%s", o.rpgSvcUrl, task.TaskId, "undone")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(taskData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}

func (o *orchestratorTask) OrchestrateTaskUndo(w http.ResponseWriter, r *http.Request, userId string) {
}
