// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: feeds.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createFeed = `-- name: CreateFeed :one
insert into feeds (name, url)
values($1, $2)
returning id, created_at, name, url, user_id
`

type CreateFeedParams struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type CreateFeedRow struct {
	ID        int64              `json:"id"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	Name      string             `json:"name"`
	Url       string             `json:"url"`
	UserID    pgtype.Int8        `json:"user_id"`
}

func (q *Queries) CreateFeed(ctx context.Context, arg CreateFeedParams) (CreateFeedRow, error) {
	row := q.db.QueryRow(ctx, createFeed, arg.Name, arg.Url)
	var i CreateFeedRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Name,
		&i.Url,
		&i.UserID,
	)
	return i, err
}
