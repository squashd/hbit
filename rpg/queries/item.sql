-- name: ListItems :many
SELECT
    *
FROM
    item;

-- name: CreateItem :one
INSERT INTO
    item (id, name, description, item_type, attributes)
VALUES
    (?, ?, ?, ?, ?) RETURNING *;

-- name: ReadItem :one
SELECT
    *
FROM
    item
WHERE
    id = ?;

-- name: UpdateItem :one
UPDATE
    item
SET
    name = ?,
    description = ?,
    item_type = ?,
    attributes = ?
WHERE
    id = ? RETURNING *;

-- name: DeleteItem :exec
DELETE FROM
    item
WHERE
    id = ?;