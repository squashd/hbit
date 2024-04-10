-- name: DeleteUserData :exec
DELETE FROM
    user_settings
WHERE
    user_id = ?;