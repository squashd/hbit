// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package database

import (
	"time"
)

type Achievement struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Requirement string    `json:"requirement"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type QuestLog struct {
	EventID   string      `json:"event_id"`
	UserID    string      `json:"user_id"`
	QuestID   string      `json:"quest_id"`
	Timestamp time.Time   `json:"timestamp"`
	Payload   interface{} `json:"payload"`
}

type TaskLog struct {
	EventID   string      `json:"event_id"`
	UserID    string      `json:"user_id"`
	TaskID    string      `json:"task_id"`
	Timestamp time.Time   `json:"timestamp"`
	Payload   interface{} `json:"payload"`
}

type UserAchievement struct {
	UserID        string    `json:"user_id"`
	AchievementID string    `json:"achievement_id"`
	CreatedAt     time.Time `json:"created_at"`
}
