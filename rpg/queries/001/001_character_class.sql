-- name: ListCharacterClasses :many
SELECT
    *
FROM
    character_class;

-- name: CreateCharacterClass :one
INSERT INTO
    character_class (id, name, description)
VALUES
    (?, ?, ?) RETURNING *;

-- name: ReadCharacterClass :one
SELECT
    *
FROM
    character_class
WHERE
    name = ?;

-- name: UpdateCharacterClass :one
UPDATE
    character_class
SET
    name = ?,
    description = ?
WHERE
    id = ? RETURNING *;

-- name: DeleteCharacterClass :exec
DELETE FROM
    character_class
WHERE
    id = ?;