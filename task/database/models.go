// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package database

import (
	"time"
)

type Task struct {
	ID        string      `json:"id"`
	UserID    string      `json:"user_id"`
	Title     string      `json:"title"`
	Text      string      `json:"text"`
	Data      interface{} `json:"data"`
	TaskType  string      `json:"task_type"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}
