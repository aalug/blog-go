-- name: AddTagToPost :exec
INSERT INTO post_tags
    (post_id, tag_id)
VALUES ($1, $2);

-- name: RemoveTagFromPost :exec
DELETE
FROM post_tags
WHERE post_id = $1
  AND tag_id = $2;