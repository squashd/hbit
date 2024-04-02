package task

import (
	"time"

	"github.com/SQUASHD/hbit/task/taskdb"
)

const (
	RepeatDaily   RepeatType = "daily"
	RepeatWeekly  RepeatType = "weekly"
	RepeatMonthly RepeatType = "monthly"
	RepeatYearly  RepeatType = "yearly"

	Monday    Day = "m"
	Tuesday   Day = "t"
	Wednesday Day = "w"
	Thursday  Day = "th"
	Friday    Day = "f"
	Saturday  Day = "sa"
	Sunday    Day = "su"
)

type (
	Task  = taskdb.Task
	Tasks = []Task

	RepeatType string
	Day        string

	DTO struct {
		ID        string    `json:"id"`
		UserID    string    `json:"user_id"`
		Title     string    `json:"title"`
		Text      string    `json:"text"`
		Data      string    `json:"data"`
		TaskType  string    `json:"task_type"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
)

func toDTO(task Task) DTO {
	return DTO{
		ID:        task.ID,
		UserID:    task.UserID,
		Title:     task.Title,
		Text:      task.Text,
		TaskType:  task.TaskType,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}
}

func toDTOs(tasks Tasks) []DTO {
	tasksDTO := make([]DTO, len(tasks))
	for i, task := range tasks {
		tasksDTO[i] = toDTO(task)
	}
	return tasksDTO
}
