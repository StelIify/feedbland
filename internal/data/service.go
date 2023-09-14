package data

import (
	"context"
	"fmt"

	"github.com/StelIify/feedbland/internal/database"
)

type CustomQueries struct {
	db database.DBTX
}

func NewCustomQueries(db database.DBTX) *CustomQueries {
	return &CustomQueries{db: db}
}

func (q *CustomQueries) ListAllPosts(ctx context.Context, filters Filters) ([]database.ListPostsRow, error) {
	listPosts := fmt.Sprintf(`
	select p.id, f.name as feed_name, p.created_at, p.updated_at, p.title, p.url, p.description, p.published_at
	from posts p
	join feeds f on p.feed_id=f.id
	where (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) or $1 = '')
	order by %s %s
	limit $2 offset $3
	`, filters.SortColumn(), filters.SortDirection())

	rows, err := q.db.Query(ctx, listPosts, filters.Title, filters.Limit, filters.Offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var items []database.ListPostsRow
	for rows.Next() {
		var i database.ListPostsRow
		if err := rows.Scan(
			&i.ID,
			&i.FeedName,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
			&i.Url,
			&i.Description,
			&i.PublishedAt,
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

func (q *CustomQueries) ListAllFeeds(ctx context.Context, filters Filters) ([]database.ListFeedsRow, error) {
	listFeeds := fmt.Sprintf(`
	select f.id, f.created_at, f.name, f.url, f.user_id, f.description,
	coalesce(images.url, 'https://feebland.s3.eu-west-3.amazonaws.com/feedsImg/default-feed-image.jpg') as image_url, coalesce(images.name, 'default image') as image_alt
	from feeds f
	left join images on images.id=f.image_id
	where (to_tsvector('simple', f.name) @@ plainto_tsquery('simple', $1) or $1 = '')
	order by %s %s
	limit $2 offset $3
	`, filters.SortColumn(), filters.SortDirection())

	rows, err := q.db.Query(ctx, listFeeds, filters.Title, filters.Limit, filters.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []database.ListFeedsRow
	for rows.Next() {
		var i database.ListFeedsRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.Name,
			&i.Url,
			&i.UserID,
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
