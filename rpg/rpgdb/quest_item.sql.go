// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: quest_item.sql

package rpgdb

import (
	"context"
)

const createQuestItem = `-- name: CreateQuestItem :one
INSERT INTO
    quest_item (quest_id, item_id, quantity_required)
VALUES
    (?, ?, ?) RETURNING quest_id, item_id, quantity_required
`

type CreateQuestItemParams struct {
	QuestID          string `json:"quest_id"`
	ItemID           string `json:"item_id"`
	QuantityRequired int64  `json:"quantity_required"`
}

func (q *Queries) CreateQuestItem(ctx context.Context, arg CreateQuestItemParams) (QuestItem, error) {
	row := q.db.QueryRowContext(ctx, createQuestItem, arg.QuestID, arg.ItemID, arg.QuantityRequired)
	var i QuestItem
	err := row.Scan(&i.QuestID, &i.ItemID, &i.QuantityRequired)
	return i, err
}

const deleteQuestItem = `-- name: DeleteQuestItem :exec
DELETE FROM
    quest_item
WHERE
    quest_id = ?
`

func (q *Queries) DeleteQuestItem(ctx context.Context, questID string) error {
	_, err := q.db.ExecContext(ctx, deleteQuestItem, questID)
	return err
}

const listQuestItems = `-- name: ListQuestItems :many
SELECT
    quest_id, item_id, quantity_required
FROM
    quest_item
`

func (q *Queries) ListQuestItems(ctx context.Context) ([]QuestItem, error) {
	rows, err := q.db.QueryContext(ctx, listQuestItems)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []QuestItem
	for rows.Next() {
		var i QuestItem
		if err := rows.Scan(&i.QuestID, &i.ItemID, &i.QuantityRequired); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const readQuestItem = `-- name: ReadQuestItem :one
SELECT
    quest_id, item_id, quantity_required
FROM
    quest_item
WHERE
    quest_id = ?
`

func (q *Queries) ReadQuestItem(ctx context.Context, questID string) (QuestItem, error) {
	row := q.db.QueryRowContext(ctx, readQuestItem, questID)
	var i QuestItem
	err := row.Scan(&i.QuestID, &i.ItemID, &i.QuantityRequired)
	return i, err
}

const updateQuestItem = `-- name: UpdateQuestItem :one
UPDATE
    quest_item
SET
    quest_id = ?,
    item_id = ?,
    quantity_required = ?
WHERE
    quest_id = ? RETURNING quest_id, item_id, quantity_required
`

type UpdateQuestItemParams struct {
	QuestID          string `json:"quest_id"`
	ItemID           string `json:"item_id"`
	QuantityRequired int64  `json:"quantity_required"`
	QuestID_2        string `json:"quest_id_2"`
}

func (q *Queries) UpdateQuestItem(ctx context.Context, arg UpdateQuestItemParams) (QuestItem, error) {
	row := q.db.QueryRowContext(ctx, updateQuestItem,
		arg.QuestID,
		arg.ItemID,
		arg.QuantityRequired,
		arg.QuestID_2,
	)
	var i QuestItem
	err := row.Scan(&i.QuestID, &i.ItemID, &i.QuantityRequired)
	return i, err
}