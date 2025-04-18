-- name: SetFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id) VALUES (
                                                                           $1,
                                                                           $2,
                                                                           $3,
                                                                           $4,
                                                                           $5,
                                                                           $6
                                                                          )
RETURNING *;


-- name: GetFeeds :many
SELECT feeds.name AS feedName, feeds.url AS feedURL, users.name AS userName FROM feeds INNER JOIN users ON feeds.user_id = users.id;

-- name: GetFeedFromURL :one
SELECT * FROM feeds WHERE url = $1;