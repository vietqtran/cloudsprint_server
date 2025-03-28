// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: session.sql

package sqlc

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createSession = `-- name: CreateSession :one
INSERT INTO sessions (
  id,
  account_id,
  refresh_token,
  user_agent,
  client_ip,
  expires_at
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING id, account_id, refresh_token, user_agent, client_ip, is_blocked, expires_at, status, created_at
`

type CreateSessionParams struct {
	ID           uuid.UUID `json:"id"`
	AccountID    uuid.UUID `json:"accountId"`
	RefreshToken string    `json:"refreshToken"`
	UserAgent    string    `json:"userAgent"`
	ClientIp     string    `json:"clientIp"`
	ExpiresAt    time.Time `json:"expiresAt"`
}

func (q *Queries) CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error) {
	row := q.db.QueryRowContext(ctx, createSession,
		arg.ID,
		arg.AccountID,
		arg.RefreshToken,
		arg.UserAgent,
		arg.ClientIp,
		arg.ExpiresAt,
	)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.RefreshToken,
		&i.UserAgent,
		&i.ClientIp,
		&i.IsBlocked,
		&i.ExpiresAt,
		&i.Status,
		&i.CreatedAt,
	)
	return i, err
}

const getSession = `-- name: GetSession :one
SELECT id, account_id, refresh_token, user_agent, client_ip, is_blocked, expires_at, status, created_at FROM sessions
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetSession(ctx context.Context, id uuid.UUID) (Session, error) {
	row := q.db.QueryRowContext(ctx, getSession, id)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.RefreshToken,
		&i.UserAgent,
		&i.ClientIp,
		&i.IsBlocked,
		&i.ExpiresAt,
		&i.Status,
		&i.CreatedAt,
	)
	return i, err
}

const listSessions = `-- name: ListSessions :many
SELECT id, account_id, refresh_token, user_agent, client_ip, is_blocked, expires_at, status, created_at FROM sessions
WHERE account_id = $1
ORDER BY created_at DESC
`

func (q *Queries) ListSessions(ctx context.Context, accountID uuid.UUID) ([]Session, error) {
	rows, err := q.db.QueryContext(ctx, listSessions, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Session{}
	for rows.Next() {
		var i Session
		if err := rows.Scan(
			&i.ID,
			&i.AccountID,
			&i.RefreshToken,
			&i.UserAgent,
			&i.ClientIp,
			&i.IsBlocked,
			&i.ExpiresAt,
			&i.Status,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateSessionRefreshToken = `-- name: UpdateSessionRefreshToken :one
UPDATE sessions
SET
  refresh_token = $2,
  expires_at = $3
WHERE id = $1
RETURNING id, account_id, refresh_token, user_agent, client_ip, is_blocked, expires_at, status, created_at
`

type UpdateSessionRefreshTokenParams struct {
	ID           uuid.UUID `json:"id"`
	RefreshToken string    `json:"refreshToken"`
	ExpiresAt    time.Time `json:"expiresAt"`
}

func (q *Queries) UpdateSessionRefreshToken(ctx context.Context, arg UpdateSessionRefreshTokenParams) (Session, error) {
	row := q.db.QueryRowContext(ctx, updateSessionRefreshToken, arg.ID, arg.RefreshToken, arg.ExpiresAt)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.RefreshToken,
		&i.UserAgent,
		&i.ClientIp,
		&i.IsBlocked,
		&i.ExpiresAt,
		&i.Status,
		&i.CreatedAt,
	)
	return i, err
}
