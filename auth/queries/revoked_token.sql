-- name: FindRevokedToken :one
SELECT
    *
FROM
    revoked_token
WHERE
    token = ?;

-- name: CreateRevokedToken :exec
INSERT INTO
    revoked_token (token, user_id)
VALUES
    (?, ?);