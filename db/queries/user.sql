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

-- name: UpdateUser :one
UPDATE users
SET hashed_password     = COALESCE(sqlc.narg('hashed_password'), hashed_password),
    password_changed_at = COALESCE(sqlc.narg('password_changed_at'), password_changed_at),
    username            = COALESCE(sqlc.narg('username'), username),
    is_email_verified   = COALESCE(sqlc.narg('is_email_verified'), is_email_verified)
WHERE email = sqlc.arg('email')
RETURNING *;