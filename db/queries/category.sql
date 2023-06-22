-- name: CreateCategory :one
INSERT INTO categories
    (name)
VALUES ($1)
RETURNING *;

-- name: GetOrCreateCategory :one
WITH new_category AS (
    INSERT INTO categories (name)
        VALUES ($1)
        ON CONFLICT (name) DO NOTHING
        RETURNING id)
SELECT id
FROM new_category
UNION
SELECT id
FROM categories
WHERE name = $1;

-- name: DeleteCategory :exec
DELETE
FROM categories
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

