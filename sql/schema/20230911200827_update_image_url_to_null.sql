-- +goose Up
-- +goose StatementBegin
alter table images
    alter column url drop not null;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table images
    alter column url set not null;
-- +goose StatementEnd
