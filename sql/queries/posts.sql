-- name: CreatePost :exec
INSERT INTO posts(id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: GetPostsByUser :many
SELECT p.* FROM posts p
INNER JOIN feeds f ON f.id = p.feed_id
WHERE f.user_id = $1
ORDER BY p.published_at
LIMIT $2;