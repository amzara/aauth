-- name: Register :exec
INSERT INTO creds (username, password) VALUES ($1, $2);

-- name: CheckUserExists :one
SELECT EXISTS(SELECT 1 FROM creds WHERE username = $1);

-- name: GetUserByUsername :one
SELECT username, password FROM creds WHERE username = $1;


