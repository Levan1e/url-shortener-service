-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS urls 
(
original_url TEXT PRIMARY KEY,
short_url TEXT UNIQUE NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS urls
-- +goose StatementEnd
