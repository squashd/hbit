-- name: ListQuests :many
SELECT
    *
FROM
    quest;

-- name: CreateQuest :one
INSERT INTO
    quest (id, title, description, requirements, rewards)
VALUES
    (uuid4(), ?, ?, ?, ?) RETURNING *;

-- name: ReadQuest :one
SELECT
    *
FROM
    quest
WHERE
    id = ?;

-- name: UpdateQuest :one
UPDATE
    quest
SET
    title = ?,
    description = ?,
    requirements = ?,
    rewards = ?
WHERE
    id = ? RETURNING *;

-- name: DeleteQuest :exec
DELETE FROM
    quest
WHERE
    id = ?;