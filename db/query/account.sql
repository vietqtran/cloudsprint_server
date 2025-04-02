-- name: CreateAccount :one
INSERT INTO accounts (
    email,
    hashed_password,
    user_id
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetAccountByEmail :one
SELECT * FROM accounts
WHERE email = $1 AND status != 3 LIMIT 1;

-- name: GetAccountById :one
SELECT * FROM accounts
WHERE id = $1 AND status != 3 LIMIT 1;

-- name: GetAccountByUserId :one
SELECT * FROM accounts
WHERE user_id = $1 AND status != 3 LIMIT 1;

-- name: ListAccounts :many
SELECT * FROM accounts
WHERE status != 3
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: ListAccountWithoutStatus :many
SELECT * FROM accounts
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateAccount :one
UPDATE accounts
SET
    hashed_password = COALESCE(sqlc.narg(hashed_password), hashed_password),
    status = COALESCE(sqlc.narg(status), status),
    updated_at = now()
WHERE id = sqlc.arg(id) AND status != 3
RETURNING *;

-- name: DeleteAccount :exec
UPDATE accounts
SET status = 3
WHERE accounts.id = $1;