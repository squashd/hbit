-- name: CreateUserSettings :one
INSERT INTO
    user_settings (
        user_id,
        theme,
        display_name,
        email,
        email_notifications
    )
VALUES
    (?, ?, ?, ?, ?) RETURNING *;

-- name: ReadUserSettings :one
SELECT
    *
FROM
    user_settings
WHERE
    user_id = ?;

-- name: UpdateUserSettings :one
UPDATE
    user_settings
SET
    theme = ?,
    display_name = ?,
    email = ?,
    email_notifications = ?
WHERE
    user_id = ? RETURNING *;

-- name: DeleteUserSettings :exec
DELETE FROM
    user_settings
WHERE
    user_id = ?;