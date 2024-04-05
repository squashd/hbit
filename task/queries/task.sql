-- name: ListTasks :many
SELECT
    *
FROM
    task
WHERE
    user_id = ?;

-- name: CreateTask :one
INSERT INTO
    task (id, user_id, title, data)
VALUES
    (uuid4(), ?, ?, ?) RETURNING *;

-- name: ReadTask :one
SELECT
    *
FROM
    task
WHERE
    id = ?;

-- name: UpdateTask :one
UPDATE
    task
SET
    title = ?,
    text = ?,
    data = ?,
    updated_at = ?
WHERE
    id = ? RETURNING *;

-- name: DeleteTask :exec
DELETE FROM
    task
WHERE
    id = ?;

-- name: DeleteUserTasks :exec
DELETE FROM
    task
WHERE
    user_id = ?;