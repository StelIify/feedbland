-- name: CreatePost :exec 
insert into posts (feed_id, title, url, description, published_at)
values ($1, $2, $3, $4, $5);

-- name: GetPostsFollowedByUser :many 
select p.id, p.title, img.url from posts p 
join feed_follows fw on fw.feed_id=p.feed_id
left join images img on img.id=p.image_id
where fw.user_id = $1
order by p.published_at desc;

-- name: GetPostsForFeed :many
select p.id, p.title, p.url, p.description, p.published_at from posts p
where feed_id=$1
order by p.published_at desc;

-- name: ListPosts :many
select p.id, f.name as feed_name, p.created_at, p.updated_at, p.title, p.url, p.description, p.published_at
from posts p
join feeds f on p.feed_id=f.id
where (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) or $1 = '')
order by published_at desc
limit $2 offset $3;

-- name: CountPosts :one
select count(*) from posts
where (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) or $1 = '');