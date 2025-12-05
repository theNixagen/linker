-- name: CreateUser :one
INSERT INTO users(id, email, password, name) values(default, $1, $2,$3) RETURNING id;

-- name: GetUserByEmail :one
SELECT * FROM users where email = $1;
