package hbit

import (
	"encoding/json"
)

type (
	EventType string

	EventMessage struct {
		Type    EventType       `json:"type"`
		UserId  string          `json:"user_id"`
		EventId string          `json:"event_id"`
		Payload json.RawMessage `json:"payload"`
	}

	Publisher interface {
		Publish(msg EventMessage, routingKey []string) error
	}
)

const (
	TaskCompleteEvent EventType = "task_complete"

	AuthDeleteEvent EventType = "auth_delete"
	AuthLoginEvent  EventType = "auth_login"

	CharacterLevelUpEvent EventType = "character_level_up"

	QuestCompleteEvent EventType = "quest_complete"

	TaskRewardEvent EventType = "task_reward"

	FeatUnlocked EventType = "feat_unlocked"
)

func NewEventMessage(
	eventType EventType,
	userId, eventId string,
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
