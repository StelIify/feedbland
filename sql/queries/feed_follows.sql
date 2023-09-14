-- name: CreateFeedFollow :one
insert into feed_follows(user_id, feed_id)
values ($1, $2)
returning *;

-- name: ListFeedFollow :many
select f.id, f.created_at, f.name, f.description,
coalesce(img.url, 'https://feebland.s3.eu-west-3.amazonaws.com/feedsImg/default-feed-image.jpg') as image_url, coalesce(img.name, 'default image') as image_alt from feed_follows fw
join feeds f on f.id = fw.feed_id
left join images img on img.id=f.image_id
where fw.user_id=$1
order by fw.created_at desc;

-- name: DeleteFeedFollow :exec
delete from feed_follows
where user_id=$1 and feed_id=$2;