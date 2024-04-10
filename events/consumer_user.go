package events

import (
	"context"
	"encoding/json"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/user"
	"github.com/wagslane/go-rabbitmq"
)

func NewUserEventConsumer(url string) (*rabbitmq.Consumer, *rabbitmq.Conn, error) {
	conn, err := rabbitmq.NewConn(url)
	if err != nil {
		return nil, nil, err
	}
	consumer, err := rabbitmq.NewConsumer(
		conn,
		"user",
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

func UserMessageHandler(svc user.Service) func(d rabbitmq.Delivery) rabbitmq.Action {
	return func(d rabbitmq.Delivery) rabbitmq.Action {
		var event hbit.EventMessage
		if err := json.Unmarshal(d.Body, &event); err != nil {
			return rabbitmq.NackDiscard
		}
		ctx := context.Background()
		switch event.Type {
		case hbit.AUTHDELETE:
			if err := svc.DeleteData(ctx, event.UserId); err != nil {
				return rabbitmq.NackRequeue
			}
			return rabbitmq.Ack
		default:
			return rabbitmq.NackDiscard
		}
	}
}
