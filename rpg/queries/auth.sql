-- name: DeleteUserData :exec
DELETE FROM
    user_quest
WHERE
    user_id = ?;

DELETE FROM
    character_state
WHERE
    user_id = ?;