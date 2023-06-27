-- +goose Up
-- +goose StatementBegin
alter table users
    alter column activated set default false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table users
    alter column activated drop default;
-- +goose StatementEnd
