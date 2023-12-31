// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: tag.sql

package db

import (
	"context"

	"github.com/lib/pq"
)

const createTag = `-- name: CreateTag :one
INSERT INTO tags
    (name)
VALUES ($1)
RETURNING id, name
`

func (q *Queries) CreateTag(ctx context.Context, name string) (Tag, error) {
	row := q.db.QueryRowContext(ctx, createTag, name)
	var i Tag
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const deleteTag = `-- name: DeleteTag :exec
DELETE
FROM tags
WHERE name = $1
`

func (q *Queries) DeleteTag(ctx context.Context, name string) error {
	_, err := q.db.ExecContext(ctx, deleteTag, name)
	return err
}

const getOrCreateTags = `-- name: GetOrCreateTags :many
WITH input_tags AS (SELECT UNNEST($1::text[]) AS name),
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
FROM created_tags
`

func (q *Queries) GetOrCreateTags(ctx context.Context, tagNames []string) ([]int32, error) {
	rows, err := q.db.QueryContext(ctx, getOrCreateTags, pq.Array(tagNames))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []int32{}
	for rows.Next() {
		var id int32
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listTagIDsByNames = `-- name: ListTagIDsByNames :many
SELECT id
FROM tags
WHERE name = ANY ($1::text[])
`

func (q *Queries) ListTagIDsByNames(ctx context.Context, tagNames []string) ([]int32, error) {
	rows, err := q.db.QueryContext(ctx, listTagIDsByNames, pq.Array(tagNames))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []int32{}
	for rows.Next() {
		var id int32
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listTags = `-- name: ListTags :many
SELECT id, name
FROM tags
ORDER BY name
LIMIT $1 OFFSET $2
`

type ListTagsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListTags(ctx context.Context, arg ListTagsParams) ([]Tag, error) {
	rows, err := q.db.QueryContext(ctx, listTags, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Tag{}
	for rows.Next() {
		var i Tag
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
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

const updateTag = `-- name: UpdateTag :one
UPDATE tags
SET name = $2
WHERE name = $1
RETURNING id, name
`

type UpdateTagParams struct {
	Name   string `json:"name"`
	Name_2 string `json:"name_2"`
}

func (q *Queries) UpdateTag(ctx context.Context, arg UpdateTagParams) (Tag, error) {
	row := q.db.QueryRowContext(ctx, updateTag, arg.Name, arg.Name_2)
	var i Tag
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}
