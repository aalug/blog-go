-- name: CreatePost :one
INSERT INTO posts
    (title, description, content, author_id, category_id, image)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetPostByID :one
SELECT p.id,
       p.title,
       p.description,
       p.content,
       u.username AS author_username,
       c.name     AS category_name,
       p.image,
       p.created_at
FROM posts p
         JOIN users u ON p.author_id = u.id
         JOIN categories c ON p.category_id = c.id
WHERE p.id = $1;

-- name: GetPostByTitle :one
SELECT p.id,
       p.title,
       p.description,
       p.content,
       u.username AS author_username,
       c.name     AS category_name,
       p.image,
       p.created_at
FROM posts p
         JOIN users u ON p.author_id = u.id
         JOIN categories c ON p.category_id = c.id
WHERE p.title = $1;

-- name: ListPosts :many
SELECT p.title,
       p.description,
       u.username AS author_username,
       c.name     AS category_name,
       p.image,
       p.created_at
FROM posts p
         JOIN users u ON p.author_id = u.id
         JOIN categories c ON p.category_id = c.id
ORDER BY p.created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListPostsByCategory :many
SELECT p.title,
       p.description,
       u.username AS author_username,
       c.name     AS category_name,
       p.image,
       p.created_at
FROM posts p
         JOIN users u ON p.author_id = u.id
         JOIN categories c ON p.category_id = c.id
WHERE c.id = $1
ORDER BY p.created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListPostsByAuthor :many
SELECT p.title,
       p.description,
       u.username AS author_username,
       c.name     AS category_name,
       p.image,
       p.created_at
FROM posts p
         JOIN users u ON p.author_id = u.id
         JOIN categories c ON p.category_id = c.id
WHERE p.author_id = $1
ORDER BY p.created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListPostsByTags :many
SELECT p.title,
       p.description,
       u.username AS author_username,
       c.name     AS category_name,
       p.image,
       p.created_at
FROM posts p
         JOIN users u ON p.author_id = u.id
         JOIN categories c ON p.category_id = c.id
         JOIN post_tags pt ON p.id = pt.post_id
         JOIN tags t ON pt.tag_id = t.id
WHERE t.id = ANY (@tag_ids::int[])
ORDER BY p.created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdatePost :one
UPDATE posts
SET title       = COALESCE($2, title),
    description = COALESCE($3, description),
    content     = COALESCE($4, content),
    category_id = COALESCE($5, category_id),
    image       = COALESCE($6, image),
    updated_at  = $7
WHERE id = $1
RETURNING *;

-- name: DeletePost :exec
DELETE
FROM posts
WHERE id = $1;
