package task

import (
	"context"
	"encoding/json"

	"github.com/SQUASHD/hbit"
	"github.com/wagslane/go-rabbitmq"
)

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

// This is used in response from auth publishing a delete event
func (s *service) DeleteUserTasks(userId string) error {
	return s.queries.DeleteUserTasks(context.Background(), userId)
}

func (s *service) CleanUp() error {
	return s.db.Close()
}
