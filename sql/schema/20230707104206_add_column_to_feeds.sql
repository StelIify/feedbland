-- +goose Up
-- +goose StatementBegin
alter table feeds
    add column last_fetched_at timestamp(0) with time zone;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table feeds
    drop column last_fetched_at;
-- +goose StatementEnd
