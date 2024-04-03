package hbit

import "encoding/json"

type (
	EventType string

	EventMessage struct {
		Type    EventType       `json:"type"`
		UserId  string          `json:"user_id"`
		EventId string          `json:"event_id"`
		Payload json.RawMessage `json:"payload"`
	}

	Publisher interface {
		Publish(msg EventMessage, routingKey []string, headers map[string]any) error
	}

	TaskCompleteMessage struct {
		TaskID string `json:"userId"`
	}

	AuthDeleteMessage struct {
		UserID string `json:"userId"`
	}

	UserLoginMessage struct {
		UserID string `json:"userId"`
	}

	QuestCompleteMessage struct {
		UserID string `json:"userId"`
	}

	AchievementUnlockMessage struct {
		UserID string `json:"userId"`
	}
	UserDeleteMessage struct {
		UserID string `json:"userId"`
	}
)

const (
	TaskCompleteEvent EventType = "task_complete"

	AuthDeleteEvent EventType = "auth_delete"
	AuthLoginEvent  EventType = "auth_login"

	CharacterLevelUpEvent EventType = "character_level_up"

	QuestCompleteEvent EventType = "quest_complete"

	TaskRewardEvent EventType = "task_reward"

	AchievementUnlockEvent EventType = "achievement_unlock"
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

	msg := EventMessage{
		Type:    eventType,
		UserId:  userId,
		EventId: eventId,
		Payload: payloadBytes,
	}

	return msg, nil
}
