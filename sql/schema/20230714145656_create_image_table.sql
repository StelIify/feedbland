-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS images(
    id bigserial primary key,
    url text unique not null,
    name varchar(150),
    created_at timestamp(0) with time zone not null default now(),
    updated_at timestamp(0) with time zone not null default now()
);

alter table feeds
    add column image_id bigint references images(id);

alter table posts
    add column image_id bigint references images(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS images;

alter table feeds
    drop column image_id;

alter table posts
    drop column image_fd;
-- +goose StatementEnd

