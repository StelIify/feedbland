// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: feed_follows.sql

package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

const createFeedFollow = `-- name: CreateFeedFollow :one
insert into feed_follows(user_id, feed_id)
values ($1, $2)
returning feed_id, user_id, created_at, updated_at
`

type CreateFeedFollowParams struct {
	UserID int64 `json:"user_id"`
	FeedID int64 `json:"feed_id"`
}

func (q *Queries) CreateFeedFollow(ctx context.Context, arg CreateFeedFollowParams) (FeedFollow, error) {
	row := q.db.QueryRow(ctx, createFeedFollow, arg.UserID, arg.FeedID)
	var i FeedFollow
	err := row.Scan(
		&i.FeedID,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteFeedFollow = `-- name: DeleteFeedFollow :exec
delete from feed_follows
where user_id=$1 and feed_id=$2
`

type DeleteFeedFollowParams struct {
	UserID int64 `json:"user_id"`
	FeedID int64 `json:"feed_id"`
}

func (q *Queries) DeleteFeedFollow(ctx context.Context, arg DeleteFeedFollowParams) error {
	_, err := q.db.Exec(ctx, deleteFeedFollow, arg.UserID, arg.FeedID)
	return err
}

const listFeedFollow = `-- name: ListFeedFollow :many
select f.id, f.created_at, f.name, f.description,
coalesce(img.url, 'https://feebland.s3.eu-west-3.amazonaws.com/feedsImg/default-feed-image.jpg') as image_url, coalesce(img.name, 'default image') as image_alt from feed_follows fw
join feeds f on f.id = fw.feed_id
left join images img on img.id=f.image_id
where fw.user_id=$1
order by fw.created_at desc
`

type ListFeedFollowRow struct {
	ID          int64       `json:"id"`
	CreatedAt   time.Time   `json:"created_at"`
	Name        string      `json:"name"`
	Description pgtype.Text `json:"description"`
	ImageUrl    string      `json:"image_url"`
	ImageAlt    string      `json:"image_alt"`
}

func (q *Queries) ListFeedFollow(ctx context.Context, userID int64) ([]ListFeedFollowRow, error) {
	rows, err := q.db.Query(ctx, listFeedFollow, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListFeedFollowRow
	for rows.Next() {
		var i ListFeedFollowRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.Name,
			&i.Description,
			&i.ImageUrl,
			&i.ImageAlt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
