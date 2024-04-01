// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: revoked_token.sql

package database

import (
	"context"
)

const createRevokedToken = `-- name: CreateRevokedToken :exec
INSERT INTO
    revoked_token (token, user_id)
VALUES
    (?, ?)
`

type CreateRevokedTokenParams struct {
	Token  string `json:"token"`
	UserID string `json:"user_id"`
}

func (q *Queries) CreateRevokedToken(ctx context.Context, arg CreateRevokedTokenParams) error {
	_, err := q.db.ExecContext(ctx, createRevokedToken, arg.Token, arg.UserID)
	return err
}

const findRevokedToken = `-- name: FindRevokedToken :one
SELECT
    token, user_id
FROM
    revoked_token
WHERE
    token = ?
`

func (q *Queries) FindRevokedToken(ctx context.Context, token string) (RevokedToken, error) {
	row := q.db.QueryRowContext(ctx, findRevokedToken, token)
	var i RevokedToken
	err := row.Scan(&i.Token, &i.UserID)
	return i, err
}
