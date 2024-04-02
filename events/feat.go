package events

import (
	"encoding/json"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/config"
	"github.com/SQUASHD/hbit/feat"
	"github.com/wagslane/go-rabbitmq"
)

func NewFeatEventConsumer(rabbitmqConf config.RabbitMQ) (*rabbitmq.Consumer, *rabbitmq.Conn, error) {
	connStr := config.NewRabbitConnectionString(rabbitmqConf)
	conn, err := rabbitmq.NewConn(connStr)
	if err != nil {
		return nil, nil, err
	}
	consumer, err := rabbitmq.NewConsumer(
		conn,
		"events",
		rabbitmq.WithConsumerOptionsRoutingKey("task.*"),
		rabbitmq.WithConsumerOptionsRoutingKey("quest.*"),
		rabbitmq.WithConsumerOptionsRoutingKey("auth.delete"),
		rabbitmq.WithConsumerOptionsRoutingKey("char.*"),
		rabbitmq.WithConsumerOptionsExchangeKind("topic"),
		rabbitmq.WithConsumerOptionsExchangeName("events"),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
	)
	if err != nil {
		return nil, nil, err
	}

	return consumer, conn, nil
}

type EventHandler func(d rabbitmq.Delivery) rabbitmq.Action

type featConsumerHandler struct {
	featSvc feat.Service
}

func NewFeatConsumerHandler(svc feat.Service) *featConsumerHandler {
	return &featConsumerHandler{featSvc: svc}
}

func (h *featConsumerHandler) HandleEvents(d rabbitmq.Delivery) rabbitmq.Action {
	var event hbit.EventMessage
	if err := json.Unmarshal(d.Body, &event); err != nil {
		return rabbitmq.NackDiscard
	}
	switch event.Type {
	case "task.created":
		// do something
	case "task.updated":
		// do something
	case "task.deleted":
		// do something
	case "quest.created":
		// do something
	case "quest.updated":
		// do something
	case "quest.deleted":
		// do something
	case "auth.delete":
		// do something
	case "char.created":
		// do something
	case "char.updated":
		// do something
	case "char.deleted":
		// do something
	default:
		return rabbitmq.NackDiscard
	}

	return rabbitmq.Ack

}
