-- name: CreateCategory :one
INSERT INTO categories
    (name)
VALUES ($1)
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories
WHERE name = $1;

-- name: ListCategories :many
SELECT *
FROM categories
ORDER BY name
LIMIT $1 OFFSET $2;

-- name: UpdateCategory :one
UPDATE categories
SET name = $2
WHERE name = $1
RETURNING *;

