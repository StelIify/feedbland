-- name: CreateFeedFollow :one
insert into feed_follows(user_id, feed_id)
values ($1, $2)
returning *;

-- name: ListFeedFollow :many
select * from feed_follows;

-- name: DeleteFeedFollow :exec
delete from feed_follows
where user_id=$1 and feed_id=$2;