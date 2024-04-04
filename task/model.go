package task

import (
	"time"

	"github.com/SQUASHD/hbit/task/taskdb"
)

const (
	DAILY   RepeatType = "daily"
	WEEKLY  RepeatType = "weekly"
	MONTHLY RepeatType = "monthly"
	YEARLY  RepeatType = "yearly"

	MONDAY    Day = "m"
	TUESDAY   Day = "t"
	WEDNESDAY Day = "w"
	THURSDAY  Day = "th"
	FRIDAY    Day = "f"
	SATURDAY  Day = "sa"
	SUNDAY    Day = "su"
)

type (
	RepeatType string
	Day        string

	DTO struct {
		ID        string    `json:"id"`
		Title     string    `json:"title"`
		Text      string    `json:"text"`
		Data      string    `json:"data"`
		TaskType  string    `json:"task_type"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
)

func toDTO(task taskdb.Task) DTO {
	return DTO{
		ID:        task.ID,
		Title:     task.Title,
		Text:      task.Text,
		TaskType:  task.TaskType,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}
}

func toDTOs(tasks []taskdb.Task) []DTO {
	tasksDTO := make([]DTO, len(tasks))
	for i, task := range tasks {
		tasksDTO[i] = toDTO(task)
	}
	return tasksDTO
}
