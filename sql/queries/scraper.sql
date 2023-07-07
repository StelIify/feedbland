-- name: GenerateNextFeedsToFetch :many
select id, name, url from feeds
where last_fetched_at is null or last_fetched_at < now() - interval '1 day'
order by last_fetched_at desc;

-- name: MarkFeedFetched :exec
update feeds
set last_fetched_at=$1, updated_at=$2
where id in ($3);
