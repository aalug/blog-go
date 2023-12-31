-- name: AddTagToPost :exec
INSERT INTO post_tags
    (post_id, tag_id)
VALUES ($1, $2);

-- name: DeleteTagsFromPost :exec
WITH deleted_tags AS (
    DELETE FROM post_tags
        WHERE post_id = @post_id::int
            AND tag_id = ANY (@tag_ids::int[])
        RETURNING tag_id)
DELETE
FROM tags
WHERE id IN (SELECT dt.tag_id
             FROM deleted_tags dt
             WHERE dt.tag_id NOT IN (SELECT tag_id
                                     FROM post_tags));

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