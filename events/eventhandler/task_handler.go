package eventhandler

import (
	"encoding/json"
	"log"

	"github.com/SQUASHD/hbit/task"
	"github.com/wagslane/go-rabbitmq"
)

type TaskMessage struct {
	Type    string          `json:"type"`    // e.g., "create", "update", "delete"
	Payload json.RawMessage `json:"payload"` // Task details
}

func TaskHandler(svc task.Service) rabbitmq.Handler {
	return func(d rabbitmq.Delivery) rabbitmq.Action {
		var taskMsg TaskMessage
		if err := json.Unmarshal(d.Body, &taskMsg); err != nil {
			log.Printf("Error unmarshalling task message: %v", err)
			return rabbitmq.NackRequeue
		}

		switch taskMsg.Type {
		case "create":
			// Handle create task...
			if err := svc.DeleteTasks(taskMsg.Payload); err != nil {
				log.Printf("Error creating task: %v", err)
				return rabbitmq.NackRequeue
			}
		case "update":
			// Handle update task...
		case "delete":
			// Handle delete task...
		default:
			log.Printf("Unhandled task message type: %s", taskMsg.Type)
			return rabbitmq.NackRequeue
		}

		return rabbitmq.Ack
	}
}
