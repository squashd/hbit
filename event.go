package hbit

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Event types
// Because Go does not have enums... sad
const (
	TASKDONE    EventType = "task_done"
	TASKUNDO    EventType = "task_undo"
	TASKCREATED EventType = "task_created"
	TASKDELETED EventType = "task_deleted"

	AUTHDELETE EventType = "auth_delete"
	AUTHLOGIN  EventType = "auth_login"

	LEVELUP       EventType = "character_level_up"
	QUESTCOMPLETE EventType = "quest_complete"
	RPGREWARD     EventType = "rpg_reward"

	FEATUNLOCKED EventType = "feat_unlocked"
)

type (
	EventType string
	EventId   string
	UserId    string

	// EventMessage is a unified event message that can be published
	// to the message broker. The payload is a json.RawMessage to allow
	// for any type of payload and avoid unnecessary marshalling
	EventMessage struct {
		Type    EventType       `json:"type"`
		UserId  string          `json:"user_id"`
		EventId EventId         `json:"event_id"`
		Payload json.RawMessage `json:"payload"`
	}

	Publisher interface {
		Publish(event EventMessage, routingKeys []string) error
	}

	// All services need to implement this interface
	UserDataHandler interface {
		DeleteData(ctx context.Context, userId string) error
	}
)

func (u UserId) String() string {
	return string(u)
}

// NewUUID is a convenience function to generate a new UUID
// SQLite does not have a UUID type, so we use a string
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

// NewEventIdWithTimestamp is a convenience function to generate a new event ID
// with a timestamp. The domain is used to identify the event type and it's
// intended to be used when propograting a new event. If an event is published
// as a response to another event the eventID should be copied from the original
func NewEventIdWithTimestamp(domain string) EventId {
	timestamp := time.Now().UTC().Format("20060102-150405")
	return EventId(fmt.Sprintf("%s-%s-%s", domain, timestamp, NewUUID()))
}

// ExtractTimestampFromEventID  — not yet sure if needed
func ExtractTimestampFromEventID(eventId string) string {
	return eventId[len(eventId)-36 : len(eventId)-6]
}

type NewEventMessageParams struct {
	EventType EventType
	UserId    string
	EventId   EventId
	Payload   any // can be nil
}

// NewEventMessage — we try a bit of type enforcement. As a treat.
func NewEventMessage(params NewEventMessageParams) (EventMessage, error) {
	if params.Payload == nil {
		params.Payload = map[string]interface{}{}
	}

	payloadBytes, err := json.Marshal(params.Payload)
	if err != nil {
		return EventMessage{}, err
	}

	event := EventMessage{
		Type:    params.EventType,
		UserId:  params.UserId,
		EventId: params.EventId,
		Payload: payloadBytes,
	}

	return event, nil
}
