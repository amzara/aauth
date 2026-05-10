-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS creds (
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS creds;
-- +goose StatementEnd
