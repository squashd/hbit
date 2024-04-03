package events

import (
	"encoding/json"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/rpg"
	"github.com/wagslane/go-rabbitmq"
)

func NewRPGEventConsumer(url string) (*rabbitmq.Consumer, *rabbitmq.Conn, error) {
	conn, err := rabbitmq.NewConn(url)
	if err != nil {
		return nil, nil, err
	}
	consumer, err := rabbitmq.NewConsumer(
		conn,
		"rpg",
		rabbitmq.WithConsumerOptionsRoutingKey("task.*"),
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

type rpgConsumerHandler struct {
	rpgSvc rpg.EventService
}

func NewRPGConsumerHandler(svc rpg.Service) *rpgConsumerHandler {
	return &rpgConsumerHandler{rpgSvc: svc}
}

func (h *rpgConsumerHandler) HandleEvents(d rabbitmq.Delivery) rabbitmq.Action {
	var event hbit.EventMessage
	if err := json.Unmarshal(d.Body, &event); err != nil {
		return rabbitmq.NackDiscard
	}
	switch event.Type {
	case hbit.TaskCompleteEvent:
		if err := h.rpgSvc.HandleTaskCompleted(event.UserId); err != nil {
			return rabbitmq.NackDiscard
		}
		return rabbitmq.Ack
	default:
		return rabbitmq.NackDiscard
	}

}
