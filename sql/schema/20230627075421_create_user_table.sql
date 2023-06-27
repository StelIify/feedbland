-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id bigserial primary key,
    created_at timestamp(0) with time zone not null default now(),
    updated_at timestamp(0) with time zone not null default now(),
    name text not null,
    email citext unique not null,
    password_hash bytea not null,
    activated bool not null default false,
    version integer not null default 1
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
