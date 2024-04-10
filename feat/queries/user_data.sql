-- name: DeleteUserFeats :exec
DELETE FROM
    user_feats
WHERE
    user_id = ?;

-- name: DeleteUserTaskLogs :exec
DELETE FROM
    task_log
WHERE
    user_id = ?;

-- name: DeleteUserQuestLogs :exec
DELETE FROM
    quest_log
WHERE
    user_id = ?;