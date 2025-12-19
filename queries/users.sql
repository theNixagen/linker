-- name: CreateUser :one
INSERT INTO users(id, email, password, name) values(default, $1, $2,$3) RETURNING id;

-- name: GetUserByEmail :one
SELECT * FROM users where email = $1;

-- name: UpdateBio :exec
UPDATE users set bio = $1 where email = $2;

-- name: UpdateProfilePhoto :exec
UPDATE users set profile_picture = $1 where email = $2;

-- name: UpdateBannerPhoto :exec
UPDATE users set banner_picture = $1 where email = $2;
