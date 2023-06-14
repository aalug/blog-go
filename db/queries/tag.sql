-- name: CreateTag :one
INSERT INTO tags
    (name)
VALUES ($1)
RETURNING *;

-- name: DeleteTag :exec
DELETE
FROM tags
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

-- name: GetOrCreateTags :many
WITH input_tags AS (SELECT UNNEST(@tag_names::text[]) AS name),
     created_tags AS (
         INSERT INTO tags (name)
             SELECT name FROM input_tags
             ON CONFLICT (name) DO NOTHING
             RETURNING id)
SELECT id
FROM tags
WHERE name IN (SELECT name FROM input_tags)
UNION ALL
SELECT id
FROM created_tags;