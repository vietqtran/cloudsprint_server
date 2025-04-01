// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: account.sql

package sqlc

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createAccount = `-- name: CreateAccount :one
INSERT INTO accounts (
    email,
    hashed_password,
    user_id
) VALUES (
    $1, $2, $3
) RETURNING id, email, hashed_password, user_id, status, reset_password_token, reset_password_token_expires_at, login_failed_attempts, first_failed_login_at, verify_email_token, created_at, updated_at
`

type CreateAccountParams struct {
	Email          string    `json:"email"`
	HashedPassword string    `json:"hashedPassword"`
	UserID         uuid.UUID `json:"userId"`
}

func (q *Queries) CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error) {
	row := q.db.QueryRowContext(ctx, createAccount, arg.Email, arg.HashedPassword, arg.UserID)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.HashedPassword,
		&i.UserID,
		&i.Status,
		&i.ResetPasswordToken,
		&i.ResetPasswordTokenExpiresAt,
		&i.LoginFailedAttempts,
		&i.FirstFailedLoginAt,
		&i.VerifyEmailToken,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteAccount = `-- name: DeleteAccount :exec
WITH account_user AS (
  SELECT user_id FROM accounts WHERE accounts.id = $1
)
DELETE FROM users
WHERE users.id = (SELECT user_id FROM account_user)
`

func (q *Queries) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteAccount, id)
	return err
}

const getAccountByEmail = `-- name: GetAccountByEmail :one
SELECT id, email, hashed_password, user_id, status, reset_password_token, reset_password_token_expires_at, login_failed_attempts, first_failed_login_at, verify_email_token, created_at, updated_at FROM accounts
WHERE email = $1 LIMIT 1
`

func (q *Queries) GetAccountByEmail(ctx context.Context, email string) (Account, error) {
	row := q.db.QueryRowContext(ctx, getAccountByEmail, email)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.HashedPassword,
		&i.UserID,
		&i.Status,
		&i.ResetPasswordToken,
		&i.ResetPasswordTokenExpiresAt,
		&i.LoginFailedAttempts,
		&i.FirstFailedLoginAt,
		&i.VerifyEmailToken,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getAccountById = `-- name: GetAccountById :one
SELECT id, email, hashed_password, user_id, status, reset_password_token, reset_password_token_expires_at, login_failed_attempts, first_failed_login_at, verify_email_token, created_at, updated_at FROM accounts
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetAccountById(ctx context.Context, id uuid.UUID) (Account, error) {
	row := q.db.QueryRowContext(ctx, getAccountById, id)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.HashedPassword,
		&i.UserID,
		&i.Status,
		&i.ResetPasswordToken,
		&i.ResetPasswordTokenExpiresAt,
		&i.LoginFailedAttempts,
		&i.FirstFailedLoginAt,
		&i.VerifyEmailToken,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getAccountByUserId = `-- name: GetAccountByUserId :one
SELECT id, email, hashed_password, user_id, status, reset_password_token, reset_password_token_expires_at, login_failed_attempts, first_failed_login_at, verify_email_token, created_at, updated_at FROM accounts
WHERE user_id = $1 LIMIT 1
`

func (q *Queries) GetAccountByUserId(ctx context.Context, userID uuid.UUID) (Account, error) {
	row := q.db.QueryRowContext(ctx, getAccountByUserId, userID)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.HashedPassword,
		&i.UserID,
		&i.Status,
		&i.ResetPasswordToken,
		&i.ResetPasswordTokenExpiresAt,
		&i.LoginFailedAttempts,
		&i.FirstFailedLoginAt,
		&i.VerifyEmailToken,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateAccount = `-- name: UpdateAccount :one
UPDATE accounts
SET
    hashed_password = COALESCE($1, hashed_password),
    status = COALESCE($2, status),
    updated_at = now()
WHERE id = $3
RETURNING id, email, hashed_password, user_id, status, reset_password_token, reset_password_token_expires_at, login_failed_attempts, first_failed_login_at, verify_email_token, created_at, updated_at
`

type UpdateAccountParams struct {
	HashedPassword sql.NullString `json:"hashedPassword"`
	Status         sql.NullInt32  `json:"status"`
	ID             uuid.UUID      `json:"id"`
}

func (q *Queries) UpdateAccount(ctx context.Context, arg UpdateAccountParams) (Account, error) {
	row := q.db.QueryRowContext(ctx, updateAccount, arg.HashedPassword, arg.Status, arg.ID)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.HashedPassword,
		&i.UserID,
		&i.Status,
		&i.ResetPasswordToken,
		&i.ResetPasswordTokenExpiresAt,
		&i.LoginFailedAttempts,
		&i.FirstFailedLoginAt,
		&i.VerifyEmailToken,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
