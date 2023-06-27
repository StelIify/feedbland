-- name: CreateUser :one
insert into users (name, email, password_hash)
values ($1, $2, $3)
returning id;