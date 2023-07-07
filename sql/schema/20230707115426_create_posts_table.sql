-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS posts(
    id bigserial primary key,
    feed_id bigint not null references feeds(id) on delete cascade,
    created_at timestamp(0) with time zone not null default now(),
    updated_at timestamp(0) with time zone not null default now(),
    title varchar(250) not null,
    url text unique not null,
    description text not null,
    published_at timestamp(0) with time zone
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS posts;
-- +goose StatementEnd
