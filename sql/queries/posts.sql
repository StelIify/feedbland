-- name: CreatePost :exec 
insert into posts (feed_id, title, url, description, published_at)
values ($1, $2, $3, $4, $5);

-- name: GetPostsFollowedByUser :many 
select p.* from posts p 
join feed_follows fw on fw.feed_id=p.feed_id
where fw.user_id = $1
order by p.published_at desc;

-- name: ListPosts :many
select * from posts
where (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) or $1 = '')
order by published_at desc
limit $2 offset $3;