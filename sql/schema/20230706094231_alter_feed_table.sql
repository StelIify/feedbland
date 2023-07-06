-- +goose Up
-- +goose StatementBegin
alter table feeds
    drop column slug,
    alter column user_id set not null;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table feeds
    add column slug varchar(250);
-- +goose StatementEnd
