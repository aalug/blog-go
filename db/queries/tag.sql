-- name: CreateTag :one
INSERT INTO tags
    (name)
VALUES ($1)
RETURNING *;

-- name: DeleteTag :exec
DELETE FROM tags
WHERE name = $1;

-- name: ListTags :many
SELECT *
FROM tags
ORDER BY name
LIMIT $1 OFFSET $2;

-- name: UpdateTag :one
UPDATE tags
SET name = $2
WHERE name = $1
RETURNING *;

