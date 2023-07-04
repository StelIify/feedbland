-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS feed_follows(
    feed_id bigint references feeds(id),
    user_id bigint references users(id),
    created_at timestamp(0) with time zone not null default now(),
    updated_at timestamp(0) with time zone not null default now(),
    primary key(feed_id, user_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS feed_follows;
-- +goose StatementEnd
