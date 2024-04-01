package task

import (
	"context"
	"database/sql"

	"github.com/SQUASHD/hbit/task/database"
)

type repository struct {
	queries *database.Queries
}

func NewRepository(db *sql.DB) Repository {
	queries := database.New(db)
	return &repository{queries: queries}
}

func (r *repository) List(ctx context.Context, userId string) (Tasks, error) {
	return r.queries.ListTasks(ctx, userId)
}

func (r *repository) Create(ctx context.Context, data CreateTaskData) (database.Task, error) {
	return r.queries.CreateTask(ctx, data)
}

func (r *repository) Read(ctx context.Context, id string) (database.Task, error) {
	return r.queries.ReadTask(ctx, id)
}

func (r *repository) Update(ctx context.Context, data UpdateTaskData) (database.Task, error) {
	return r.queries.UpdateTask(ctx, data)
}

func (r *repository) Delete(ctx context.Context, id string) error {
	return r.queries.DeleteTask(ctx, id)
}
