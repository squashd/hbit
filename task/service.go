package task

import (
	"context"
	"encoding/json"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/events/eventpub"
)

type (
	Repository interface {
		List(ctx context.Context, userId string) (Tasks, error)
		Create(ctx context.Context, data CreateTaskData) (Task, error)
		Read(ctx context.Context, id string) (Task, error)
		Update(ctx context.Context, data UpdateTaskData) (Task, error)
		Delete(ctx context.Context, id string) error
	}

	service struct {
		repo      Repository
		publisher eventpub.Publisher
	}
)

func NewService(repo Repository, publisher eventpub.Publisher) Service {
	return &service{
		repo:      repo,
		publisher: publisher,
	}
}

func (s *service) List(ctx context.Context, RequestedById string) ([]DTO, error) {
	todos, err := s.repo.List(ctx, RequestedById)
	if err != nil {
		return nil, err
	}
	dtos := toDTOs(todos)
	return dtos, nil
}

func (s *service) Create(ctx context.Context, form CreateTaskForm) (DTO, error) {
	task, err := s.repo.Create(ctx, form.CreateTaskData)
	if err != nil {
		return DTO{}, err
	}

	dto := toDTO(task)
	return dto, nil
}

func (s *service) Update(ctx context.Context, form UpdateTaskForm) (DTO, error) {
	if form.UpdateTaskData.ID != form.TaskId {
		return DTO{}, &hbit.Error{Code: hbit.EUNAUTHORIZED, Message: "unauthorized"}
	}
	task, err := s.repo.Update(ctx, form.UpdateTaskData)
	if err != nil {
		return DTO{}, err
	}
	dto := toDTO(task)
	return dto, nil
}

func (s *service) Delete(ctx context.Context, form DeleteTaskForm) error {
	task, err := s.repo.Read(ctx, form.TaskId)
	if err != nil {
		return err
	}
	if task.UserID != form.RequestedById {
		return &hbit.Error{Code: hbit.EUNAUTHORIZED, Message: "unauthorized"}
	}

	err = s.repo.Delete(ctx, form.TaskId)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteTasks(msg json.RawMessage) error {
	panic("unimplemented")
}
