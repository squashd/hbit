-- name: DeleteUserTaskData :exec
DELETE FROM
    task
WHERE
    user_id = ?;