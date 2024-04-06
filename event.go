package hbit

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Event types
// Because Go does not have enums... sad
// I do not want event types to be routing keys as consumer should be able to
// subscribe to multiple event types without knowing the routing key
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
	// TODO: Consider setting all userId params as a type alias
	UserId string

	// EventMessage is a unified event message that can be published
	// to the message broker. The payload is a json.RawMessage to allow
	// for any type of payload and avoid unnecessary marshalling
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

// TODO: Refactor middleware to use this function
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

// NewEventMessage — we try a bit of type enforcement. As a treat.
func NewEventMessage(
	eventType EventType,
	userId UserId,
	eventId EventId,
	payload any,
) (EventMessage, error) {
	if payload == nil {
		payload = map[string]interface{}{}
	}

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
