-- name: CreateUser :one
INSERT INTO users (
  username,
  email,
  hashed_password,
  full_name
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET
  username = COALESCE(sqlc.narg(username), username),
  email = COALESCE(sqlc.narg(email), email),
  hashed_password = COALESCE(sqlc.narg(hashed_password), hashed_password),
  full_name = COALESCE(sqlc.narg(full_name), full_name),
  updated_at = now()
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

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