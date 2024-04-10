package events

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/task"
	"github.com/wagslane/go-rabbitmq"
)

func NewTaskEventConsumer(url string) (*rabbitmq.Consumer, *rabbitmq.Conn, error) {
	conn, err := rabbitmq.NewConn(url)
	if err != nil {
		return nil, nil, err
	}
	consumer, err := rabbitmq.NewConsumer(
		conn,
		"Task",
		rabbitmq.WithConsumerOptionsRoutingKey("auth.delete"),
		rabbitmq.WithConsumerOptionsExchangeKind("topic"),
		rabbitmq.WithConsumerOptionsExchangeName("events"),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
	)
	if err != nil {
		return nil, nil, err
	}

	return consumer, conn, nil
}

func TaskMessageHandler(svc task.Service) func(d rabbitmq.Delivery) rabbitmq.Action {
	return func(d rabbitmq.Delivery) rabbitmq.Action {
		var event hbit.EventMessage
		if err := json.Unmarshal(d.Body, &event); err != nil {
			return rabbitmq.NackDiscard
		}
		switch event.Type {
		case hbit.AUTHDELETE:
			ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
			defer cancel()
			if err := svc.DeleteData(ctx, event.UserId); err != nil {
				log.Printf("Failed to delete user data: %v", err)
			}
			return rabbitmq.Ack
		default:
			return rabbitmq.NackDiscard
		}
	}
}
