-- name: CreateUser :one
INSERT INTO users(id, email, password, name, username) values(default, $1, $2,$3, $4) RETURNING id;

-- name: GetUserByUsername :one
SELECT * FROM users where username = $1;

-- name: UpdateBio :exec
UPDATE users set bio = $1 where username = $2;

-- name: UpdateProfilePhoto :exec
UPDATE users set profile_picture = $1 where username = $2;

-- name: UpdateBannerPhoto :exec
UPDATE users set banner_picture = $1 where username = $2;
