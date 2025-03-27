-- name: CreateSession :one
INSERT INTO sessions (
  id,
  user_id,
  refresh_token,
  user_agent,
  client_ip,
  expires_at
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions
WHERE id = $1 LIMIT 1;

-- name: ListSessions :many
SELECT * FROM sessions
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateSessionRefreshToken :one
UPDATE sessions
SET
  refresh_token = $2,
  expires_at = $3
WHERE id = $1
RETURNING *;