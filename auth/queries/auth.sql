-- name: FindUserByUsername :one
SELECT
    *
FROM
    auth
WHERE
    username = ?;

-- name: CreateAuth :one
INSERT INTO
    auth (username, hashed_password)
VALUES
    (?, ?) RETURNING *;

-- name: UpdateUser :one
UPDATE
    auth
SET
    username = ?,
    hashed_password = ?,
    updated_at = ?
WHERE
    user_id = ? RETURNING *;

-- name: DeleteUser :exec
DELETE FROM
    auth
WHERE
    user_id = ?;