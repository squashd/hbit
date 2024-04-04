package hbit

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type (
	EventType string
	EventId   string
	UserId    string

	EventMessage struct {
		Type    EventType       `json:"type"`
		UserId  UserId          `json:"user_id"`
		EventId EventId         `json:"event_id"`
		Payload json.RawMessage `json:"payload"`
	}

	Publisher interface {
		Publish(msg EventMessage, routingKey []string) error
	}

	// All service need to implement this interface
	UserDataHandler interface {
		DeleteData(userId string) error
	}
)

func (u UserId) String() string {
	return string(u)
}

func NewUUID() string {
	return uuid.New().String()
}

func (u UserId) Valid() error {
	_, err := uuid.Parse(string(u))
	if err != nil {
		return err
	}
	return nil
}

func NewUserId() UserId {
	return UserId(NewUUID())
}

func NewEventIdWithTimestamp(domain string) EventId {
	timestamp := time.Now().UTC().Format("20060102-150405")
	return EventId(fmt.Sprintf("%s-%s-%s", domain, timestamp, NewUUID()))
}

func ExtractTimestampFromEventID(eventId string) string {
	return eventId[len(eventId)-36 : len(eventId)-6]
}

// Event types
// Because Go does not have enums... sad
const (
	TASKDONE EventType = "task_done"
	TASKUNDO EventType = "task_undo"

	AUTHDELETE EventType = "auth_delete"
	AUTHLOGIN  EventType = "auth_login"

	CHARACTERLEVELUP EventType = "character_level_up"
	QUESTCOMPLETE    EventType = "quest_complete"
	RPGREWARD        EventType = "rpg_reward"

	FEATUNLOCKED EventType = "feat_unlocked"
)

// NewEventMessage is a helper function to create a unified event message
// the payload is marshalled to json.RawMessage to allow for any type of payload
func NewEventMessage(
	eventType EventType,
	userId UserId,
	eventId EventId,
	payload any,
) (EventMessage, error) {

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return EventMessage{}, err
	}

	event := EventMessage{
		Type:    eventType,
		UserId:  userId,
		EventId: eventId,
		Payload: payloadBytes,
	}

	return event, nil
}
