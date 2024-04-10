-- name: GetUserRpgState :one
SELECT
    char_state.*,
    user_quest.quest_id AS user_quest_id,
    quest_info.quest_type,
    quest_info.description,
    quest_info.title,
    quest_info.details AS quest_details,
    user_quest.details AS user_quest_details,
    user_quest.completed AS quest_completed
FROM
    character_state char_state
    LEFT JOIN (
        SELECT
            *
        FROM
            user_quest
        WHERE
            user_quest.user_id = ?
        ORDER BY
            timestamp DESC
        LIMIT
            1
    ) AS user_quest ON char_state.user_id = user_quest.user_id
    INNER JOIN quest AS quest_info ON user_quest.quest_id = quest_info.quest_id
WHERE
    char_state.user_id = ?
ORDER BY
    char_state.timestamp DESC
LIMIT
    1;