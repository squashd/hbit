// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: achievement.sql

package database

import (
	"context"
	"time"
)

const createAchievement = `-- name: CreateAchievement :one
INSERT INTO
    achievement (id, name, requirement, created_at, updated_at)
VALUES
    (?, ?, ?, ?, ?) returning id, name, requirement, created_at, updated_at
`

type CreateAchievementParams struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Requirement string    `json:"requirement"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (q *Queries) CreateAchievement(ctx context.Context, arg CreateAchievementParams) (Achievement, error) {
	row := q.db.QueryRowContext(ctx, createAchievement,
		arg.ID,
		arg.Name,
		arg.Requirement,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i Achievement
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Requirement,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteAchievement = `-- name: DeleteAchievement :exec
DELETE FROM
    achievement
WHERE
    id = ?
`

func (q *Queries) DeleteAchievement(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, deleteAchievement, id)
	return err
}

const listAchievements = `-- name: ListAchievements :many
SELECT
    id, name, requirement, created_at, updated_at
FROM
    achievement
`

func (q *Queries) ListAchievements(ctx context.Context) ([]Achievement, error) {
	rows, err := q.db.QueryContext(ctx, listAchievements)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Achievement
	for rows.Next() {
		var i Achievement
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Requirement,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const readAchievement = `-- name: ReadAchievement :one
SELECT
    id, name, requirement, created_at, updated_at
FROM
    achievement
WHERE
    id = ?
`

func (q *Queries) ReadAchievement(ctx context.Context, id string) (Achievement, error) {
	row := q.db.QueryRowContext(ctx, readAchievement, id)
	var i Achievement
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Requirement,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateAchievement = `-- name: UpdateAchievement :one
UPDATE
    achievement
SET
    name = ?,
    requirement = ?,
    updated_at = ?
WHERE
    id = ? RETURNING id, name, requirement, created_at, updated_at
`

type UpdateAchievementParams struct {
	Name        string    `json:"name"`
	Requirement string    `json:"requirement"`
	UpdatedAt   time.Time `json:"updated_at"`
	ID          string    `json:"id"`
}

func (q *Queries) UpdateAchievement(ctx context.Context, arg UpdateAchievementParams) (Achievement, error) {
	row := q.db.QueryRowContext(ctx, updateAchievement,
		arg.Name,
		arg.Requirement,
		arg.UpdatedAt,
		arg.ID,
	)
	var i Achievement
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Requirement,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
