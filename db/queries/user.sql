-- name: CreateUser :one
INSERT INTO users
    (username, email, hashed_password)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUser :one
SELECT *
FROM users
WHERE email = $1
LIMIT 1;

-- name: DeleteUser :exec
DELETE
FROM users
WHERE email = $1;

-- name: ListUsersContainingString :many
SELECT *
FROM users
WHERE username ILIKE '%' || @str::text || '%'
   OR email ILIKE '%' || @str::text || '%';