-- name: IsAdmin :one
SELECT
    *
FROM
    admin
WHERE
    user_id = ?;