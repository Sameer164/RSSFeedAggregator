-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id) VALUES (
                                                                                                       $1,
                                                                                                       $2,
                                                                                                       $3,
                                                                                                       $4,
                                                                                                       $5,
                                                                                                       $6,
                                                                                                       $7,
                                                                                                       $8
                                                                                                      )
    ON CONFLICT (url) DO NOTHING
RETURNING *;

-- name: GetPostsForUser :many
SELECT posts.* FROM users INNER JOIN feed_follows ON feed_follows.user_id = users.id AND users.id = $1 INNER JOIN posts ON posts.feed_id = feed_follows.feed_id ORDER BY posts.updated_at NULLS FIRST LIMIT $2;