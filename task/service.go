package task

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/task/taskdb"
	"github.com/google/uuid"
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
	if err != nil {
		return nil, err
	}
	dtos := toDTOs(todos)
	return dtos, nil
}

func (s *service) Create(ctx context.Context, form CreateTaskForm) (DTO, error) {
	if form.CreateTaskParams.UserID != form.RequestedById {
		return DTO{}, &hbit.Error{Code: hbit.EUNAUTHORIZED, Message: "unauthorized"}
	}

	id := form.CreateTaskParams.ID
	if _, err := uuid.Parse(id); err != nil {
		return DTO{}, &hbit.Error{Code: hbit.EINVALID, Message: "invalid task id"}
	}

	taskId := hbit.NewUUID()
	form.CreateTaskParams.ID = taskId

	task, err := s.queries.CreateTask(ctx, form.CreateTaskParams)
	if err != nil {
		return DTO{}, err
	}

	if err := s.Publish(hbit.EventMessage{
		Type:    hbit.TASKCREATED,
		UserId:  hbit.UserId(form.RequestedById),
		EventId: hbit.NewEventIdWithTimestamp("task"),
		Payload: []byte{},
	}, []string{"task.created"}); err != nil {
		return DTO{}, err
	}

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

func (s *service) Publish(event hbit.EventMessage, routingKeys []string) error {
	msg, err := json.Marshal(event)
	if err != nil {
		return err
	}
	err = s.publisher.Publish(
		msg,
		routingKeys,
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsExchange("events"),
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteUserTasks(userId string) error {
	return s.queries.DeleteUserTasks(context.Background(), userId)
}

func (s *service) CleanUp() error {
	return s.db.Close()
}
