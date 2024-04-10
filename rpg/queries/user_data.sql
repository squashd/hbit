-- name: DeleteUserQuestData :exec
DELETE FROM
    user_quest
WHERE
    user_id = ?;

-- name: DeleteUserCharacterData :exec
DELETE FROM
    character_state
WHERE
    user_id = ?;