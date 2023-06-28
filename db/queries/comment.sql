-- name: CreateComment :one
INSERT INTO "comments"
    (content, user_id, post_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListCommentsForPost :many
SELECT c.id, c.content, c.user_id, u.username, c.created_at
FROM "comments" c
         JOIN "users" u ON c.user_id = u.id
WHERE c.post_id = $1
ORDER BY c.created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateComment :one
UPDATE "comments"
SET content = $2
WHERE id = $1
RETURNING *;

-- name: DeleteComment :exec
DELETE
FROM "comments"
WHERE id = $1;

-- name: GetComment :one
SELECT id, user_id
FROM "comments"
WHERE id = $1;
