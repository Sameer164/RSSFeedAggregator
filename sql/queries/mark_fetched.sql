-- name: MarkFeedFetched :one
UPDATE feeds
SET last_fetched_at = $2, updated_at = $2
WHERE id = $1
RETURNING *;

-- name: GetNextFeedToFetch :one
SELECT feeds.url, feeds.id FROM feed_follows INNER JOIN feeds ON feed_follows.feed_id = feeds.id AND feed_follows.user_id = $1 ORDER BY feeds.last_fetched_at NULLS FIRST LIMIT 1;