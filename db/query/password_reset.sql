-- name: CreatePasswordReset :one
INSERT INTO password_resets (
    id,
    email,
    token,
    expires_at
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetPasswordResetByToken :one
SELECT * FROM password_resets
WHERE token = $1 AND expires_at > NOW() AND used = false
LIMIT 1;

-- name: MarkPasswordResetUsed :exec
UPDATE password_resets
SET used = true
WHERE id = $1; 