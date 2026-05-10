-- +goose Up
-- +goose StatementBegin
INSERT INTO creds (username,password) VALUES ('abc','123') ON CONFLICT DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM creds WHERE username = 'abc'; 
-- +goose StatementEnd
