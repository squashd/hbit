-- name: CreateCharacter :one
INSERT INTO
    character_state (user_id, class_id, event_id)
VALUES
    (?, ?, ?) RETURNING *;

-- name: ReadCharacter :one
SELECT
    *
FROM
    character_state
WHERE
    user_id = ?
ORDER BY
    timestamp DESC
LIMIT
    1;

-- name: UpdateCharacter :one
INSERT INTO
    character_state (
        user_id,
        class_id,
        event_id,
        character_level,
        experience,
        health,
        mana,
        strength,
        dexterity,
        intelligence
    )
VALUES
    (
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?
    ) RETURNING *;

-- name: DeleteCharacter :exec
DELETE FROM
    character_state
WHERE
    user_id = ?;