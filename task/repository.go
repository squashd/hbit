package task

import (
	"context"
	"database/sql"

	"github.com/SQUASHD/hbit/task/taskdb"
)

type repository struct {
	queries *taskdb.Queries
	db      *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	queries := taskdb.New(db)
	return &repository{queries: queries, db: db}
}

func (r *repository) List(ctx context.Context, userId string) (Tasks, error) {
	return r.queries.ListTasks(ctx, userId)
}

func (r *repository) Create(ctx context.Context, data taskdb.CreateTaskParams) (taskdb.Task, error) {
	return r.queries.CreateTask(ctx, data)
}

func (r *repository) Read(ctx context.Context, id string) (taskdb.Task, error) {
	return r.queries.ReadTask(ctx, id)
}

func (r *repository) Update(ctx context.Context, data taskdb.UpdateTaskParams) (taskdb.Task, error) {
	return r.queries.UpdateTask(ctx, data)
}

func (r *repository) Delete(ctx context.Context, id string) error {
	return r.queries.DeleteTask(ctx, id)
}

func (r *repository) DeleteUserTasks(userId string) error {
	return r.queries.DeleteUserTasks(context.Background(), userId)
}

func (r *repository) Cleanup() error {
	return r.db.Close()
}
