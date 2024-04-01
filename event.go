package hbit

import "encoding/json"

type (
	EventMessage struct {
		Type    string          `json:"type"`
		UserID  string          `json:"userId"`
		Payload json.RawMessage `json:"payload"`
	}

	Publisher interface {
		Publish(msg EventMessage, routingKey []string, headers map[string]any) error
	}

	TaskCompleteMessage struct {
		UserID string `json:"userId"`
	}

	AuthDeleteMessage struct {
		UserID string `json:"userId"`
	}

	UserLoginMessage struct {
		UserID string `json:"userId"`
	}

	CharacterLevelUpMessage struct {
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
