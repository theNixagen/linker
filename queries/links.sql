-- name: CreateLink :exec
INSERT INTO links(
  user_id,
	url,
	title,
	description,
	created_at
) values(
  $1,$2,$3,$4,NOW()
);

-- name: FindAllLinksFromAUser :many
SELECT * FROM links where user_id = $1;


