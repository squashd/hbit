-- name: ListQuestItems :many
SELECT
    *
FROM
    quest_item;

-- name: CreateQuestItem :one
INSERT INTO
    quest_item (quest_id, item_id, quantity_required)
VALUES
    (?, ?, ?) RETURNING *;

-- name: ReadQuestItem :one
SELECT
    *
FROM
    quest_item
WHERE
    quest_id = ?;

-- name: UpdateQuestItem :one
UPDATE
    quest_item
SET
    quest_id = ?,
    item_id = ?,
    quantity_required = ?
WHERE
    quest_id = ? RETURNING *;

-- name: DeleteQuestItem :exec
DELETE FROM
    quest_item
WHERE
    quest_id = ?;