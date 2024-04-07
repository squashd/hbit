package task

import (
	"context"
	"time"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/task/taskdb"
)

// TaskDone marks a task as done. Only the completed field is updated.
func (s *service) TaskDone(ctx context.Context, payload TaskStatePayload) (DTO, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return DTO{}, err
	}
	defer tx.Rollback()
	task, err := s.queries.WithTx(tx).ReadTask(ctx, payload.TaskId)
	if err != nil {
		return DTO{}, err
	}
	// TODO: wrap in validation function
	if task.UserID != payload.UserId {
		return DTO{}, &hbit.Error{Code: hbit.EUNAUTHORIZED, Message: "unauthorized"}
	}
	if task.IsCompleted {
		return DTO{}, &hbit.Error{Code: hbit.EINVALID, Message: "task already done"}
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

// TaskUndone marks a task as undone. Only the completed field is updated.
func (s *service) TaskUndone(ctx context.Context, payload TaskStatePayload) (DTO, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return DTO{}, err
	}
	defer tx.Rollback()
	task, err := s.queries.WithTx(tx).ReadTask(ctx, payload.TaskId)
	if err != nil {
		return DTO{}, err
	}

	// TODO: wrap in validation function and return MultiErr
	if task.UserID != payload.UserId {
		return DTO{}, &hbit.Error{Code: hbit.EUNAUTHORIZED, Message: "unauthorized"}
	}
	if !task.IsCompleted {
		return DTO{}, &hbit.Error{Code: hbit.EINVALID, Message: "task already undone"}
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
