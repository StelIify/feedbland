-- name: CreateImage :one
insert into images (url, name)
values($1, $2)
returning id;