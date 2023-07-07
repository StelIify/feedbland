-- name: CreateFeed :one
insert into feeds (name, url, user_id)
values($1, $2, $3)
returning id, created_at, name, url, user_id;

-- name: ListFeeds :many
select id, created_at, name, url, user_id from feeds
order by created_at;

-- name: GenerateNextFeedsToFetch :many
select id, name, url from feeds
where last_fetched_at is null or last_fetched_at < now() - interval '1 day'
order by last_fetched_at asc nulls first
limit $1;

-- name: MarkFeedFetched :exec
update feeds
set last_fetched_at=now(), updated_at=now()
where id=$1;