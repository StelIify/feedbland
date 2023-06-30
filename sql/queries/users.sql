-- name: CreateUser :one
insert into users (name, email, password_hash)
values ($1, $2, $3)
returning id, created_at, name, email, activated;

-- name: GetUserByEmail :one
select id, created_at, name, email, password_hash, activated, version
from users
where email = $1;

-- name: UpdateUser :one
update users
set name=$1, email=$2, password_hash=$3, activated=$4, version= version + 1
where id=$5 and version=$6 
returning version;

-- name: GetUserByToken :one
select id, created_at, name, email, password_hash, activated, version
from users
join tokens on tokens.user_id = users.id
where tokens.hash=$1
and tokens.scope=$2
and tokens.expiry > $3;