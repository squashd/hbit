package events

import (
	"encoding/json"
	"fmt"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/config"
	"github.com/SQUASHD/hbit/rpg"
	"github.com/wagslane/go-rabbitmq"
)

func NewUpdateEventConsumer(rabbitmqConf config.RabbitMQ) (*rabbitmq.Consumer, *rabbitmq.Conn, error) {
	connStr := config.NewRabbitConnectionString(rabbitmqConf)
	conn, err := rabbitmq.NewConn(connStr)
	if err != nil {
		return nil, nil, err
	}
	consumer, err := rabbitmq.NewConsumer(
		conn,
		"updates",
		rabbitmq.WithConsumerOptionsRoutingKey("rpg.#"),
		rabbitmq.WithConsumerOptionsRoutingKey("task.#"),
		rabbitmq.WithConsumerOptionsExchangeName("events"),
		rabbitmq.WithConsumerOptionsExchangeKind("topic"),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
	)
	if err != nil {
		return nil, nil, err
	}

	return consumer, conn, nil
}

type updateConsumerHandler struct {
}

func NewUpdateConsumerHandler() *updateConsumerHandler {
	return &updateConsumerHandler{}
}

func (h *updateConsumerHandler) HandleEvents(d rabbitmq.Delivery) rabbitmq.Action {
	var event hbit.EventMessage
	if err := json.Unmarshal(d.Body, &event); err != nil {
		fmt.Printf("failed to unmarshal event: %v\n", err)
		return rabbitmq.NackDiscard
	}
	switch event.Type {
	case "task_reward":
		var payload rpg.TaskRewardPayload
		err := json.Unmarshal(event.Payload, &payload)
		if err != nil {
			fmt.Printf("failed to unmarshal payload: %v\n", err)
			return rabbitmq.NackDiscard
		}
		fmt.Printf("update consumer received task reward event: %v\n", payload)
		return rabbitmq.Ack
	case "level_up":
		var payload rpg.CharacterLevelUpPayload
		err := json.Unmarshal(event.Payload, &payload)
		if err != nil {
			fmt.Printf("failed to unmarshal payload: %v\n", err)
			return rabbitmq.NackDiscard
		}
		fmt.Printf("update consumer received task reward event: %v\n", payload)
		return rabbitmq.Ack
	case "task_complete":
		fmt.Printf("update consumer received task completed event: %v\n", event)
		return rabbitmq.Ack
	default:
		fmt.Printf("update consumer received unknown event: %v\n", event.Type)
		return rabbitmq.NackDiscard

	}

}
