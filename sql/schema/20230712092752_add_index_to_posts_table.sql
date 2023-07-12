-- +goose Up
-- +goose StatementBegin
create index if not exists posts_title_idx on posts using GIN (to_tsvector('simple', title));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index if exists posts_title_idx;
-- +goose StatementEnd
