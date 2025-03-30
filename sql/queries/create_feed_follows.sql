-- name: CreateFeedFollow :one
WITH inserted_feed AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id) VALUES (
                                                                                   $1,
                                                                                   $2,
                                                                                   $3,
                                                                                   $4,
                                                                                   $5)
        RETURNING *
) SELECT inserted_feed.*, users.name AS user_name, feeds.name AS feed_name FROM inserted_feed INNER JOIN users ON users.id = inserted_feed.user_id INNER JOIN feeds ON feeds.id = inserted_feed.feed_id;

-- name: GetFeedFollowsForUser :many
SELECT feeds.name FROM users INNER JOIN feed_follows ON feed_follows.user_id = users.id AND users.name = $1 INNER JOIN feeds ON feed_follows.feed_id = feeds.id;