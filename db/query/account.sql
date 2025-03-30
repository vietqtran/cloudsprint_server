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
WHERE email = $1 LIMIT 1;

-- name: GetAccountById :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: GetAccountByUserId :one
SELECT * FROM accounts
WHERE user_id = $1 LIMIT 1;

-- name: UpdateAccount :one
UPDATE accounts
SET
    hashed_password = COALESCE(sqlc.narg(hashed_password), hashed_password),
    status = COALESCE(sqlc.narg(status), status),
    updated_at = now()
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteAccount :exec
WITH account_user AS (
  SELECT user_id FROM accounts WHERE accounts.id = $1
)
DELETE FROM users
WHERE users.id = (SELECT user_id FROM account_user);