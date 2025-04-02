-- name: CreateUser :one
INSERT INTO users (
  email,
  first_name,
  last_name
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 AND status != 3 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 AND status != 3 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
WHERE status != 3
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: ListUsersWithoutStatus :many
SELECT * FROM users
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET
  email = COALESCE(sqlc.narg(email), email),
  first_name = COALESCE(sqlc.narg(first_name), first_name),
  last_name = COALESCE(sqlc.narg(last_name), last_name),
  updated_at = now()
WHERE id = sqlc.arg(id) AND status != 3
RETURNING *;

-- name: DeleteUser :exec
UPDATE users
SET status = 3
WHERE id = $1;