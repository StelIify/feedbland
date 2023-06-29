-- name: CreateFeed :one
insert into feeds (name, url)
values($1, $2)
returning id, created_at, name, url, user_id;