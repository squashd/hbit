package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
	if taskId == "" {
		Error(w, r, &hbit.Error{Code: hbit.EINVALID, Message: "Invalid task ID"})
	}

	var request hbit.TaskOrchestrationRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		Error(w, r, &hbit.Error{Code: hbit.EINVALID, Message: "Invalid JSON Body"})
		return
	}

	taskStatePayload := task.TaskStatePayload{
		TaskId: taskId,
		UserId: userId,
	}
	rpgRewardPayload := rpg.CaclulateRewardPayload{
		Difficulty: task.TaskDifficulty(request.Difficulty),
		UserId:     userId,
	}

	var wg sync.WaitGroup

	var taskRes, rpgRes *http.Response
	var taskErr, rpgErr error

	wg.Add(2)
	go func() {
		defer wg.Done()
		taskRes, taskErr = o.callTaskDone(taskStatePayload)
	}()
	go func() {
		defer wg.Done()
		rpgRes, rpgErr = o.callRPGTaskDone(rpgRewardPayload)
	}()
	wg.Wait()

	if rpgErr != nil && taskErr == nil {
		go func() {
			if _, err := o.callTaskUndo(taskStatePayload); err != nil {
				log.Println(err)
				return
			}
		}()
		Error(w, r, rpgErr)
		return
	}

	if taskErr != nil && rpgErr == nil {
		go func() {
			if _, err := o.callRPGTaskUndo(rpgRewardPayload); err != nil {
				log.Println(err)
				return
			}
		}()
		Error(w, r, taskErr)
		return
	}

	if taskErr != nil && rpgErr != nil {
		Error(w, r, taskErr)
		return
	}

	// Because a response of status code 300 or higher is considered an error
	// we need to check the status code of the response and handle it accordingly
	var operationErrors []error

	if taskRes.StatusCode != http.StatusOK {
		err := parseResponseError(taskRes)
		operationErrors = append(operationErrors, err)
		go func() {
			if _, err := o.callRPGTaskUndo(rpgRewardPayload); err != nil {
				log.Println(err)
				return
			}
		}()
	}

	if rpgRes.StatusCode != http.StatusOK {
		err := parseResponseError(rpgRes)
		operationErrors = append(operationErrors, err)
		go func() {
			if _, err := o.callTaskUndo(taskStatePayload); err != nil {
				log.Println(err)
				return
			}
		}()
	}

	if len(operationErrors) > 0 {
		Error(w, r, operationErrors[0])
		return
	}

	var taskDTO task.DTO
	var rpgPayload rpg.TaskRewardPayload

	var resErrs []error
	var taskResErr, rpgResErr error

	wg.Add(2)
	go func() {
		defer wg.Done()
		taskResErr = json.NewDecoder(taskRes.Body).Decode(&taskDTO)
	}()
	go func() {
		defer wg.Done()
		rpgResErr = json.NewDecoder(rpgRes.Body).Decode(&rpgPayload)
	}()
	defer taskRes.Body.Close()
	defer rpgRes.Body.Close()
	wg.Wait()

	for _, err := range []error{taskResErr, rpgResErr} {
		if err != nil {
			resErrs = append(resErrs, err)
		}
	}

	if len(resErrs) > 0 {
		Error(w, r, errors.Join(resErrs...))
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

func (o *orchestratorTask) callTaskDone(payload task.TaskStatePayload) (*http.Response, error) {
	taskData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/done", o.taskSvcUrl)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(taskData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	setInternalHeader(req)
	resp, err := o.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (o *orchestratorTask) callTaskUndo(payload task.TaskStatePayload) (*http.Response, error) {
	taskData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/undo", o.taskSvcUrl)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(taskData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	setInternalHeader(req)

	resp, err := o.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (o *orchestratorTask) callRPGTaskDone(payload rpg.CaclulateRewardPayload) (*http.Response, error) {
	taskData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/rewards/calculate", o.rpgSvcUrl)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(taskData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	setInternalHeader(req)

	resp, err := o.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *orchestratorTask) callRPGTaskUndo(payload rpg.CaclulateRewardPayload) (*http.Response, error) {
	taskData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/rewards/undo", o.rpgSvcUrl)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(taskData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	setInternalHeader(req)

	resp, err := o.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *orchestratorTask) OrchestrateTaskUndo(w http.ResponseWriter, r *http.Request, userId string) {
	taskId := r.PathValue("id")
	if taskId == "" {
		Error(w, r, &hbit.Error{Code: hbit.EINVALID, Message: "Invalid task ID"})
	}

	var request hbit.TaskOrchestrationRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		Error(w, r, &hbit.Error{Code: hbit.EINVALID, Message: "Invalid request payload"})
		return
	}

	taskStatePayload := task.TaskStatePayload{
		TaskId: taskId,
		UserId: userId,
	}
	rpgRewardPayload := rpg.CaclulateRewardPayload{
		Difficulty: task.TaskDifficulty(request.Difficulty),
		UserId:     userId,
	}

	var wg sync.WaitGroup

	var taskRes, rpgRes *http.Response
	var taskErr, rpgErr error

	wg.Add(2)
	go func() {
		defer wg.Done()
		taskRes, taskErr = o.callTaskUndo(taskStatePayload)
	}()
	go func() {
		defer wg.Done()
		rpgRes, rpgErr = o.callRPGTaskUndo(rpgRewardPayload)
	}()
	wg.Wait()

	if rpgErr != nil && taskErr == nil {
		go func() {
			if _, err := o.callTaskDone(taskStatePayload); err != nil {
				log.Println(err)
				return
			}
		}()
		Error(w, r, rpgErr)
		return
	}

	if taskErr != nil && rpgErr == nil {
		go func() {
			if _, err := o.callRPGTaskDone(rpgRewardPayload); err != nil {
				log.Println(err)
				return
			}
		}()
		Error(w, r, taskErr)
		return
	}

	if taskErr != nil && rpgErr != nil {
		Error(w, r, taskErr)
		return
	}

	var operationErrors []error

	if taskRes.StatusCode != http.StatusOK {
		err := parseResponseError(taskRes)
		operationErrors = append(operationErrors, err)
		go func() {
			if _, err := o.callRPGTaskUndo(rpgRewardPayload); err != nil {
				log.Println(err)
				return
			}
		}()
	}

	if rpgRes.StatusCode != http.StatusOK {
		err := parseResponseError(rpgRes)
		operationErrors = append(operationErrors, err)
		go func() {
			if _, err := o.callTaskUndo(taskStatePayload); err != nil {
				log.Println(err)
				return
			}
		}()
	}

	if len(operationErrors) > 0 {
		Error(w, r, operationErrors[0])
		return
	}

	var taskDTO task.DTO
	var rpgPayload rpg.TaskRewardPayload

	var resErrs []error
	var taskResErr, rpgResErr error

	wg.Add(2)
	go func() {
		defer wg.Done()
		taskResErr = json.NewDecoder(taskRes.Body).Decode(&taskDTO)
	}()
	go func() {
		defer wg.Done()
		rpgResErr = json.NewDecoder(rpgRes.Body).Decode(&rpgPayload)
	}()
	defer taskRes.Body.Close()
	defer rpgRes.Body.Close()
	wg.Wait()

	for _, err := range []error{taskResErr, rpgResErr} {
		if err != nil {
			resErrs = append(resErrs, err)
		}
	}

	if len(resErrs) > 0 {
		Error(w, r, errors.Join(resErrs...))
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
