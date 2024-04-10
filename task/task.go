package task

import (
	"database/sql"
	"encoding/json"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/task/taskdb"
	"github.com/wagslane/go-rabbitmq"
)

type (
	Service interface {
		UserTaskService
		TaskResolutionService
		hbit.UserDataHandler
		CleanUp() error
	}

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

func (s *service) CleanUp() error {
	return s.db.Close()
}
