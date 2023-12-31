// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: comment.sql

package db

import (
	"context"
	"time"
)

const createComment = `-- name: CreateComment :one
INSERT INTO "comments"
    (content, user_id, post_id)
VALUES ($1, $2, $3)
RETURNING id, content, user_id, post_id, created_at
`

type CreateCommentParams struct {
	Content string `json:"content"`
	UserID  int32  `json:"user_id"`
	PostID  int32  `json:"post_id"`
}

func (q *Queries) CreateComment(ctx context.Context, arg CreateCommentParams) (Comment, error) {
	row := q.db.QueryRowContext(ctx, createComment, arg.Content, arg.UserID, arg.PostID)
	var i Comment
	err := row.Scan(
		&i.ID,
		&i.Content,
		&i.UserID,
		&i.PostID,
		&i.CreatedAt,
	)
	return i, err
}

const deleteComment = `-- name: DeleteComment :exec
DELETE
FROM "comments"
WHERE id = $1
`

func (q *Queries) DeleteComment(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteComment, id)
	return err
}

const getComment = `-- name: GetComment :one
SELECT id, user_id
FROM "comments"
WHERE id = $1
`

type GetCommentRow struct {
	ID     int64 `json:"id"`
	UserID int32 `json:"user_id"`
}

func (q *Queries) GetComment(ctx context.Context, id int64) (GetCommentRow, error) {
	row := q.db.QueryRowContext(ctx, getComment, id)
	var i GetCommentRow
	err := row.Scan(&i.ID, &i.UserID)
	return i, err
}

const listCommentsForPost = `-- name: ListCommentsForPost :many
SELECT c.id, c.content, c.user_id, u.username, c.created_at
FROM "comments" c
         JOIN "users" u ON c.user_id = u.id
WHERE c.post_id = $1
ORDER BY c.created_at DESC
LIMIT $2 OFFSET $3
`

type ListCommentsForPostParams struct {
	PostID int32 `json:"post_id"`
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type ListCommentsForPostRow struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	UserID    int32     `json:"user_id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

func (q *Queries) ListCommentsForPost(ctx context.Context, arg ListCommentsForPostParams) ([]ListCommentsForPostRow, error) {
	rows, err := q.db.QueryContext(ctx, listCommentsForPost, arg.PostID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListCommentsForPostRow{}
	for rows.Next() {
		var i ListCommentsForPostRow
		if err := rows.Scan(
			&i.ID,
			&i.Content,
			&i.UserID,
			&i.Username,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateComment = `-- name: UpdateComment :one
UPDATE "comments"
SET content = $2
WHERE id = $1
RETURNING id, content, user_id, post_id, created_at
`

type UpdateCommentParams struct {
	ID      int64  `json:"id"`
	Content string `json:"content"`
}

func (q *Queries) UpdateComment(ctx context.Context, arg UpdateCommentParams) (Comment, error) {
	row := q.db.QueryRowContext(ctx, updateComment, arg.ID, arg.Content)
	var i Comment
	err := row.Scan(
		&i.ID,
		&i.Content,
		&i.UserID,
		&i.PostID,
		&i.CreatedAt,
	)
	return i, err
}
