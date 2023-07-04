-- name: CreateFeed :one
insert into feeds (name, url, user_id)
values($1, $2, $3)
returning id, created_at, name, url, user_id;

-- name: ListFeeds :many
select id, created_at, name, url, user_id from feeds
order by created_at;