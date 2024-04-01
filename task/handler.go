package task

import (
	"context"
	"encoding/json"
)

type (
	CreateTaskForm struct {
		CreateTaskData
		RequestedById string
	}

	UpdateTaskForm struct {
		UpdateTaskData
		TaskId        string
		RequestedById string
	}

	DeleteTaskForm struct {
		TaskId        string
		RequestedById string
	}

	Service interface {
		List(context context.Context, userId string) ([]DTO, error)
		Create(context context.Context, form CreateTaskForm) (DTO, error)
		Update(context context.Context, form UpdateTaskForm) (DTO, error)
		Delete(context context.Context, form DeleteTaskForm) error

		DeleteTasks(msg json.RawMessage) error
	}
)
