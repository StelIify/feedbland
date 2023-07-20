-- +goose Up
-- +goose StatementBegin
alter table feeds
    add column description text;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table feeds
    drop column description;
-- +goose StatementEnd
