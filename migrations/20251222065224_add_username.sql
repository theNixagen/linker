-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN username varchar(255) unique;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN username;
-- +goose StatementEnd
