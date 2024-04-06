package task

import (
	"context"

	"github.com/SQUASHD/hbit/task/taskdb"
)

type (
	Service interface {
		UserTaskService
		InternalService
	}

	CreateTaskForm struct {
		CreateTaskRequest
		RequestedById string
	}

	UpdateTaskForm struct {
		taskdb.UpdateTaskParams
		TaskId        string
		RequestedById string
	}

	DeleteTaskForm struct {
		TaskId        string
		RequestedById string
	}

	UserTaskService interface {
		List(ctx context.Context, userId string) ([]DTO, error)
		Create(ctx context.Context, form CreateTaskForm) (DTO, error)
		Update(ctx context.Context, form UpdateTaskForm) (DTO, error)
		Delete(ctx context.Context, form DeleteTaskForm) error
	}

	InternalService interface {
		DeleteUserTasks(userId string) error
		CleanUp() error
	}
)
