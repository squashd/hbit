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

type featEventHandler struct {
	featSvc feat.Service
}

func NewFeatEventHandler(svc feat.Service) *featEventHandler {
	return &featEventHandler{featSvc: svc}
}

func (h *featEventHandler) HandleEvents(d rabbitmq.Delivery) rabbitmq.Action {
	var event hbit.EventMessage
	if err := json.Unmarshal(d.Body, &event); err != nil {
		fmt.Printf("failed to unmarshal event: %v\n", err)
		return rabbitmq.NackDiscard
	}
	switch event.Type {
	case "task_complete":
		fmt.Printf("feats service received task complete event for user: %s\n", event.UserId)
		return rabbitmq.Ack
	default:
		return rabbitmq.NackDiscard

	}
}
