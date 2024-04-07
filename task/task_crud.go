package task

import (
	"context"
	"database/sql"
	"log"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/task/taskdb"
	"github.com/wagslane/go-rabbitmq"
)

type (
	service struct {
		db        *sql.DB
		queries   *taskdb.Queries
		publisher *rabbitmq.Publisher
	}
)

func NewService(
	db *sql.DB,
	queries *taskdb.Queries,
	publisher *rabbitmq.Publisher,
) Service {
	return &service{
		db:        db,
		queries:   queries,
		publisher: publisher,
	}
}

func (s *service) List(ctx context.Context, requestedById string) ([]DTO, error) {
	todos, err := s.queries.ListTasks(ctx, requestedById)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	dtos := toDTOs(todos)
	return dtos, nil
}

func (s *service) Create(ctx context.Context, form CreateTaskForm) (DTO, error) {
	model := CreateFormtoModel(form)
	task, err := s.queries.CreateTask(ctx, model)
	if err != nil {
		return DTO{}, err
	}

	go func() {
		if err := s.Publish(hbit.EventMessage{
			Type:    hbit.TASKCREATED,
			UserId:  hbit.UserId(form.RequestedById),
			EventId: hbit.NewEventIdWithTimestamp("task"),
			Payload: []byte{},
		}, []string{"task.created"}); err != nil {
			log.Printf("failed to publish task created event: %v", err)
		}
	}()

	dto := toDTO(task)
	return dto, nil
}

func (s *service) Update(ctx context.Context, form UpdateTaskForm) (DTO, error) {
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

func (s *service) Delete(ctx context.Context, form DeleteTaskForm) error {
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

