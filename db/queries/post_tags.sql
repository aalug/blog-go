-- name: AddTagToPost :exec
INSERT INTO post_tags
    (post_id, tag_id)
VALUES ($1, $2);

-- name: RemoveTagFromPost :exec
DELETE
FROM post_tags
WHERE post_id = $1
  AND tag_id = $2;

-- name: AddMultipleTagsToPost :exec
WITH input_tags AS (SELECT UNNEST(@tag_ids::int[]) AS tag_id)
INSERT
INTO post_tags (post_id, tag_id)
SELECT @post_id, tag_id
FROM input_tags;

-- name: GetTagsOfPost :many
SELECT t.*
FROM tags AS t
         JOIN post_tags AS pt ON pt.tag_id = t.id
WHERE pt.post_id = $1;