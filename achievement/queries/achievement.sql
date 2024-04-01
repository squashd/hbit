-- +goose Up
CREATE TABLE IF NOT EXISTS achievement (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    requirement TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- name: ListAchievements :many
SELECT
    *
FROM
    achievement;

-- name: CreateAchievement :one
INSERT INTO
    achievement (id, name, requirement, created_at, updated_at)
VALUES
    (?, ?, ?, ?, ?) returning *;

-- name: ReadAchievement :one
SELECT
    *
FROM
    achievement
WHERE
    id = ?;

-- name: UpdateAchievement :one
UPDATE
    achievement
SET
    name = ?,
    requirement = ?,
    updated_at = ?
WHERE
    id = ? RETURNING *;

-- name: DeleteAchievement :exec
DELETE FROM
    achievement
WHERE
    id = ?;