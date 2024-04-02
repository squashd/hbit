-- name: ListCharacters :many
SELECT
    *
FROM
    character;

-- name: CreateCharacter :one
INSERT INTO
    character (user_id, class_id)
VALUES
    (?, ?) RETURNING *;

-- name: ReadCharacter :one
SELECT
    *
FROM
    character
WHERE
    user_id = ?;

-- name: UpdateCharacter :one
UPDATE
    character
SET
    class_id = ?,
    character_level = ?,
    experience = ?,
    health = ?,
    mana = ?,
    strength = ?,
    dexterity = ?,
    intelligence = ?
WHERE
    user_id = ? RETURNING *;

-- name: DeleteCharacter :exec
DELETE FROM
    character
WHERE
    user_id = ?;

-- name: GetUsersCharacters :many
SELECT
    *
FROM
    character
WHERE
    user_id = ?;