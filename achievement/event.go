package achievement

import (
	"encoding/json"
	"log"

	"github.com/SQUASHD/hbit"
	"github.com/wagslane/go-rabbitmq"
)

func NewTaskDoneHandler(svc Service) func(rabbitmq.Delivery) rabbitmq.Action {
	return func(d rabbitmq.Delivery) rabbitmq.Action {
		var msg hbit.TaskCompleteMessage
		if err := json.Unmarshal(d.Body, &msg); err != nil {
			log.Printf("Error unmarshalling task complete message: %v", err)
			return rabbitmq.NackRequeue
		}

		if err := svc.TaskDone(msg.UserID); err != nil {
			log.Printf("Error handling task done for user %s: %v", msg.UserID, err)
			return rabbitmq.NackRequeue
		}

		return rabbitmq.Ack
	}
}
