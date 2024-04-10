package task

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/task/taskdb"
)

type UserTaskService interface {
	ListTasks(ctx context.Context, userId string) ([]DTO, error)
	CreateTask(ctx context.Context, form CreateTaskForm) (DTO, error)
	UpdateTask(ctx context.Context, form UpdateTaskForm) (DTO, error)
	DeleteTask(ctx context.Context, form DeleteTaskForm) error
}

func (s *service) ListTasks(ctx context.Context, requestedById string) ([]DTO, error) {
	todos, err := s.queries.ListTasks(ctx, requestedById)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	dtos := toDTOs(todos)
	return dtos, nil
}

type CreateTaskForm struct {
	CreateTaskRequest
	RequestedById string
}

func (s *service) CreateTask(ctx context.Context, form CreateTaskForm) (DTO, error) {
	model := CreateFormtoModel(form)
	task, err := s.queries.CreateTask(ctx, model)
	if err != nil {
		return DTO{}, err
	}

	go func() {
		if err := s.Publish(hbit.EventMessage{
			Type:    hbit.TASKCREATED,
			UserId:  form.RequestedById,
			EventId: hbit.NewEventIdWithTimestamp("task"),
			Payload: []byte{},
		}, []string{"task.created"}); err != nil {
			log.Printf("failed to publish task created event: %v", err)
		}
	}()

	dto := toDTO(task)
	return dto, nil
}

type UpdateTaskForm struct {
	taskdb.UpdateTaskParams
	TaskId        string
	RequestedById string
}

func (s *service) UpdateTask(ctx context.Context, form UpdateTaskForm) (DTO, error) {
	if form.UpdateTaskParams.ID != form.TaskId {
		return DTO{}, &hbit.Error{Code: hbit.EUNAUTHORIZED, Message: "unauthorized"}
	}
	task, err := s.queries.ReadTask(ctx, form.TaskId)
	if err != nil {
		return DTO{}, err
	}

	if task.UserID != form.RequestedById {
		return DTO{}, &hbit.Error{Code: hbit.EUNAUTHORIZED, Message: "unauthorized"}
	}

	updatedTask, err := s.queries.UpdateTask(ctx, form.UpdateTaskParams)
	if err != nil {
		return DTO{}, err
	}

	dto := toDTO(updatedTask)
	return dto, nil
}

type DeleteTaskForm struct {
	TaskId        string
	RequestedById string
}

func (s *service) DeleteTask(ctx context.Context, form DeleteTaskForm) error {
	task, err := s.queries.ReadTask(ctx, form.TaskId)
	if err != nil {
		return err
	}
	if task.UserID != form.RequestedById {
		return &hbit.Error{Code: hbit.EUNAUTHORIZED, Message: "unauthorized"}
	}

	err = s.queries.DeleteTask(ctx, form.TaskId)
	if err != nil {
		return err
	}
	return nil
}
