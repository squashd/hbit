package task

import (
	"context"
	"time"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/task/taskdb"
)

type (
	TaskResolutionService interface {
		TaskDone(ctx context.Context, payload TaskStateRequest) (DTO, error)
		TaskUndone(ctx context.Context, payload TaskStateRequest) (DTO, error)
	}

	TaskStateRequest struct {
		TaskId string `json:"taskId"`
		UserId string `json:"userId"`
	}
)

// TaskDone marks a task as done. Only the completed field is updated.
func (s *service) TaskDone(ctx context.Context, payload TaskStateRequest) (DTO, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return DTO{}, err
	}
	defer tx.Rollback()
	task, err := s.queries.WithTx(tx).ReadTask(ctx, payload.TaskId)
	if err != nil {
		return DTO{}, err
	}

	err = validateTaskDone(task, payload)
	if err != nil {
		return DTO{}, err
	}

	task, err = s.queries.UpdateTask(ctx, taskdb.UpdateTaskParams{
		ID:          payload.TaskId,
		Title:       task.Title,
		Text:        task.Text,
		IsCompleted: true,
		Difficulty:  task.Difficulty,
		UpdatedAt:   time.Now(),
	})
	if err != nil {
		return DTO{}, err
	}
	if err := tx.Commit(); err != nil {
		return DTO{}, err
	}
	dto := toDTO(task)
	return dto, nil
}

func validateTaskDone(task taskdb.Task, payload TaskStateRequest) error {
	if task.UserID != payload.UserId {
		return &hbit.Error{Code: hbit.EUNAUTHORIZED, Message: "unauthorized"}
	}
	if task.IsCompleted {
		return &hbit.Error{Code: hbit.EINVALID, Message: "task already done"}
	}

	return nil
}

// TaskUndone marks a task as undone. Only the completed field is updated.
func (s *service) TaskUndone(ctx context.Context, payload TaskStateRequest) (DTO, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return DTO{}, err
	}
	defer tx.Rollback()
	task, err := s.queries.WithTx(tx).ReadTask(ctx, payload.TaskId)
	if err != nil {
		return DTO{}, err
	}

	err = validateTaskUndone(task, payload)
	if err != nil {
		return DTO{}, err
	}

	task, err = s.queries.UpdateTask(ctx, taskdb.UpdateTaskParams{
		ID:          payload.TaskId,
		Title:       task.Title,
		Text:        task.Text,
		IsCompleted: false,
		Difficulty:  task.Difficulty,
		UpdatedAt:   time.Now(),
	})
	if err != nil {
		return DTO{}, err
	}
	if err := tx.Commit(); err != nil {
		return DTO{}, err
	}
	dto := toDTO(task)
	return dto, nil
}

func validateTaskUndone(task taskdb.Task, payload TaskStateRequest) error {
	if task.UserID != payload.UserId {
		return &hbit.Error{Code: hbit.EUNAUTHORIZED, Message: "unauthorized"}
	}
	if !task.IsCompleted {
		return &hbit.Error{Code: hbit.EINVALID, Message: "task already done"}
	}

	return nil
}
