-- name: CreateFeed :one
insert into feeds (name, description, url, user_id, image_id)
values($1, $2, $3, $4, $5)
returning id, created_at, name, url, user_id;

-- name: ListFeeds :many
select f.id, f.created_at, f.name, f.url, f.user_id, f.description,
coalesce(images.url, 'https://feebland.s3.eu-west-3.amazonaws.com/feedsImg/default-feed-image.jpg') as image_url, coalesce(images.name, 'default image') as image_alt
from feeds f
left join images on images.id=f.image_id
where (to_tsvector('simple', f.name) @@ plainto_tsquery('simple', $1) or $1 = '')
order by created_at desc
limit $2 offset $3;

-- name: GenerateNextFeedsToFetch :many
select id, name, url
from feeds
order by last_fetched_at asc nulls first
limit $1;

-- name: MarkFeedFetched :exec
update feeds
set last_fetched_at=now(), updated_at=now()
where id=$1;

-- name: CountFeeds :one
select count(*) from feeds;