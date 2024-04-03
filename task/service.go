package task

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/task/taskdb"
	"github.com/wagslane/go-rabbitmq"
)

type (
	Repository interface {
		TaskCRUD
		DeleteUserTasks(userId string) error
		Cleanup() error
	}

	TaskCRUD interface {
		List(ctx context.Context, userId string) (Tasks, error)
		Create(ctx context.Context, data taskdb.CreateTaskParams) (Task, error)
		Read(ctx context.Context, id string) (Task, error)
		Update(ctx context.Context, data taskdb.UpdateTaskParams) (Task, error)
		Delete(ctx context.Context, id string) error
	}

	service struct {
		repo      Repository
		publisher *rabbitmq.Publisher
	}
)

func NewService(repo Repository, publisher *rabbitmq.Publisher) Service {
	if publisher == nil {
		fmt.Println("publisher is nil")
	}

	return &service{
		repo:      repo,
		publisher: publisher,
	}
}

func (s *service) List(ctx context.Context, requestedById string) ([]DTO, error) {
	todos, err := s.repo.List(ctx, requestedById)
	if err != nil {
		return nil, err
	}
	dtos := toDTOs(todos)
	return dtos, nil
}

func (s *service) Create(ctx context.Context, form CreateTaskForm) (DTO, error) {
	task, err := s.repo.Create(ctx, form.CreateTaskParams)
	if err != nil {
		return DTO{}, err
	}

	if s.publisher == nil {
		fmt.Println("publisher is nil")
		return DTO{}, nil
	}

	if err := s.publisher.Publish(
		[]byte("test"),
		[]string{"task.done"},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsExchange("events"),
	); err != nil {
		fmt.Println("cannot publish message")
	}

	dto := toDTO(task)
	return dto, nil
}

func (s *service) Update(ctx context.Context, form UpdateTaskForm) (DTO, error) {
	if form.UpdateTaskParams.ID != form.TaskId {
		return DTO{}, &hbit.Error{Code: hbit.EUNAUTHORIZED, Message: "unauthorized"}
	}
	task, err := s.repo.Update(ctx, form.UpdateTaskParams)
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

func (s *service) DeleteUserTasks(userId string) error {
	return s.repo.DeleteUserTasks(userId)
}

func (s *service) Cleanup() error {
	return s.repo.Cleanup()
}

func (s *service) Test(ctx context.Context, userId string) error {
	event := hbit.EventMessage{
		Type:   "task_complete",
		UserID: "42069",
	}
	msg, err := json.Marshal(&event)
	if err != nil {
		return err
	}

	err = s.publisher.Publish(
		msg,
		[]string{"task.done"},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsExchange("events"),
	)
	if err != nil {
		return err
	}

	fmt.Println("published task complete event")

	return nil
}
