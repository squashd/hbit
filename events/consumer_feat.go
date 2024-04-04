package events

import (
	"encoding/json"
	"fmt"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/feat"
	"github.com/wagslane/go-rabbitmq"
)

func NewFeatEventConsumer(url string) (*rabbitmq.Consumer, *rabbitmq.Conn, error) {
	conn, err := rabbitmq.NewConn(url)
	if err != nil {
		return nil, nil, err
	}
	consumer, err := rabbitmq.NewConsumer(
		conn,
		"feats",
		rabbitmq.WithConsumerOptionsRoutingKey("task.*"),
		rabbitmq.WithConsumerOptionsRoutingKey("auth.delete"),
		rabbitmq.WithConsumerOptionsRoutingKey("auth.login"),
		rabbitmq.WithConsumerOptionsRoutingKey("character.level_up"),
		rabbitmq.WithConsumerOptionsRoutingKey("quest.complete"),
		rabbitmq.WithConsumerOptionsExchangeName("events"),
		rabbitmq.WithConsumerOptionsExchangeKind("topic"),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
	)
	if err != nil {
		return nil, nil, err
	}

	return consumer, conn, nil
}

func FeatsMessageHandler(svc feat.Service) func(d rabbitmq.Delivery) rabbitmq.Action {
	return func(d rabbitmq.Delivery) rabbitmq.Action {
		var event hbit.EventMessage
		if err := json.Unmarshal(d.Body, &event); err != nil {
			fmt.Printf("failed to unmarshal event: %v\n", err)
			return rabbitmq.NackDiscard
		}
		// TODO: Map events to dispatcher
		switch event.Type {
		case hbit.TASKDONE:
			// TODO: Log task done
			return rabbitmq.Ack

		case hbit.TASKUNDO:
			// TODO: Log task undone
			return rabbitmq.Ack

		case hbit.RPGREWARD:
			// TODO: Log RPG reward
			return rabbitmq.Ack

		case hbit.LEVELUP:
			// TODO: Log level up
			return rabbitmq.Ack

		default:
			return rabbitmq.NackDiscard

		}
	}
}
