// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: feat.sql

package featdb

import (
	"context"
	"time"
)

const listUserFeats = `-- name: ListUserFeats :many
SELECT
    feat.id, feat.name, feat.requirement, feat.created_at, feat.updated_at,
    user_feats.user_id,
    user_feats.created_at AS achieved_at
FROM
    feat
    LEFT JOIN user_feats ON user_feats.feat_id = feat.id
WHERE
    user_id = ?
`

type ListUserFeatsRow struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Requirement string     `json:"requirement"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	UserID      *string    `json:"user_id"`
	AchievedAt  *time.Time `json:"achieved_at"`
}

func (q *Queries) ListUserFeats(ctx context.Context, userID string) ([]ListUserFeatsRow, error) {
	rows, err := q.db.QueryContext(ctx, listUserFeats, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListUserFeatsRow
	for rows.Next() {
		var i ListUserFeatsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Requirement,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
			&i.AchievedAt,
		); err != nil {
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