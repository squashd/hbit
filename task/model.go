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
		ID          string    `json:"id"`
		UserID      string    `json:"user_id"`
		Title       string    `json:"title"`
		Text        string    `json:"text"`
		IsCompleted bool      `json:"is_completed"`
		Difficulty  string    `json:"difficulty"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}

	CreateTaskRequest struct {
		Title string `json:"title"`
	}

	UpdateTaskRequest struct{}
)

func CreateFormtoModel(form CreateTaskForm) taskdb.CreateTaskParams {
	return taskdb.CreateTaskParams{
		Title:  form.Title,
		UserID: form.RequestedById,
	}
}

func toDTO(task taskdb.Task) DTO {
	return DTO{
		ID:          task.ID,
		Title:       task.Title,
		Text:        task.Text,
		IsCompleted: task.IsCompleted,
		Difficulty:  task.Difficulty,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}

func toDTOs(tasks []taskdb.Task) []DTO {
	tasksDTO := make([]DTO, len(tasks))
	for i, task := range tasks {
		tasksDTO[i] = toDTO(task)
	}
	return tasksDTO
}
