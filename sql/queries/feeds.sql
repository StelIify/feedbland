-- name: CreateFeed :one
insert into feeds (name, description, url, user_id, image_id)
values($1, $2, $3, $4, $5)
returning id, created_at, name, url, user_id;

-- name: ListFeeds :many
select f.id, f.created_at, f.name, f.url, f.user_id, images.url as image_url, images.name as image_alt from feeds f
join images on images.id=f.image_id
order by created_at;

-- name: GenerateNextFeedsToFetch :many
select id, name, url from feeds
order by last_fetched_at asc nulls first
limit $1;

-- name: MarkFeedFetched :exec
update feeds
set last_fetched_at=now(), updated_at=now()
where id=$1;