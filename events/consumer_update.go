package events

import (
	"encoding/json"
	"fmt"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/rpg"
	"github.com/SQUASHD/hbit/updates"
	"github.com/wagslane/go-rabbitmq"
)

func NewUpdateEventConsumer(url string) (*rabbitmq.Consumer, *rabbitmq.Conn, error) {
	conn, err := rabbitmq.NewConn(url)
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
func UpdatesMessageHandler(svc *updates.Service) func(d rabbitmq.Delivery) rabbitmq.Action {
	return func(d rabbitmq.Delivery) rabbitmq.Action {
		var event hbit.EventMessage
		if err := json.Unmarshal(d.Body, &event); err != nil {
			fmt.Printf("failed to unmarshal event: %v\n", err)
			return rabbitmq.NackDiscard
		}
		// TODO: Map events to dispatcher
		switch event.Type {
		case hbit.RPGREWARD:
			var payload rpg.TaskRewardPayload
			err := json.Unmarshal(event.Payload, &payload)
			if err != nil {
				return rabbitmq.NackDiscard
			}
			svc.SendMessageToUser(event.UserId, payload)
			return rabbitmq.Ack

		case hbit.LEVELUP:
			var payload rpg.CharacterLevelUpPayload
			err := json.Unmarshal(event.Payload, &payload)
			if err != nil {
				return rabbitmq.NackDiscard
			}
			svc.SendMessageToUser(event.UserId, payload)
			return rabbitmq.Ack

		default:
			return rabbitmq.NackDiscard

		}

	}
}
